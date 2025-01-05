package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Coien-rr/CommitWhisper/internal/comm"
	selfErr "github.com/Coien-rr/CommitWhisper/pkg/errors"
)

const ARK_BASE_URL = "https://ark.cn-beijing.volces.com/api/v3"

type DoubaoModel struct {
	BaseModel
	contextID     string
	isRefineStage bool
}

type DoubaoCreateSessionReqBody struct {
	Mode string `json:"mode"`
	genericLLMsServiceReqBody
	TTL int `json:"ttl"`
}

type DoubaoSessionChatReqBody struct {
	ContextID string `json:"context_id"`
	genericLLMsServiceReqBody
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
}

type Choices struct {
	Message Message `json:"message"`
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

func NewDoubaoModelAgent(modelName, baseUrl, apiKey string) (*DoubaoModel, error) {
	model := &DoubaoModel{
		BaseModel:     BaseModel{modelName: modelName, url: baseUrl, key: apiKey},
		isRefineStage: false,
	}

	err := model.initContextSession()

	return model, err
}

func (m *DoubaoModel) setContextID(cxtID string) {
	m.contextID = cxtID
}

func (m *DoubaoModel) setRefineStage() {
	m.isRefineStage = true
}

func (m *DoubaoModel) initContextSession() error {
	client := comm.NewLLMsServiceClient(m.key, m.url)

	// TODO: refactor
	requestBodyBytes, err := json.Marshal(m.prepareSessionCreateReqBody())
	if err != nil {
		return fmt.Errorf(
			"ERROR(initContextSession): Failed to marshal request body: %w",
			err,
		)
	}

	resp, statusCode, err := client.CreateLLMsContextSession(requestBodyBytes)

	contextID, err := m.handleSessionCreateResponse(resp, statusCode)
	if err != nil {
		return fmt.Errorf(
			"ERROR(initContextSession): %w",
			err,
		)
	}

	m.setContextID(contextID)

	return nil
}

func (m *DoubaoModel) prepareSessionCreateReqBody() DoubaoCreateSessionReqBody {
	return DoubaoCreateSessionReqBody{
		genericLLMsServiceReqBody: genericLLMsServiceReqBody{
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
	if !m.isRefineStage {
		prompt = getCommitGeneratePrompt(diffInfo)
		m.setRefineStage()
	} else {
		prompt = getCommitRefinePrompt(diffInfo)
	}

	return DoubaoSessionChatReqBody{
		ContextID: m.contextID,
		genericLLMsServiceReqBody: genericLLMsServiceReqBody{
			Model: m.modelName,
			Messages: []Message{
				{Role: "user", Content: prompt},
			},
		},
	}
}

func (m *DoubaoModel) handleSessionCreateResponse(
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

func (m *DoubaoModel) GenerateCommitMessage(diffInfo string) (string, error) {
	client := comm.NewLLMsServiceClient(m.key, ARK_BASE_URL)

	requestBody, err := json.Marshal(m.prepareSessionChatReqBody(diffInfo))
	if err != nil {
		return "", fmt.Errorf(
			"ERROR(GenerateCommitMessage): Failed to marshal request body: %w",
			err,
		)
	}

	resp, statusCode, err := client.CreateSessionChatReqWithLLMs(requestBody)
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
}
