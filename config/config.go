package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger LoggerConfig `toml:"logger"`
	OpenAI OpenAIConfig `toml:"openai"`
	Server ServerConfig `toml:"server"`
}

type LoggerConfig struct {
	Level string `toml:"level"`
	Concise bool `toml:"concise"`
}

type OpenAIConfig struct {
	ApiKey string `toml:"api_key"`
	MaxTokens int `toml:"max_tokens"`
	SystemPrompt string `toml:"system_prompt"`
}

type ServerConfig struct {
	Port int `toml:"port"`
}

func NewFromFile(file string, cfg *Config) error {
	if _, err := toml.DecodeFile(file, cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}
	return nil
}
