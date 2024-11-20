package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// type Model interface {
// 	GetModelResponse(diffInfo string) string
// }

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Model struct {
	url       string
	modelName string
	key       string
}

func (m *Model) prepareRequestBody(diffInfo string) RequestBody {
	return RequestBody{
		Model: m.modelName,
		Messages: []Message{
			{Role: "system", Content: GetSystemPrompt()},
			{Role: "user", Content: prepareQuestionContent(diffInfo)},
		},
	}
}

func prepareQuestionContent(diffInfo string) string {
	return "Please write a commit message for these git changes, " + diffInfo
}

func CreateModel(url, modelName, key string) *Model {
	return &Model{url, modelName, key}
}

func (m *Model) PrePareRequest(diffInfo string) (*http.Request, error) {
	reqBody := m.prepareRequestBody(diffInfo)

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, m.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
