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

var WhisperPrinter = NewPrinter()

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
	fmt.Println(p.errorTheme.Render(fmt.Sprintf(" RuntimeError: %v", msg)))
}

func (p *Printer) WarningDisplayLists(msgTitle string, list []string) {
	p.Warning(msgTitle)
	for _, item := range list {
		if item == "" {
			continue
		}
		fmt.Println(p.warningTheme.Render(fmt.Sprintf("   %v", item)))
	}
}

func (p *Printer) InfoDisplayLists(msgTitle string, list []string) {
	p.Info(msgTitle)
	for _, item := range list {
		if item == "" {
			continue
		}
		fmt.Println(p.infoTheme.Render(fmt.Sprintf("   %v", item)))
	}
}
