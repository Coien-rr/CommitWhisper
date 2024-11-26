package models

const MyAPIKEY = "sk-502bbd61da624094920ce1c00375f45f"

var ModelsList = map[string][]string{
	"Qwen": {
		"qwen-coder-plus",
		"qwen-coder-plus-latest",
		"qwen-coder-plus-2024-11-06",
		"qwen-coder-turbo",
		"qwen-coder-turbo-latest",
		"qwen-coder-turbo-2024-09-19",
		"qwen2.5-coder-32b-instruct",
		"qwen2.5-coder-14b-instruct",
		"qwen2.5-coder-7b-instruct",
		"qwen2.5-coder-3b-instruct",
		"qwen2.5-coder-1.5b-instruct",
		"qwen2.5-coder-0.5b-instruct",
	},
	"Doubao": {
		"Doubao-lite-4k",
		"Doubao-lite-32k",
		"Doubao-pro-4k",
		"Doubao-pro-32k",
		"Doubao-pro-128k",
	},
	"Skylark": {
		"Skylark2-pro-character-4k",
		"Skylark2-pro-32k",
		"Skylark2-pro-4k",
		"Skylark2-pro-turbo-8k",
		"Skylark2-lite-8k",
	},

	// TODO: OpenAI
	// TODO: Claude
	// TODO: Gemini
}

var ModelsURLList = map[string]string{
	"Qwen":    "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
	"Doubao":  "doubao TODO",
	"Skylark": "doubao TODO",
}
