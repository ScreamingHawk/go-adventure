package openai_test

import (
	"context"
	"testing"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/ScreamingHawk/go-adventure/narrator"
	"github.com/ScreamingHawk/go-adventure/openai"
	"github.com/go-chi/httplog/v2"
)

var systemPrompt = narrator.EnforceJSONPrompt("Begin the story. Make it super short.")

func makeOpenAI(t *testing.T) *openai.OpenAI {
	logger := httplog.NewLogger("test")
	cfg := &config.Config{}
	err := config.NewFromFile("../etc/app.test.conf", cfg)
	if err != nil {
		t.Fatal("failed to load config")
	}
	openAI, err := openai.NewOpenAI(&cfg.OpenAI, logger)
	if err != nil {
		t.Fatalf("failed to create openai: %v", err)
	}
	return openAI
}

func doCreateChat(t *testing.T, openAI *openai.OpenAI, storyKey string, prompt string) string {
	result, err := openAI.CreateChat(context.Background(), storyKey, prompt)
	if err != nil {
		t.Fatalf("failed to create completion: %v", err)
	}
	if result == "" {
		t.Fatal("expected completion to be non-empty")
	}
	// Log result
	t.Log(result)
	return result
}

func TestOpenAICreateChat(t *testing.T) {
	openAI := makeOpenAI(t)
	storyKey := "test"
	doCreateChat(t, openAI, storyKey, systemPrompt)
}

func TestOpenAIUpdateChat(t *testing.T) {
	openAI := makeOpenAI(t)
	storyKey := "test"
	result := doCreateChat(t, openAI, storyKey, systemPrompt)

	chatResponse, err := narrator.DecodeChatResponse(result)
	if err != nil {
		t.Fatalf("failed to decode chat response: %v", err)
	}

	// Continue with the first response
	result, err = openAI.UpdateChat(context.Background(), storyKey, chatResponse.Choices[0])
	if err != nil {
		t.Fatalf("failed to create completion: %v", err)
	}
	if result == "" {
		t.Fatal("expected completion to be non-empty")
	}
	// Log result
	t.Log(result)
}
