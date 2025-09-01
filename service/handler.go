package main

import (
	"context"
	"errors"

	genaidemo "github.com/example/genai-foundation-demo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service describes an API for managing chat interactions with LLM.
type Service interface {
	Chat(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error)
	ChatWithTool(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error)
	ChatWithAgent(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error)
	ChatWithDoc(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error)
	Close() error
}

// ChatResult represents the result of a chat interaction
type ChatResult struct {
	Content    string
	TokenUsage *TokenUsageInfo
}

// TokenUsageInfo contains token usage statistics
type TokenUsageInfo struct {
	InputTokens  int32
	OutputTokens int32
	TotalTokens  int32
}


// Handler is handling incoming gRPC requests
type Handler struct {
	genaidemo.UnimplementedChatServiceServer
	service Service
}

// newHandler creates a new handler with the given service
func newHandler(service Service) (*Handler, error) {
	if service == nil {
		return nil, errors.New("service must be set")
	}

	return &Handler{
		service: service,
	}, nil
}

// Chat handles the Chat gRPC method
func (h *Handler) Chat(ctx context.Context, req *genaidemo.ChatRequest) (*genaidemo.ChatResponse, error) {
	if req.Messages == nil || len(req.Messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// Validate messages
	for i, msg := range req.Messages {
		if msg.Content == "" {
			return nil, status.Errorf(codes.InvalidArgument, "message content cannot be empty at index %d", i)
		}
		if msg.Role == genaidemo.Role_ROLE_UNKNOWN {
			return nil, status.Errorf(codes.InvalidArgument, "invalid message role at index %d", i)
		}
	}

	result, err := h.service.Chat(ctx, req.Messages, req.Temperature, req.MaxTokens)
	if err != nil {
		return nil, err
	}

	response := &genaidemo.ChatResponse{
		Content: result.Content,
	}

	if result.TokenUsage != nil {
		response.TokenUsage = &genaidemo.TokenUsage{
			InputTokenNum:  result.TokenUsage.InputTokens,
			OutputTokenNum: result.TokenUsage.OutputTokens,
			TotalTokenNum:  result.TokenUsage.TotalTokens,
		}
	}

	return response, nil
}

// ChatWithTool handles the ChatWithTool gRPC method
func (h *Handler) ChatWithTool(ctx context.Context, req *genaidemo.ChatRequest) (*genaidemo.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// Validate messages
	for i, msg := range req.Messages {
		if msg.Content == "" {
			return nil, status.Errorf(codes.InvalidArgument, "message content cannot be empty at index %d", i)
		}
		if msg.Role == genaidemo.Role_ROLE_UNKNOWN {
			return nil, status.Errorf(codes.InvalidArgument, "invalid message role at index %d", i)
		}
	}

	result, err := h.service.ChatWithTool(ctx, req.Messages, req.Temperature, req.MaxTokens)
	if err != nil {
		return nil, err
	}

	response := &genaidemo.ChatResponse{
		Content: result.Content,
	}

	if result.TokenUsage != nil {
		response.TokenUsage = &genaidemo.TokenUsage{
			InputTokenNum:  result.TokenUsage.InputTokens,
			OutputTokenNum: result.TokenUsage.OutputTokens,
			TotalTokenNum:  result.TokenUsage.TotalTokens,
		}
	}

	return response, nil
}

// ChatWithAgent handles the ChatWithAgent gRPC method
func (h *Handler) ChatWithAgent(ctx context.Context, req *genaidemo.ChatRequest) (*genaidemo.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// Validate messages
	for i, msg := range req.Messages {
		if msg.Content == "" {
			return nil, status.Errorf(codes.InvalidArgument, "message content cannot be empty at index %d", i)
		}
		if msg.Role == genaidemo.Role_ROLE_UNKNOWN {
			return nil, status.Errorf(codes.InvalidArgument, "invalid message role at index %d", i)
		}
	}

	result, err := h.service.ChatWithAgent(ctx, req.Messages, req.Temperature, req.MaxTokens)
	if err != nil {
		return nil, err
	}

	response := &genaidemo.ChatResponse{
		Content: result.Content,
	}

	if result.TokenUsage != nil {
		response.TokenUsage = &genaidemo.TokenUsage{
			InputTokenNum:  result.TokenUsage.InputTokens,
			OutputTokenNum: result.TokenUsage.OutputTokens,
			TotalTokenNum:  result.TokenUsage.TotalTokens,
		}
	}

	return response, nil
}

// ChatWithDoc handles the ChatWithDoc gRPC method
func (h *Handler) ChatWithDoc(ctx context.Context, req *genaidemo.ChatRequest) (*genaidemo.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// Validate messages
	for i, msg := range req.Messages {
		if msg.Content == "" {
			return nil, status.Errorf(codes.InvalidArgument, "message content cannot be empty at index %d", i)
		}
		if msg.Role == genaidemo.Role_ROLE_UNKNOWN {
			return nil, status.Errorf(codes.InvalidArgument, "invalid message role at index %d", i)
		}
	}

	result, err := h.service.ChatWithDoc(ctx, req.Messages, req.Temperature, req.MaxTokens)
	if err != nil {
		return nil, err
	}

	response := &genaidemo.ChatResponse{
		Content: result.Content,
	}

	if result.TokenUsage != nil {
		response.TokenUsage = &genaidemo.TokenUsage{
			InputTokenNum:  result.TokenUsage.InputTokens,
			OutputTokenNum: result.TokenUsage.OutputTokens,
			TotalTokenNum:  result.TokenUsage.TotalTokens,
		}
	}

	return response, nil
}

// Close all resources created by the handler
func (h *Handler) Close() error {
	return h.service.Close()
}
