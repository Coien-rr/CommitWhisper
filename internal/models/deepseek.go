package models

func NewDeepSeekModelAgent(modelName, baseUrl, apiKey string) (Model, error) {
	return NewGenericModelAgent(modelName, baseUrl, apiKey)
}
