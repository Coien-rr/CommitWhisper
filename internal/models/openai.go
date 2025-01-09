package models

func NewOpenAIModelAgent(modelName, baseUrl, apiKey string) (Model, error) {
	return NewGenericModelAgent(modelName, baseUrl, apiKey)
}
