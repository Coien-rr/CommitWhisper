package models

var AiProviderList = []string{
	"OpenAI",
	"Qwen",
	"Doubao",
	"DeepSeek",
}

var ModelsList = map[string][]string{
	"OpenAI": {
		"gpt-4o-mini",
		"gpt-4o-latest",
		"o1-preview",
		"o1-mini",
		"gpt-4-turbo",
		"gpt-4",
	},
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
	"Doubao": {},
	"DeepSeek": {
		"deepseek-chat",
	},

	// TODO: Claude
	// TODO: Gemini
}

var AiProviderBaseUrlList = map[string]string{
	"Qwen":     "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
	"OpenAI":   "https://api.openai.com/v1/chat/completions",
	"Doubao":   "https://ark.cn-beijing.volces.com/api/v3",
	"DeepSeek": "https://api.deepseek.com/v1/chat/completions",
}
