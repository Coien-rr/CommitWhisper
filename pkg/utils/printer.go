package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Printer struct {
	infoTheme    lipgloss.Style
	warningTheme lipgloss.Style
	errorTheme   lipgloss.Style
}

func NewPrinter() *Printer {
	return &Printer{
		infoTheme: lipgloss.NewStyle().
			Foreground(lipgloss.Color(DefaultTheme.Cyan)),
		warningTheme: lipgloss.NewStyle().Italic(true).Bold(true).
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
