package typer

import (
	"fmt"
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
		case "ctrl+w", "ctrl+\\":
			// ctrl+\\: ctrl+backspace
			m.app.DeleteWord()
		case "backspace", "ctrl+h":
			m.app.DeleteChar()
		default:
			char := msg.Runes[0]
			if unicode.IsPrint(char) {
				m.app.HandleKey(char)
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.width).Render(m.renderText()))
			}
		}
	default:
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.inputView(), m.headerView(), m.viewport.View(), m.footerView())
}

func (m Model) renderText() string {
	var buf strings.Builder
	currentWordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9BCED7")).Bold(true)
	for idx, word := range m.app.Words() {
		text := string(word.Text)
		if idx == m.app.CurrentWordIndex() {
			text = currentWordStyle.Render(text)
		}
		buf.WriteString(text)
		buf.WriteString(" ")
	}
	return buf.String()
}

func (m Model) headerView() string {
	title := titleStyle.Render("Chapter")
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
