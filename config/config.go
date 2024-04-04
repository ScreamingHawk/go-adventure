package config

import (
	"fmt"
	"os"
	"strconv"

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

func NewFromEnv(cfg *Config) error {
	// Logger
	if env, ok := os.LookupEnv("LOGGER_LEVEL"); ok {
		cfg.Logger.Level = env
	}
	if env, ok := os.LookupEnv("LOGGER_CONCISE"); ok {
		cfg.Logger.Concise = env == "true"
	}
	// OpenAI
	if env, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		cfg.OpenAI.ApiKey = env
	} else {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}
	if env, ok := os.LookupEnv("OPENAI_MAX_TOKENS"); ok {
		maxTokens, err := strconv.Atoi(env)
		if err != nil {
			return fmt.Errorf("failed to convert OPENAI_MAX_TOKENS to integer: %w", err)
		}
		cfg.OpenAI.MaxTokens = maxTokens
	}
	if env, ok := os.LookupEnv("OPENAI_SYSTEM_PROMPT"); ok {
		cfg.OpenAI.SystemPrompt = env
	} else {
		return fmt.Errorf("OPENAI_SYSTEM_PROMPT is required")
	}
	// Server
	if env, ok := os.LookupEnv("SERVER_PORT"); ok {
		port, err := strconv.Atoi(env)
		if err != nil {
			return fmt.Errorf("failed to convert SERVER_PORT to integer: %w", err)
		}
		cfg.Server.Port = port
	} else {
		cfg.Server.Port = 0
	}
	return nil
}
