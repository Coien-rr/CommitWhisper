package utils

import (
	"sync"

	"github.com/charmbracelet/lipgloss"
)

type Palette struct {
	BackGround string
	ForeGound  string
	Comment    string
	Cyan       string
	Green      string
	Orange     string
	Pink       string
	Purple     string
	Red        string
	Yellow     string
}

// NOTE: Using Dracula Theme as DefaultTheme
// URL: https://draculatheme.com/contribute
var DefaultTheme = Palette{
	BackGround: "#282A36",
	ForeGound:  "#F8F8F2",
	Comment:    "#6272A4",
	Cyan:       "#8BE9FD",
	Green:      "#50FA7B",
	Orange:     "#FFB86C",
	Pink:       "#FF79C6",
	Purple:     "#BD93F9",
	Red:        "#FF5555",
	Yellow:     "#F1FA8C",
}

var (
	InfoTheme    lipgloss.Style
	WarningTheme lipgloss.Style
	ErrorTheme   lipgloss.Style
	once         sync.Once
)
