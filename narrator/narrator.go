package narrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/ScreamingHawk/go-adventure/openai"
	"github.com/go-chi/httplog/v2"
)

var initPrompts = []string{
	"Write a story about a hero who is on a quest to save the world.",
	"Write a story about a unicorn who is on a quest to find their family.",
	"Write a story about a dragon who is on a quest to find treasure.",
	"Write a story about a wizard who is on a quest to find a lost spell.",
	"Write a story about a knight who is on a quest to find a lost sword.",
	"Write a story about a princess who is on a quest to find a lost crown.",
	"Write a story about a king who is on a quest to find a lost kingdom.",
	"Write a story about a queen who is on a quest to find a lost throne.",
	"Write a story about a prince who is on a quest to find a lost castle.",
	"Write a story about a villain who is on a quest to take over the world.",
}

type Narrator struct {
	logger *httplog.Logger

	OpenAI *openai.OpenAI
}

type ChatResponse struct {
	Plot    string   `json:"plot"`
	Choices []string `json:"choices"`
	Hidden  string   `json:"hidden"`
}

func NewNarrator(cfg *config.OpenAIConfig, logger *httplog.Logger) (*Narrator, error) {
	cfg.SystemPrompt = EnforceJSONPrompt(cfg.SystemPrompt)
	openAI, err := openai.NewOpenAI(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &Narrator{
		logger: logger,
		OpenAI: openAI,
	}, nil
}

func (n *Narrator) CreateStory(ctx context.Context, sessionKey string) (*ChatResponse, error) {
	// Pseudo randomly select a prompt
	prompt := initPrompts[int(sessionKey[0])%len(initPrompts)]
	n.logger.Debug("Creating story", "op", "create_story", "prompt", prompt)

	response, err := n.OpenAI.CreateChat(ctx, sessionKey, prompt)
	if err != nil {
		return nil, err
	}
	return DecodeChatResponse(response)
}

func (n *Narrator) UpdateStory(ctx context.Context, sessionKey, choice string) (*ChatResponse, error) {
	messages, err := n.OpenAI.GetChat(sessionKey)
	if err != nil {
		return nil, err
	}

	// Check choice is valid
	lastMessage := messages[len(messages)-1]
	lastResponse, err := DecodeChatResponse(lastMessage.Content)
	if err != nil {
		return nil, err
	}
	found := false
	for _, c := range lastResponse.Choices {
		if c == choice {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("choice %q is not valid", choice)
	}

	response, err := n.OpenAI.UpdateChat(ctx, sessionKey, choice)
	if err != nil {
		return nil, err
	}
	return DecodeChatResponse(response)
}

func EnforceJSONPrompt(systemPrompt string) (string) {
	return systemPrompt + ". You MUST ALWAYS and ONLY respond with JSON. Your JSON contains the keys 'plot', 'choices', and 'hidden'. 'plot' is a string of the story so far. 'choices' is a string list of choices the user can make. 'hidden' is a string of extra information to keep in mind when story telling that won't be displayed to the user."
}

func DecodeChatResponse(response string) (*ChatResponse, error) {
	// Decode JSON
	var chatResponse ChatResponse
	if err := json.Unmarshal([]byte(response), &chatResponse); err != nil {
		return nil, fmt.Errorf("failed to decode chat response: %w", err)
	}
	return &chatResponse, nil
}
