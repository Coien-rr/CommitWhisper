package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DoubaoModel struct {
	BaseModel
}

func (m *DoubaoModel) prepareRequestBody(diffInfo string) RequestBody {
	return RequestBody{
		Model: m.modelName,
		Messages: []Message{
			{Role: "system", Content: GetSystemPrompt()},
			{Role: "user", Content: prepareQuestionContent(diffInfo)},
		},
	}
}

func (m *DoubaoModel) PrepareRequest(diffInfo string) (*http.Request, error) {
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
