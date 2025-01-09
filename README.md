# :robot: CommitWhisper _(Ssshh 💤 )_

_Whisper From `LLMs` For Your `Commit` Msg_

`CommitWhisper` utilizes `LLMs` generate git commit messages based on code diff information.

## ✨ Features

- 🔍 Get Diff Info From Staged Files
- 🌈 Interactive configuration selection for easy model switching.
- 🌅 Supports **Mainstream AI providers**, Such as:
  - [ChatGPT](https://platform.openai.com/docs/models) From `OpenAI`
  - 💤 [Claude](https://claude.ai/chats) From `Anthropic`
  - 💤 [Gemini](https://ai.google.dev/gemini-api/docs) From `Google`
  - 💤 More
- 🚀 Especially supports **Chinese AI providers**, Such as:
  - [Tongyi Qianwen(通义千问)](https://www.aliyun.com/product/bailian) From `Alibaba`
  - [Doubao(豆包)](https://www.volcengine.com/product/doubao) From `ByteDance`
  - [DeepSeek](https://www.deepseek.com/) From `DeepSeek`
  - 💤 More
- 💬 [TODO] Interactive prompts allow the model to refine generated messages.
- ⏳ Coming Soon More

### :jigsaw: Supported Ai Provider

| LLMs              | Refer To Get API KEY |
| ----------------- | -------------------- |
| Qianwen           | ✅ [Get Your Key](https://www.aliyun.com/product/bailian)|
| OpenAI            | ✅ [Get Your Key](https://platform.openai.com)|
| Doubao            | ✅ [Get Your Endpoint and Key](https://console.volcengine.com)|
| DeepSeek          | ✅ [Get Your Key](https://www.deepseek.com/)|
| Claude            | 💤 Coming Soon       |

## 🔥 Status
>
> [!WARNING]
> This CLI is _beta_ quality. Expect breaking changes and many bugs

## ⚡️ Requirements

- Nerd Font
- Git
- Go

## 📦 Installation

### Git Version（Under Dev）

```sh
go install github.com/Coien-rr/CommitWhisper@latest
ln -s $(go env GOPATH)/bin/CommitWhisper $(go env GOPATH)/bin/cw
cw # Using Commit Whisper in Git-Repo
```

### Release Version (Recommend)

1. DownLoad Released Package Based On Your System From [Release Page](https://github.com/Coien-rr/CommitWhisper/releases)
2. Extract the package to obtain the CLI `cw`

3. ```sh
   sudo ln -s ./cw /usr/local/bin/
   cw # Using Commit Whisper in Git-Repo
   ```

## ⚙️ Configuration

**Cofig File Locates** _`~/.commitwhisper`_

When you first start `cw`, it will prompt you interactively to select configurations, including

- AI provider,
- LLMs,
- API key,
and other related information.

Also, You Can Use `cw rc` to reconfig.

Demo Config:

```yaml
AiProvider: Qwen
ModelName: qwen2.5-coder-32b-instruct
APIUrl: https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions
APIKey: [Please Put Your API Key Here]
```

## :crystal_ball: RoadMap

_💡 Ideas to explore_

- [ ] Support More AiProvider, Such as OpenAI✅, Claude💤, etc.(_In Process_)
- [ ] Generate more detailed commit descriptions(multi-lines commit).
- [ ] Explain the purpose of code changes compare two commit.
- [ ] Interactive prompts allow the model to refine generated messages.
- [ ] Integrate `commitlint` to evaluate the quality of commit messages generated by LLMs.
- [ ] And More
