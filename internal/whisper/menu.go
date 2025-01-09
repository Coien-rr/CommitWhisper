package whisper

import (
	"github.com/Coien-rr/CommitWhisper/internal/models"
	"github.com/charmbracelet/huh"
)

func configMenu(config *Config) {
	// aiProvider, modelName, apiKey := config.AiProvider, config.ModelName, config.APIKey
	var confirm bool
	var endpoint string

	if config.AiProvider == "Doubao" {
		endpoint = config.ModelName
	}

	endpointInput := huh.NewInput().Title("Enter Your Endpoint").Value(&endpoint)

	modelSelector := huh.NewSelect[string]().
		Value(&config.ModelName).
		Height(10).
		Title("Choose Your Model")

	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&config.AiProvider).
				Title("Choose Your AiProvider").
				Height(10).OptionsFunc(func() []huh.Option[string] {
				return huh.NewOptions(models.AiProviderList...)
			}, nil),
		),

		huh.NewGroup(
			modelSelector,
		).WithHideFunc(func() bool {
			if config.AiProvider == "Doubao" {
				return true
			} else {
				modelOptions := huh.NewOptions(models.ModelsList[config.AiProvider]...)
				modelSelector.Options(modelOptions...)
				return false
			}
		}),

		huh.NewGroup(
			endpointInput,
		).WithHideFunc(func() bool {
			return config.AiProvider != "Doubao"
		}),

		huh.NewGroup(
			huh.NewInput().Title("Enter Your API Key").Value(&config.APIKey),
			huh.NewConfirm().Title("Confirm Config?").Value(&confirm),
		),
	).Run()

	if config.AiProvider == "Doubao" {
		config.ModelName = endpoint
	}
	config.APIUrl = models.AiProviderBaseUrlList[config.AiProvider]
	// TODO: add config check
}
