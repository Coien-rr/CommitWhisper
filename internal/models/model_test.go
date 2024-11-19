package models

import (
	"reflect"
	"testing"
)

func TestPrepareRequestBody(t *testing.T) {
	got := PrepareRequestBody("qwen2.5-coder-3b-instruct", "Hello")

	want := RequestBody{
		"qwen2.5-coder-3b-instruct",
		[]Message{
			{"system", GetSystemPrompt()},
			{"user", PrepareQuestionContent("Hello")},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, but got %v", want, got)
	}
}
