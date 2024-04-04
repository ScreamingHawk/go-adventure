package openai

import (
	"context"
	"fmt"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/go-chi/httplog/v2"
	gopenai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	logger *httplog.Logger
	client *gopenai.Client

	maxTokens    int
	systemPrompt string

	chats map[string][]gopenai.ChatCompletionMessage
}

func NewOpenAI(cfg *config.OpenAIConfig, logger *httplog.Logger) (*OpenAI, error) {
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	client := gopenai.NewClient(cfg.ApiKey)

	systemPrompt := cfg.SystemPrompt + ". You MUST ALWAYS and ONLY respond with JSON."

	return &OpenAI{
		logger:       logger,
		client:       client,
		maxTokens:    cfg.MaxTokens,
		systemPrompt: systemPrompt,
		chats:        make(map[string][]gopenai.ChatCompletionMessage),
	}, nil
}

func (o *OpenAI) CreateChat(ctx context.Context, key string, prompt string) (string, error) {
	if _, ok := o.chats[key]; ok {
		return "", fmt.Errorf("chat with key %q already exists", key)
	}
	messages := []gopenai.ChatCompletionMessage{
		{Role: gopenai.ChatMessageRoleSystem, Content: o.systemPrompt},
		{Role: gopenai.ChatMessageRoleUser, Content: prompt},
	}
	o.chats[key] = messages

	response, err := o.callChatComplete(ctx, key, messages)
	if err != nil {
		return "", err
	}

	o.chats[key] = append(messages, gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleSystem, Content: response})
	return response, nil
}

func (o *OpenAI) UpdateChat(ctx context.Context, key string, prompt string) (string, error) {
	messages, ok := o.chats[key]
	if !ok {
		return "", fmt.Errorf("chat with key %q does not exist", key)
	}
	messages = append(messages, gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleUser, Content: prompt})
	o.chats[key] = messages

	response, err := o.callChatComplete(ctx, key, messages)
	if err != nil {
		return "", err
	}

	o.chats[key] = append(messages, gopenai.ChatCompletionMessage{Role: gopenai.ChatMessageRoleSystem, Content: response})
	return response, nil
}

func (o *OpenAI) GetChat(key string) ([]gopenai.ChatCompletionMessage, error) {
	messages, ok := o.chats[key]
	if !ok {
		return nil, fmt.Errorf("chat with key %q does not exist", key)
	}
	return messages, nil
}

func (o *OpenAI) callChatComplete(ctx context.Context, userId string, messages []gopenai.ChatCompletionMessage) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, gopenai.ChatCompletionRequest{
		Model:     gopenai.GPT3Dot5Turbo,
		MaxTokens: o.maxTokens,
		Messages:  messages,
		User: userId,
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
