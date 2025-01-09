package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	selfErr "github.com/Coien-rr/CommitWhisper/pkg/errors"
)

// TODO: add token usage record
type ResponseBody struct {
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Message Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type genericLLMsServiceReqBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type session struct {
	messages []Message
}

// TODO: role check
func (s *session) appendMessage(role, content string) {
	s.messages = append(s.messages, Message{role, content})
}

func (s *session) getMessages() []Message {
	return s.messages
}

func handleChatRespFromLLMs(resp []byte, statusCode int) (msg string, err error) {
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
