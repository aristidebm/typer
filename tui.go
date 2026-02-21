package typer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	app      *App
	width    int
	height   int
	viewport viewport.Model
	renderer *lipgloss.Renderer
}

func NewModel(app *App, tty io.Writer) Model {
	return Model{
		app: app,
		// For more information on why we need a custom renderer, check this link
		// https://github.com/charmbracelet/lipgloss?tab=readme-ov-file#custom-renderers
		// but basically I have noticed that without a custom renderer and the use of the tty
		// colors disappear when typer is called in bash subshell to capture it's output.
		renderer: lipgloss.NewRenderer(tty),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) titleStyle() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return m.renderer.NewStyle().BorderStyle(b).Padding(0, 1)
}

func (m Model) infoStyle() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	return m.titleStyle().BorderStyle(b)
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
		m.viewport.SetContent(m.renderer.NewStyle().Width(m.width).Render(m.renderText()))
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
	currentWordStyle := m.renderer.NewStyle().Foreground(lipgloss.Color("#9BCED7")).Bold(true)
	var text string
	for idx, word := range m.app.Words() {
		text = string(word.Text)
		// nothing is typed by the user yet
		if len(word.Progress) == 0 {
			if idx == m.app.CurrentWordIndex() {
				text = currentWordStyle.Render(text)
			} else {
				text = m.renderer.NewStyle().Foreground(lipgloss.Color("#FFFFF")).Render(text)
			}
		} else {
			currentWordStyle := m.renderer.NewStyle().Foreground(lipgloss.Color("#9BCED7")).Bold(true)
			missingStyle := m.renderer.NewStyle().Foreground(lipgloss.Color("#EA6F91"))
			correctStyle := m.renderer.NewStyle().Foreground(lipgloss.Color("#31748E"))

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
	title := m.titleStyle().Render("Page")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := m.infoStyle().Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
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
	m.viewport.SetContent(m.renderer.NewStyle().Width(m.width).Render(m.renderText()))
}
