package typer

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type Model struct {
	app      *App
	width    int
	height   int
	viewport viewport.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		inputHeight := lipgloss.Height(m.inputView())
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight + inputHeight

		m.viewport = viewport.New(m.width, m.height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		// wrap text with lipgloss.NewStyle().Width(m.width).Render(...)
		// https://github.com/charmbracelet/bubbles/issues/56#issuecomment-1073306054
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.width).Render(m.renderText()))
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+w", "ctrl+backspace":
			// ctrl+\\: ctrl+backspace
			m.app.DeleteWord()
			m.updateViewport()
		case "backspace", "ctrl+h":
			m.app.DeleteChar()
			m.updateViewport()
		default:
			if len(msg.Runes) > 0 && unicode.IsPrint(msg.Runes[0]) {
				m.app.HandleKey(msg.Runes[0])
				m.updateViewport()
			}
		}
	default:
	}
	if m.app.IsCompleted() {
		m.app.ComputeResult()
		f, err := os.Create("/tmp/result.json")
		if err == nil {
			if err := m.app.Encode(f); err != nil {
				fmt.Fprint(f, "cannot encode")
			}
		}
		defer f.Close()
		return m, tea.Quit
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView(), m.inputView())
}

func (m Model) renderText() string {
	var renderer strings.Builder
	currentWordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9BCED7")).Bold(true)
	var text string
	for idx, word := range m.app.Words() {
		text = string(word.Text)
		// nothing is typed by the user yet
		if len(word.Progress) == 0 {
			if idx == m.app.CurrentWordIndex() {
				text = currentWordStyle.Render(text)
			} else {
				text = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFF")).Render(text)
			}
		} else {
			currentWordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9BCED7")).Bold(true)
			missingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EA6F91"))
			correctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#31748E"))

			var buf strings.Builder
			for idx, key := range word.Text {
				if 0 <= idx && idx < len(word.Progress) && word.Progress[idx] == key {
					buf.WriteString(correctStyle.Render(string(key)))
				} else if 0 <= idx && idx < len(word.Progress) && word.Progress[idx] != key {
					buf.WriteString(missingStyle.Render(string(key)))
				} else {
					buf.WriteString(currentWordStyle.Render(string(key)))
				}
			}
			if len(word.Progress) > len(word.Text) {
				// mark potential remaining charathers in Progress as missings
				for _, key := range word.Progress[len(word.Text):] {
					buf.WriteString(missingStyle.Render(string(key)))
				}
			}
			text = buf.String()
		}
		renderer.WriteString(text)
		renderer.WriteString(" ")
	}
	return renderer.String()
}

func (m Model) headerView() string {
	title := titleStyle.Render("Page")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) inputView() string {
	index := m.app.CurrentWordIndex()
	if index == -1 {
		return ">"
	}
	word := m.app.Words()[index]
	return "> " + string(word.Progress)
}

func (m *Model) updateViewport() {
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.width).Render(m.renderText()))
}
