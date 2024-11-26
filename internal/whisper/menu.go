package whisper

import (
	"time"

	"github.com/Coien-rr/CommitWhisper/internal/models"
	"github.com/charmbracelet/huh"
)

func showMenu() Config {
	var aiProvider, modelName, apiKey string
	var confirm bool
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Options(huh.NewOptions(
				"Qwen", "Doubao", "Skylark")...).
				Value(&aiProvider).
				Title("Choose Your AiProvider").
				Height(5),
			huh.NewSelect[string]().
				Value(&modelName).
				Height(8).
				Title("Choose Your Model").
				OptionsFunc(func() []huh.Option[string] {
					s := models.ModelsList[aiProvider]
					time.Sleep(500 * time.Millisecond)
					return huh.NewOptions(s...)
				}, &aiProvider),
			huh.NewInput().Title("Enter Your API Key").Value(&apiKey),
			huh.NewConfirm().Title("Confirm Config?").Value(&confirm),
		),
	).Run()
	return Config{
		AiProvider: aiProvider,
		ModelName:  modelName,
		APIUrl:     models.ModelsURLList[aiProvider],
		APIKey:     apiKey,
	}
}

func reconfigMenu(config Config) Config {
	aiProvider, modelName, apiKey := config.AiProvider, config.ModelName, config.APIKey
	var confirm bool
	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Options(huh.NewOptions(
				"Qwen", "Doubao", "Skylark")...).
				Value(&aiProvider).
				Title("Choose Your AiProvider").
				Height(5),
			huh.NewSelect[string]().
				Value(&modelName).
				Height(8).
				Title("Choose Your Model").
				OptionsFunc(func() []huh.Option[string] {
					s := models.ModelsList[aiProvider]
					time.Sleep(500 * time.Millisecond)
					return huh.NewOptions(s...)
				}, &aiProvider),
			huh.NewInput().Title("Enter Your API Key").Value(&apiKey),
			huh.NewConfirm().Title("Confirm Config?").Value(&confirm),
		),
	).Run()
	// TODO: add config check
	return Config{
		AiProvider: aiProvider,
		ModelName:  modelName,
		APIUrl:     models.ModelsURLList[aiProvider],
		APIKey:     apiKey,
	}
}
