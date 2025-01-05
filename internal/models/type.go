package models

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
