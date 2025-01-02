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

const ARK_BASE_URL = "https://ark.cn-beijing.volces.com/api/v3"

type DoubaoModel struct {
	BaseModel
}

type DoubaoSessionReqBody struct {
	Model    string    `json:"model"`
	Mode     string    `json:"mode"`
	Messages []Message `json:"messages"`
	TTL      int       `json:"ttl"`
	// TruncationStrategy truncationStrategy `json:"truncation_strategy"`
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

func (m *DoubaoModel) prepareSessionCreateReqBody() DoubaoSessionReqBody {
	return DoubaoSessionReqBody{
		Model: m.modelName,
		Messages: []Message{
			{Role: "system", Content: GetSystemPrompt()},
		},
		Mode: "session",
		TTL:  3600,
	}
}

func (m *DoubaoModel) createSessionRequest() (*http.Request, error) {
	reqBody := m.prepareSessionCreateReqBody()

	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Failed to marshal request body:", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ARK_BASE_URL+"/context/create",
		bytes.NewBuffer(requestBodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (m *DoubaoModel) handleSessionCreateResponse(req *http.Request) (string, error) {
	client := &http.Client{
		Timeout: 100 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ERROR(handleSessionCreateResponse): Failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(handleSessionCreateResponse): Failed to read response body: %w",
			err,
		)
	}

	if resp.StatusCode != http.StatusOK {
		var respErr SessionError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return "", fmt.Errorf(
				"ERROR(handleSessionCreateResponse): Failed to unmarshal response: %w",
				err,
			)
		}
		return "", fmt.Errorf(respErr.Error.Message)
	} else {
		var respSession SessionResp
		err = json.Unmarshal(body, &respSession)
		if err != nil {
			return "", fmt.Errorf(
				"ERROR(handleSessionCreateResponse): Failed to unmarshal response: %w",
				err,
			)
		}
		return respSession.ID, nil
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
	req, err := m.createSessionRequest()
	if err != nil {
		return "", err
	}

	sessionID, err := m.handleSessionCreateResponse(req)

	return sessionID, err
}
