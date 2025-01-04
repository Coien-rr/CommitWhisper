package whisper

import (
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
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Message Message `json:"message"`
	// Logprobs     interface{} `json:"logprobs"`
	// FinishReason string      `json:"finish_reason"`
	// Index        int         `json:"index"`
}

type Message struct {
	// Refusal interface{} `json:"refusal"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type errResponseBody struct {
	ErrorMsg  errorMsg `json:"error"`
	RequestID string   `json:"request_id"`
}

type errorMsg struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   any    `json:"param"`
	Code    string `json:"code"`
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

func (w *Whisper) isUsePromptRefine() bool {
	return true
}

func (w *Whisper) Run() {
	w.greet()
	if w.checkIsGitRepo() {
		if w.isUsePromptRefine() {
			// TODO: 1. create context session
			// TODO: 2. generateRefinedCommitMsg()
			w.createCommitGeneratorBySession()
		} else {
			w.generateAICommitByGitDiff()
		}
	}
}

func (w *Whisper) generateCommitMessage(diffInfo string) (string, error) {
	return w.generatingCommitMessage(diffInfo)
}

func (w *Whisper) generatingCommitMessage(prompt string) (string, error) {
	var commitMsg string
	var err error

	action := func() {
		commitMsg, err = w.llmModel.GenerateCommitMessage(prompt)
	}
	_ = spinner.New().
		Title("Generating Commit Message󰒲 ").
		TitleStyle(lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#02BF87"))).
		Action(action).
		Run()

	return commitMsg, err
}

func (w *Whisper) conformGeneratedMessage(generatedCommitMsg string) bool {
	var confirm bool

	huh.NewConfirm().
		Title("Confirm the commit message?").
		Description(generatedCommitMsg).
		Affirmative("Confirm!").
		Negative("Refine!").
		Value(&confirm).Run()

	return confirm
}

func (w *Whisper) refineGeneratedMessage() string {
	var refinePrompt string

	huh.NewInput().
		Title("Input Your Prompt For Refine").
		Prompt("💡").
		// Validate(isFood).
		Value(&refinePrompt).Run()

	return refinePrompt
}

func (w *Whisper) handleGeneratedCommitMsg(diffInfo string) {
	refinePrompt := diffInfo
	for {
		commitMsg, err := w.generateCommitMessage(refinePrompt)
		if err != nil {
			// TODO: add error handle
			utils.WhisperPrinter.Error(err.Error())
			return
		}
		utils.WhisperPrinter.Info("Commit Message Generated!")
		utils.WhisperPrinter.Info("GenerateCommitMessage: " + commitMsg)
		switch w.conformGeneratedMessage(commitMsg) {
		case true:
			copyToClipboard(commitMsg)
			utils.WhisperPrinter.Info("Copied Commit Message into ClipBoard ✔")
			return
		case false:
			refinePrompt = w.refineGeneratedMessage()
			// utils.WhisperPrinter.Info(refinePrompt)
			continue
		}
	}
}

// TODO: refactor
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

func (w *Whisper) createCommitGeneratorBySession() {
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
