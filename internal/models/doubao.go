package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DoubaoModel struct {
	BaseModel
}

type DoubaoSessionReqBody struct {
	Model    string    `json:"model"`
	Mode     string    `json:"mode"`
	Messages []Message `json:"messages"`
	// TruncationStrategy truncationStrategy `json:"truncation_strategy"`
	TTL int `json:"ttl"`
}

type truncationStrategy struct {
	Type              string `json:"type"`
	LastHistoryTokens int    `json:"last_history_tokens"`
}
type SessionResp struct {
	ID                 string `json:"id"`
	Model              string `json:"model"`
	Mode               string `json:"mode"`
	TruncationStrategy struct {
		Type             string `json:"type"`
		LastHistoryToken int    `json:"last_history_token"`
	} `json:"truncation_strategy"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	TTL int `json:"ttl"`
}

type SessionError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
		Type    string `json:"type"`
	} `json:"error"`
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

// TODO: refactor Doubao LLMs communication with volcano SDK

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

func (m *DoubaoModel) CreateContextSession() (string, error) {
	reqBody := DoubaoSessionReqBody{
		Model: "ep-20241228235334-mrc9g",
		Messages: []Message{
			{Role: "system", Content: "你是李雷，你只会说“我是李雷”"},
		},
		Mode: "session",
		TTL:  3600,
		// TruncationStrategy: truncationStrategy{
		// 	Type:              "last_history_tokens",
		// 	LastHistoryTokens: 8192,
		// },
	}
	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Failed to marshal request body:", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://ark.cn-beijing.volces.com/api/v3/context/create",
		bytes.NewBuffer(requestBodyBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 100 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Failed to read response body:", err)
		}

		var response SessionError
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal("Failed to unmarshal response:", err)
		}

		fmt.Print(response)
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response body:", err)
	}

	var response SessionResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Failed to unmarshal response:", err)
	}

	fmt.Print(response)

	return response.ID, nil
}
