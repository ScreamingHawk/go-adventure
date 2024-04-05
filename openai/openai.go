package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/go-chi/httplog/v2"
	ttlcache "github.com/jellydator/ttlcache/v3"
	gopenai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	logger *httplog.Logger
	client *gopenai.Client

	maxTokens    int
	systemPrompt string

	cache *ttlcache.Cache[string, []gopenai.ChatCompletionMessage]
}

func NewOpenAI(cfg *config.OpenAIConfig, logger *httplog.Logger) (*OpenAI, error) {
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	client := gopenai.NewClient(cfg.ApiKey)

	systemPrompt := cfg.SystemPrompt + ". You MUST ALWAYS and ONLY respond with JSON."

	cache := ttlcache.New(
		ttlcache.WithTTL[string, []gopenai.ChatCompletionMessage](10 * time.Minute),
	)
	go cache.Start()

	return &OpenAI{
		logger:       logger,
		client:       client,
		maxTokens:    cfg.MaxTokens,
		systemPrompt: systemPrompt,
		cache:        cache,
	}, nil
}

func (o *OpenAI) CreateChat(ctx context.Context, key string, prompt string) (string, error) {
	if o.cache.Has(key) {
		return "", fmt.Errorf("chat with key %q already exists", key)
	}
	messages := []gopenai.ChatCompletionMessage{
		{Role: gopenai.ChatMessageRoleSystem, Content: o.systemPrompt},
		{Role: gopenai.ChatMessageRoleUser, Content: prompt},
	}

	response, err := o.callChatComplete(ctx, key, messages)
	if err != nil {
		return "", err
	}

	messages = append(messages, gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleSystem, Content: response})
	o.cache.Set(key, messages, ttlcache.DefaultTTL)

	return response, nil
}

func (o *OpenAI) UpdateChat(ctx context.Context, key string, prompt string) (string, error) {
	if !o.cache.Has(key) {
		return "", fmt.Errorf("chat with key %q does not exist", key)
	}
	cacheItem := o.cache.Get(key)
	messages := append(cacheItem.Value(), gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleUser, Content: prompt})

	response, err := o.callChatComplete(ctx, key, messages)
	if err != nil {
		return "", err
	}

	messages = append(messages, gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleSystem, Content: response})
	o.cache.Set(key, messages, ttlcache.DefaultTTL)

	return response, nil
}

func (o *OpenAI) GetChat(key string) ([]gopenai.ChatCompletionMessage, error) {
	if !o.cache.Has(key) {
		return nil, fmt.Errorf("chat with key %q does not exist", key)
	}
	return o.cache.Get(key).Value(), nil
}

func (o *OpenAI) callChatComplete(ctx context.Context, userId string, messages []gopenai.ChatCompletionMessage) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, gopenai.ChatCompletionRequest{
		Model:     gopenai.GPT3Dot5Turbo,
		MaxTokens: o.maxTokens,
		Messages:  messages,
		User:      userId,
		ResponseFormat: &gopenai.ChatCompletionResponseFormat{
			Type: gopenai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create completion: %w", err)
	}

	response := resp.Choices[0].Message.Content
	o.logger.Debug("Chat completion", "response", response)
	return response, nil
}
