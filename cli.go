package typer

import (
	"fmt"
	"io"
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
		source, err := GetSource(args)
		if err != nil {
			return fmt.Errorf("cannot read source using %v: %w", args, err)
		}
		defer source.Close()
		return runTUI(source)
	},
}

func runTUI(r io.Reader) error {
	app := &App{}

	// Open TTY directly — stdin may be a pipe, stdout is reserved for output.
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("cannot open /dev/tty: %w", err)
	}
	defer tty.Close()

	if err := app.CreateSession(r); err != nil {
		return err
	}

	// Pass TTY to model so lipgloss detects color support against the real terminal.
	m := NewModel(app, tty)
	p := tea.NewProgram(m,
		tea.WithInput(tty),
		tea.WithOutput(tty), // draw on TTY directly — keeps stdout clean
	)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func GetSource(args []string) (io.ReadCloser, error) {
	if len(args) == 0 {
		return generateRandom()
	}

	if len(args) >= 1 && args[0] == "-" {
		return os.Stdin, nil
	}

	fp, err := os.Open(args[0])
	if err != nil {
		return nil, err
	}
	return fp, nil
}

func generateRandom() (io.ReadCloser, error) {
	return &DummyReadCloser{
		strings.NewReader(strings.Repeat(`
			Lorem Ipsum is simply dummy text of the printing and typesetting industry.
			Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
			when an unknown printer took a galley of type and scrambled it to make a type specimen book.
			It has survived not only five centuries, but also the leap into electronic typesetting,
			remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages,
			and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
			`, 1)),
	}, nil
}

type DummyReadCloser struct {
	io.Reader
}

func (s *DummyReadCloser) Close() error {
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
