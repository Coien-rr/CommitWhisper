package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Coien-rr/CommitWhisper/internal/interaction"
)

type Model interface {
	GetModelResponse(diffInfo string) string
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type BaseModel struct {
	url       string
	modelName string
	key       string
}

func PrepareRequestBody(model, diffInfo string) RequestBody {
	return RequestBody{
		Model: model,
		Messages: []Message{
			{Role: "system", Content: GetSystemPrompt()},
			{Role: "user", Content: PrepareQuestionContent(diffInfo)},
		},
	}
}

func PrepareQuestionContent(diffInfo string) string {
	return "Please write a commit message for these git changes, " + diffInfo
}

func CreateModel(url, modelName, key string) *BaseModel {
	return &BaseModel{url, modelName, key}
}

func (model *BaseModel) GenerateCommitMessage(diffInfo string) (string, error) {
	reqBody := PrepareRequestBody(model.modelName, diffInfo)

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, model.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+model.key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := interaction.GeneratingCommitMessage(req)
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

func (model *BaseModel) HandleGeneratedCommitMsg(diffInfo string) {
	for {
		commitMsg, _ := model.GenerateCommitMessage(diffInfo)
		switch interaction.ConformGeneratedMessage(commitMsg) {
		case true:
			fmt.Println(commitMsg)
			return
		case false:
			fmt.Printf("GenerateCommitMessage: %v\nNot Good Enough, Retry!\n", commitMsg)
			continue
		}
	}
}
