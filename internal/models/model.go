package models

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
