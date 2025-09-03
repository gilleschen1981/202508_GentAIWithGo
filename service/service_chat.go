package main

import (
	"context"
	"fmt"
	"log"
	"time"

	genaidemo "github.com/example/genai-foundation-demo"
	"github.com/example/genai-foundation-demo/pkg/llm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// chatService implements the Service interface for LLM interactions
type chatService struct {
	vertexClient *VertexAIClient
	llmProcessor *llm.Processor
}

// newService creates a new chat service with VertexAI
func newService(ctx context.Context, cfg *serviceConfig) (*chatService, error) {
	fmt.Printf("🚀 Starting chat service with VertexAI\n")
	fmt.Printf("📍 Model: %s\n", cfg.modelName)
	fmt.Printf("📍 Project: %s\n", cfg.projectID)
	fmt.Printf("📍 Location: %s\n", cfg.location)

	// 创建 VertexAI 客户端
	vertexClient, err := NewVertexAIClientFromConfig(cfg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create VertexAI client: %v", err)
	}

	fmt.Printf("✅ VertexAI client initialized successfully\n")

	// 创建 LLM 处理器
	llmProcessor := llm.NewProcessor(vertexClient)

	return &chatService{
		vertexClient: vertexClient,
		llmProcessor: llmProcessor,
	}, nil
}

// Chat handles chat interactions with the LLM
func (s *chatService) Chat(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error) {
	startTime := time.Now()
	log.Printf("🚀 [Chat] Starting tool-enabled chat session at %s", startTime.Format("15:04:05.000"))
	if len(messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// 使用 LLM 处理器生成响应
	result, err := s.llmProcessor.ProcessMessages(ctx, messages, temperature, maxTokens)
	if err != nil {
		return nil, err
	}

	// 转换为服务层的结果格式
	tokenUsage := &TokenUsageInfo{
		InputTokens:  result.TokenUsage.InputTokens,
		OutputTokens: result.TokenUsage.OutputTokens,
		TotalTokens:  result.TokenUsage.TotalTokens,
	}

	return &ChatResult{
		Content:    result.Content,
		TokenUsage: tokenUsage,
	}, nil
}


// Close closes the service and cleans up resources
func (s *chatService) Close() error {
	return s.vertexClient.Close()
}
