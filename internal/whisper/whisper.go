package whisper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Coien-rr/CommitWhisper/internal/git"
	"github.com/Coien-rr/CommitWhisper/internal/models"
	"github.com/Coien-rr/CommitWhisper/pkg/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

type Whisper struct {
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

func NewWhisper(config Config) *Whisper {
	if err := config.checkConfig(); err != nil {
		utils.WhisperPrinter.Error(err.Error())
		utils.WhisperPrinter.Info("You should reconfig it")
		return nil
	}

	if engine, err := models.CreateModel(config.AiProvider, config.ModelName, config.APIUrl, config.APIKey); err == nil {
		return &Whisper{
			llmModel: engine,
		}
	} else {
		utils.WhisperPrinter.Warning(err.Error())
		return nil
	}
}

func (w *Whisper) greet() {
	utils.WhisperPrinter.Info("Hi, This is CommitWhisper🎉")
}

func (w *Whisper) checkIsGitRepo() bool {
	if !git.IsGitRepo() {
		utils.WhisperPrinter.Warning("The current workspace is not a valid Git repository 👻")
		utils.WhisperPrinter.Info("Please use 'cw' within a Git repository 🤙")
		return false
	}
	return true
}

func (w *Whisper) Run() {
	w.greet()
	if w.checkIsGitRepo() {
		w.generateAICommitByGitDiff()
	}
}

func (w *Whisper) generateCommitMessage(diffInfo string) (string, error) {
	req, err := w.llmModel.PrepareRequest(diffInfo)
	resp, err := w.generatingCommitMessage(req)
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

func (w *Whisper) generatingCommitMessage(req *http.Request) (*http.Response, error) {
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

	utils.WhisperPrinter.Info("Commit Message Generated!")

	return res, err
}

func (w *Whisper) conformGeneratedMessage(generatedCommitMsg string) bool {
	var confirm bool

	huh.NewConfirm().
		Title("Confirm the commit message?").
		Description(generatedCommitMsg).
		Affirmative("Confirm!").
		Negative("Retry!").
		Value(&confirm).Run()

	return confirm
}

func (w *Whisper) handleGeneratedCommitMsg(diffInfo string) {
	for {
		commitMsg, _ := w.generateCommitMessage(diffInfo)
		utils.WhisperPrinter.Info("GenerateCommitMessage: " + commitMsg)
		switch w.conformGeneratedMessage(commitMsg) {
		case true:
			copyToClipboard(commitMsg)
			utils.WhisperPrinter.Info("Copied Commit Message into ClipBoard ✔")
			return
		case false:
			utils.WhisperPrinter.Warning("Not Good Enough, Retry!")
			continue
		}
	}
}

func (w *Whisper) generateAICommitByGitDiff() {
	diff, err := git.GetGitDiff()
	if err != nil {
		utils.WhisperPrinter.Error(err.Error())
		return
	} else if diff == "" {
		utils.WhisperPrinter.Info("Working tree clean,Nothing to commit ")
		return
	}
	w.handleGeneratedCommitMsg(diff)
}
