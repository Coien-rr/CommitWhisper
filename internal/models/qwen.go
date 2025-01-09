package models

func NewQwenModelAgent(modelName, baseUrl, apiKey string) (Model, error) {
	return NewGenericModelAgent(modelName, baseUrl, apiKey)
}
