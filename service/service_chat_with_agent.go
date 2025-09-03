package main

import (
	"context"
	"log"
	"time"

	genaidemo "github.com/example/genai-foundation-demo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChatWithAgent handles chat interactions with agent capabilities
func (s *chatService) ChatWithAgent(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error) {
	startTime := time.Now()
	log.Printf("🚀 [ChatWithAgent] Starting tool-enabled chat session at %s", startTime.Format("15:04:05.000"))
	if len(messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// Use LLM processor to generate response with agent context
	result, err := s.llmProcessor.ProcessMessages(ctx, messages, temperature, maxTokens)
	if err != nil {
		return nil, err
	}

	// Add agent context to response
	enhancedContent := "[Agent Mode] " + result.Content

	tokenUsage := &TokenUsageInfo{
		InputTokens:  result.TokenUsage.InputTokens,
		OutputTokens: result.TokenUsage.OutputTokens,
		TotalTokens:  result.TokenUsage.TotalTokens,
	}

	return &ChatResult{
		Content:    enhancedContent,
		TokenUsage: tokenUsage,
	}, nil
}