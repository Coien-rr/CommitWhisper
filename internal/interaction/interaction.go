package interaction

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

func Greet() {
	// TODO: refactor using lipgloss style
	color.HiCyan("Hi, This is CommitWhisper🎉")
}

func GeneratingCommitMessage(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var res *http.Response
	var err error

	action := func() {
		res, err = client.Do(req)
	}
	_ = spinner.New().
		Title("Generating Commit Message󰒲 ").
		TitleStyle(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#02BF87"))).
		Action(action).
		Run()

	style := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"})

	fmt.Println(style.Render(" Commit Message Generated!"))

	return res, err
}

func ConformGeneratedMessage(generatedCommitMsg string) bool {
	var confirm bool

	huh.NewConfirm().
		Title("Confirm the commit message?").
		Description(generatedCommitMsg).
		Affirmative("Confirm!").
		Negative("Retry!").
		Value(&confirm).Run()

	return confirm
}
