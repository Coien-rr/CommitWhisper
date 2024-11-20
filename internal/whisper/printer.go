package whisper

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Printer struct {
	infoTheme    lipgloss.Style
	warningTheme lipgloss.Style
	errorTheme   lipgloss.Style
}

func NewWhisperPrinter() *Printer {
	return &Printer{
		infoTheme: lipgloss.NewStyle().Italic(true).
			Foreground(lipgloss.Color(DefaultTheme.Cyan)),
		warningTheme: lipgloss.NewStyle().Italic(true).
			Foreground(lipgloss.Color(DefaultTheme.Orange)),
		errorTheme: lipgloss.NewStyle().Italic(true).
			Foreground(lipgloss.Color(DefaultTheme.Red)),
	}
}

func (p *Printer) Info(msg string) {
	fmt.Println(p.infoTheme.Render(fmt.Sprintf("  %v", msg)))
}

func (p *Printer) Warning(msg string) {
	fmt.Println(p.warningTheme.Render(fmt.Sprintf("  %v", msg)))
}

func (p *Printer) Error(msg string) {
	fmt.Println(p.errorTheme.Render(fmt.Sprintf("  %v", msg)))
}

// TODO: refactor it into Printer Object
func Info(msg string) {
	fmt.Println(GetInfoTheme().Render(fmt.Sprintf("  %v", msg)))
}

func Warning(msg string) {
	fmt.Println(GetWarningTheme().Render(fmt.Sprintf("  %v", msg)))
}

func Error(msg string) {
	fmt.Println(GetErrorTheme().Render(fmt.Sprintf("  %v", msg)))
}
