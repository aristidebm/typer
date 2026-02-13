package typer

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "textfmt [text]",
	Short: "A simple text formatter",
	Long: `textfmt is a CLI tool for formatting and analyzing text.

You can provide text as arguments or pipe it in via stdin.
Multiple formatting options can be applied simultaneously.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func runTUI() error {
	app := &App{}
	r := strings.NewReader(`
	Lorem Ipsum is simply dummy text of the printing and typesetting industry.
	Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
	when an unknown printer took a galley of type and scrambled it to make a type specimen book.
	It has survived not only five centuries, but also the leap into electronic typesetting,
	remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages,
	and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
	`)
	app.CreateSession(r)
	m := Model{app: app}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
