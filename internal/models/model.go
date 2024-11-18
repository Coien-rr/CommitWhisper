package models

type Model interface {
	GetModelResponse(diffInfo string) string
}
