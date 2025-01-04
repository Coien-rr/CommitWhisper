package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Coien-rr/CommitWhisper/internal/comm"
	selfErr "github.com/Coien-rr/CommitWhisper/pkg/errors"
)

const ARK_BASE_URL = "https://ark.cn-beijing.volces.com/api/v3"

type DoubaoModel struct {
	BaseModel
	ContextID string
	chatCount int
	isPrompt  bool
}

// TODO: refactor base model
type doubaoSessionReqBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type DoubaoCreateSessionReqBody struct {
	Mode string `json:"mode"`
	doubaoSessionReqBody
	TTL int `json:"ttl"`
}

type DoubaoSessionChatReqBody struct {
	ContextID string `json:"context_id"`
	doubaoSessionReqBody
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

type ResponseBody struct {
	Choices []Choices `json:"choices"`
	// Choices []struct {
	// 	Message struct {
	// 		Role    string `json:"role"`
	// 		Content string `json:"content"`
	// 	} `json:"message"`
	// } `json:"choices"`
}

type Choices struct {
	Message Message `json:"message"`
	// Logprobs     interface{} `json:"logprobs"`
	// FinishReason string      `json:"finish_reason"`
	// Index        int         `json:"index"`
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

func (m *DoubaoModel) setContextID(cxtID string) {
	m.ContextID = cxtID
}

// func (m *DoubaoModel) initContextSession() error {
// }

func (m *DoubaoModel) prepareRequestBody(diffInfo string) RequestBody {
	return RequestBody{
		Model: m.modelName,
		Messages: []Message{
			{Role: "system", Content: GetSystemPrompt()},
			{Role: "user", Content: prepareQuestionContent(diffInfo)},
		},
	}
}

func (m *DoubaoModel) prepareSessionCreateReqBody() DoubaoCreateSessionReqBody {
	return DoubaoCreateSessionReqBody{
		doubaoSessionReqBody: doubaoSessionReqBody{
			Model: m.modelName,
			Messages: []Message{
				{Role: "system", Content: GetSystemPrompt()},
			},
		},
		Mode: "session",
		TTL:  3600,
	}
}

func (m *DoubaoModel) prepareSessionChatReqBody(diffInfo string) DoubaoSessionChatReqBody {
	var prompt string
	if m.chatCount == 0 {
		prompt = getCommitGeneratePrompt(diffInfo)
	} else {
		prompt = getRefinePrompt(diffInfo)
	}
	m.chatCount++
	return DoubaoSessionChatReqBody{
		ContextID: m.ContextID,
		doubaoSessionReqBody: doubaoSessionReqBody{
			Model: m.modelName,
			Messages: []Message{
				{Role: "user", Content: prompt},
			},
		},
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

func (m *DoubaoModel) CreateSessionChatRequest(diffInfo string) (*http.Request, error) {
	reqBody := m.prepareSessionChatReqBody(diffInfo)

	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Failed to marshal request body:", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ARK_BASE_URL+"/context/chat/completions",
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

func (m *DoubaoModel) handleSessionCreateResponse_new(
	respBody []byte,
	statusCode int,
) (string, error) {
	if statusCode != http.StatusOK {
		var respErr SessionError
		err := json.Unmarshal(respBody, &respErr)
		if err != nil {
			return "", fmt.Errorf(
				"ERROR(handleSessionCreateResponse): Failed to unmarshal response: %w",
				err,
			)
		}
		return "", fmt.Errorf(respErr.Error.Message)
	} else {
		var respSession SessionResp
		err := json.Unmarshal(respBody, &respSession)
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

	return comm.CreateLLMsRequest(reqBytes, m.key, m.url)
}

func (m *DoubaoModel) CreateContextSession() (string, error) {
	req, err := m.createSessionRequest()
	if err != nil {
		return "", err
	}

	sessionID, err := m.handleSessionCreateResponse(req)

	m.setContextID(sessionID)

	return sessionID, err
}

func (m *DoubaoModel) GenerateCommitMessage(diffInfo string) (string, error) {
	client := comm.NewLLMsServiceClient(m.key, ARK_BASE_URL)

	if m.ContextID == "" {
		// TODO: refactor
		requestBodyBytes, err := json.Marshal(m.prepareSessionCreateReqBody())
		if err != nil {
			return "", fmt.Errorf(
				"ERROR(GenerateCommitMessage): Failed to marshal request body: %w",
				err,
			)
		}

		resp, statusCode, err := client.CreateLLMsContextSession(requestBodyBytes)

		contextID, err := m.handleSessionCreateResponse_new(resp, statusCode)
		if err != nil {
			return "", fmt.Errorf(
				"ERROR(GenerateCommitMessage): %w",
				err,
			)
		}

		m.setContextID(contextID)
	}

	requestBody, err := json.Marshal(m.prepareSessionChatReqBody(diffInfo))
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(GenerateCommitMessage): Failed to marshal request body: %w",
			err,
		)
	}

	resp, statusCode, err := client.CreateChatReqWithLLMs(requestBody)
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(GenerateCommitMessage): %w",
			err,
		)
	}

	if statusCode == http.StatusOK {
		var response ResponseBody
		if err := json.Unmarshal(resp, &response); err != nil {
			return "", fmt.Errorf(
				"ERROR(generateCommitMessage): failed to parse response JSON: %v",
				err,
			)
		}
		return response.Choices[0].Message.Content, nil
	} else {
		var response errResponseBody
		if err := json.Unmarshal(resp, &response); err != nil {
			return "", fmt.Errorf(
				"ERROR(generateCommitMessage): failed to parse response JSON: %v",
				err,
			)
		}
		switch statusCode {
		case http.StatusUnauthorized:
			// TODO: key Invalid  error
			return "", selfErr.NewInvalidKeyError(response.ErrorMsg.Message)

		case http.StatusNotFound:
			// TODO: model not found error
			return "", selfErr.NewNotFoundError(response.ErrorMsg.Message)

		case http.StatusTooManyRequests:
			// TODO: rate error or bill error
			return "", selfErr.NewTooManyReqError(response.ErrorMsg.Message)

		default:
			return "", nil
		}
	}

	// req, err := w.llmModel.CreateSessionChatRequest(diffInfo)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to prepare request: %v", err)
	// }
	// resp, err := w.generatingCommitMessage(req)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to send request: %v", err)
	// }
	//
	// defer resp.Body.Close()
	//
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("ERROR(generateCommitMessage): failed to read response body: %v", err)
	// }
	//
}
