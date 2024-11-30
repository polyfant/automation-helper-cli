package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Assistant struct {
	client *openai.Client
}

func NewAssistant(apiKey string) *Assistant {
	return &Assistant{
		client: openai.NewClient(apiKey),
	}
}

func (a *Assistant) GetHelp(question string) (string, error) {
	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are an expert in ABB RAPID robotics programming language. 
					Help users understand and modify their RAPID code. Provide clear, 
					practical explanations and examples.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("AI request failed: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}