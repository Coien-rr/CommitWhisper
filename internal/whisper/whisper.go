package whisper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Coien-rr/CommitWhisper/internal/models"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Whisper struct {
	printer  Printer
	llmModel models.Model
}

type ResponseBody struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewWhisper(url, modelName, key string) *Whisper {
	return &Whisper{
		printer:  *NewWhisperPrinter(),
		llmModel: models.CreateModel("qwen", url, modelName, key),
	}
}

func (w *Whisper) Greet() {
	w.printer.Info("Hi, This is CommitWhisper🎉")
}

func (w *Whisper) GenerateCommitMessage(diffInfo string) (string, error) {
	req, err := w.llmModel.PrepareRequest(diffInfo)
	resp, err := w.GeneratingCommitMessage(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return response.Choices[0].Message.Content, nil
}

func (w *Whisper) GeneratingCommitMessage(req *http.Request) (*http.Response, error) {
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

	w.printer.Info("Commit Message Generated!")

	return res, err
}

func (w *Whisper) ConformGeneratedMessage(generatedCommitMsg string) bool {
	var confirm bool

	huh.NewConfirm().
		Title("Confirm the commit message?").
		Description(generatedCommitMsg).
		Affirmative("Confirm!").
		Negative("Retry!").
		Value(&confirm).Run()

	return confirm
}

func (w *Whisper) HandleGeneratedCommitMsg(diffInfo string) {
	for {
		commitMsg, _ := w.GenerateCommitMessage(diffInfo)
		w.printer.Info("GenerateCommitMessage: " + commitMsg)
		switch w.ConformGeneratedMessage(commitMsg) {
		case true:
			copyToClipboard(commitMsg)
			w.printer.Info("Copied Commit Message into ClipBoard💯")
			return
		case false:
			w.printer.Warning("Not Good Enough, Retry!")
			continue
		}
	}
}
