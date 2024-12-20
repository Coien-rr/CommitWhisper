package whisper

import (
	"github.com/Coien-rr/CommitWhisper/internal/models"
	"github.com/charmbracelet/huh"
)

func showMenu() Config {
	var aiProvider, modelName, apiKey, endpoint string
	var confirm bool

	modelSelector := huh.NewSelect[string]().
		Value(&modelName).
		Height(8).
		Title("Choose Your Model")

	endpointInput := huh.NewInput().Title("Enter Your Endpoint").Value(&endpoint)

	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&aiProvider).
				Title("Choose Your AiProvider").
				Height(5).
				OptionsFunc(func() []huh.Option[string] {
					return huh.NewOptions(models.AiProviderList...)
				}, nil),
		),

		huh.NewGroup(
			modelSelector,
		).WithHideFunc(func() bool {
			if aiProvider == "Doubao" {
				return true
			} else {
				modelOptions := huh.NewOptions(models.ModelsList[aiProvider]...)
				modelSelector.Options(modelOptions...)
				return false
			}
		}),

		huh.NewGroup(
			endpointInput,
		).WithHideFunc(func() bool {
			return aiProvider != "Doubao"
		}),

		huh.NewGroup(
			huh.NewInput().Title("Enter Your API Key").Value(&apiKey),
			huh.NewConfirm().Title("Confirm Config?").Value(&confirm),
		),
	).Run()

	if aiProvider == "Doubao" {
		modelName = endpoint
	}

	return Config{
		AiProvider: aiProvider,
		ModelName:  modelName,
		APIUrl:     models.ModelsURLList[aiProvider],
		APIKey:     apiKey,
	}
}

func reconfigMenu(config *Config) {
	// aiProvider, modelName, apiKey := config.AiProvider, config.ModelName, config.APIKey
	var confirm bool
	var endpoint string

	endpointInput := huh.NewInput().Title("Enter Your Endpoint").Value(&endpoint)

	modelSelector := huh.NewSelect[string]().
		Value(&config.ModelName).
		Height(8).
		Title("Choose Your Model")

	huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&config.AiProvider).
				Title("Choose Your AiProvider").
				Height(5).OptionsFunc(func() []huh.Option[string] {
				return huh.NewOptions(models.AiProviderList...)
			}, nil),
		),

		huh.NewGroup(
			modelSelector,
		).WithHideFunc(func() bool {
			// TODO: Refactor
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
	config.APIUrl = models.ModelsURLList[config.AiProvider]
	// TODO: add config check
}
