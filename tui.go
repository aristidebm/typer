package typer

import (
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
		m.viewport = viewport.New(msg.Width, msg.Height)
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
			}
		}
	default:
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) renderText() string {
	var buf strings.Builder
	for _, word := range m.app.Words() {
		buf.WriteString(string(word.Text))
		buf.WriteString(" ")
	}
	return buf.String()
}
