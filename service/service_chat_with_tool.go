package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.dentsplysirona.com/mirrors/langchaingo/tools/duckduckgo"
	genaidemo "github.com/example/genai-foundation-demo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *chatService) ChatWithTool(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error) {
	startTime := time.Now()
	log.Printf("🚀 [ChatWithTool] Starting tool-enabled chat session at %s", startTime.Format("15:04:05.000"))
	
	if len(messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	lastMessage := messages[len(messages)-1]
	if lastMessage.Role != genaidemo.Role_ROLE_USER {
		return nil, status.Error(codes.InvalidArgument, "last message must be from user")
	}

	userQuery := lastMessage.Content
	log.Printf("🔍 [ChatWithTool] Processing query: '%s'", userQuery)

	if s.shouldUseTools(userQuery) {
		return s.processWithTools(ctx, userQuery, startTime)
	}

	return s.fallbackToBasicChat(ctx, messages, temperature, maxTokens)
}

func (s *chatService) shouldUseTools(query string) bool {
	log.Printf("🤖 [shouldUseTools] Using LLM to determine tool usage for query: %s", query)
	
	ctx := context.Background()
	
	systemPrompt := `You are a tool usage analyzer. Given a user query, determine if it needs external tools to answer properly.

Available tools:
1. Search tool (DuckDuckGo) - for finding current information, news, weather, facts, etc.
2. Calculator tool - for mathematical calculations

Respond with ONLY "YES" if tools are needed, or "NO" if the query can be answered with general knowledge.

Examples:
- "What's the weather today?" -> YES (needs search)
- "Calculate 5 + 3" -> YES (needs calculator) 
- "What is machine learning?" -> NO (general knowledge)
- "Tell me a joke" -> NO (general knowledge)
- "今天天气几度?" -> YES (needs search)
- "计算 10 * 5" -> YES (needs calculator)`

	messages := []*genaidemo.Message{
		{
			Role:    genaidemo.Role_ROLE_SYSTEM,
			Content: systemPrompt,
		},
		{
			Role:    genaidemo.Role_ROLE_USER,
			Content: query,
		},
	}
	
	result, err := s.llmProcessor.ProcessMessages(ctx, messages, nil, nil)
	if err != nil {
		log.Printf("❌ [shouldUseTools] LLM call failed: %v", err)
		return false
	}
	
	response := strings.TrimSpace(strings.ToUpper(result.Content))
	needsTools := strings.Contains(response, "YES")
	
	log.Printf("🤖 [shouldUseTools] LLM response: %s -> Tools needed: %t", response, needsTools)
	return needsTools
}

func (s *chatService) processWithTools(ctx context.Context, query string, startTime time.Time) (*ChatResult, error) {
	log.Printf("🔧 [processWithTools] Starting tool processing...")
	
	var result string
	var err error
	
	if s.isCalculationQuery(query) {
		result, err = s.handleCalculation(query)
	} else {
		result, err = s.handleSearch(ctx, query)
	}
	
	if err != nil {
		log.Printf("❌ [processWithTools] Tool processing failed: %v", err)
		return nil, status.Error(codes.Internal, "tool processing failed")
	}
	
	enhancedContent := fmt.Sprintf("[Tool Mode] %s", result)
	
	tokenUsage := &TokenUsageInfo{
		InputTokens:  int32(len(query) / 4),
		OutputTokens: int32(len(enhancedContent) / 4),
		TotalTokens:  int32((len(query) + len(enhancedContent)) / 4),
	}
	
	log.Printf("✅ [processWithTools] Completed in %v", time.Since(startTime))
	
	return &ChatResult{
		Content:    enhancedContent,
		TokenUsage: tokenUsage,
	}, nil
}

func (s *chatService) isCalculationQuery(query string) bool {
	queryLower := strings.ToLower(query)
	calcIndicators := []string{
		"calculate", "compute", "math", "计算", "算",
		"+", "-", "*", "/", "=", "×", "÷",
	}
	
	for _, indicator := range calcIndicators {
		if strings.Contains(queryLower, indicator) {
			return true
		}
	}
	return false
}

func (s *chatService) handleCalculation(query string) (string, error) {
	log.Printf("🧮 [handleCalculation] Processing calculation: %s", query)
	
	queryLower := strings.ToLower(query)
	
	if strings.Contains(queryLower, "+") {
		return s.performBasicMath(query, "+")
	} else if strings.Contains(queryLower, "-") {
		return s.performBasicMath(query, "-")
	} else if strings.Contains(queryLower, "*") || strings.Contains(queryLower, "×") {
		return s.performBasicMath(query, "*")
	} else if strings.Contains(queryLower, "/") || strings.Contains(queryLower, "÷") {
		return s.performBasicMath(query, "/")
	}
	
	return "I can help with basic arithmetic operations (+, -, *, /). Please provide a clear mathematical expression.", nil
}

func (s *chatService) performBasicMath(query, operator string) (string, error) {
	parts := strings.Split(query, operator)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid mathematical expression")
	}
	
	leftStr := strings.TrimSpace(parts[0])
	rightStr := strings.TrimSpace(parts[1])
	
	leftStr = s.extractNumber(leftStr)
	rightStr = s.extractNumber(rightStr)
	
	left, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", leftStr)
	}
	
	right, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", rightStr)
	}
	
	var result float64
	switch operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = left / right
	default:
		return "", fmt.Errorf("unsupported operator: %s", operator)
	}
	
	if result == float64(int64(result)) {
		return fmt.Sprintf("%.0f %s %.0f = %.0f", left, operator, right, result), nil
	}
	return fmt.Sprintf("%.2f %s %.2f = %.2f", left, operator, right, result), nil
}

func (s *chatService) extractNumber(str string) string {
	var result strings.Builder
	for _, char := range str {
		if (char >= '0' && char <= '9') || char == '.' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func (s *chatService) handleSearch(ctx context.Context, query string) (string, error) {
	log.Printf("🔍 [handleSearch] Performing search for: %s", query)
	
	duckduckgoTool, err := duckduckgo.New(5, "Mozilla/5.0 (compatible; GenAI-Service/1.0)")
	if err != nil {
		return "", fmt.Errorf("failed to initialize DuckDuckGo tool: %v", err)
	}
	
	searchResult, err := duckduckgoTool.Call(ctx, query)
	if err != nil {
		return "", fmt.Errorf("search failed: %v", err)
	}
	
	log.Printf("✅ [handleSearch] Search completed successfully")
	return searchResult, nil
}

func (s *chatService) fallbackToBasicChat(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error) {
	log.Printf("💬 [fallbackToBasicChat] Using basic LLM processing...")
	
	result, err := s.llmProcessor.ProcessMessages(ctx, messages, temperature, maxTokens)
	if err != nil {
		return nil, err
	}
	
	enhancedContent := "[Tool Mode] " + result.Content
	
	return &ChatResult{
		Content: enhancedContent,
		TokenUsage: &TokenUsageInfo{
			InputTokens:  result.TokenUsage.InputTokens,
			OutputTokens: result.TokenUsage.OutputTokens,
			TotalTokens:  result.TokenUsage.TotalTokens,
		},
	}, nil
}