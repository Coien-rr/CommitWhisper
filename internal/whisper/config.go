package whisper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Coien-rr/CommitWhisper/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AiProvider string `yaml:"AiProvider"`
	ModelName  string `yaml:"ModelName"`
	APIUrl     string `yaml:"APIUrl"`
	APIKey     string `yaml:"APIKey"`
}

const configFileName = ".commitwhisper"

var defaultConfigPath = filepath.Join(os.Getenv("HOME"), configFileName)

func isConfigFileExist() bool {
	_, err := os.Stat(defaultConfigPath)
	return !os.IsNotExist(err)
}

func (config *Config) String() string {
	return fmt.Sprintf(
		"Config{\n   AiProvider: %s,\n   ModelName: %s,\n   APIURL: %s,\n   APIKey: %s,\n}",
		config.AiProvider,
		config.ModelName,
		config.APIUrl,
		config.APIKey,
	)
}

func getEnvConfigFromDotFile() (Config, error) {
	data, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, err
}

func setEnvConfigToDotFile(config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return nil
	}

	return os.WriteFile(defaultConfigPath, data, 0o644)
}

func GetConfig() Config {
	if isConfigFileExist() {
		config, err := getEnvConfigFromDotFile()
		if err != nil {
			println(err)
		}
		return config
	} else {
		var config Config
		for {
			config = showMenu()
			if err := config.checkConfig(); err != nil {
				utils.WhisperPrinter.Error(err.Error())
				utils.WhisperPrinter.Info("Please Config Whisper Again  ")
			} else {
				break
			}
		}
		// TODO: Set to goroutine
		setEnvConfigToDotFile(config)
		return config
	}
}

func (config *Config) checkConfig() error {
	if config.APIUrl == "" || config.APIKey == "" {
		return errors.New("ConfigError: API Url or Key can't be empty")
	}
	return nil
}

func ReConfig() {
	if isConfigFileExist() {
		config, err := getEnvConfigFromDotFile()
		if err != nil {
			utils.WhisperPrinter.Error(err.Error())
			return
		}
		utils.WhisperPrinter.Info(fmt.Sprintf("The current config:\n%s", &config))
		var newConfig Config
		for {
			newConfig = reconfigMenu(config)
			if err := newConfig.checkConfig(); err != nil {
				utils.WhisperPrinter.Error(err.Error())
				utils.WhisperPrinter.Info("Please Config Whisper Again  ")
			} else {
				break
			}
		}
		// TODO: Set to goroutine
		setEnvConfigToDotFile(newConfig)
	} else {
		utils.WhisperPrinter.Warning("The configuration file does not exist yet! ")
		utils.WhisperPrinter.Info("Please start commitwhisper and follow the prompts to configure")
	}
}

func ShowConfig() {
	if isConfigFileExist() {
		config, err := getEnvConfigFromDotFile()
		if err != nil {
			utils.WhisperPrinter.Error(err.Error())
			return
		}
		utils.WhisperPrinter.Info(fmt.Sprintf("The current config:\n%s", &config))
	} else {
		utils.WhisperPrinter.Warning("The configuration file does not exist yet! ")
		utils.WhisperPrinter.Info("Please start commitwhisper and follow the prompts to configure")
	}
}
