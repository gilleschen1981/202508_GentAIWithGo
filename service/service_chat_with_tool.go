package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"bitbucket.dentsplysirona.com/mirrors/langchaingo/llms"
	"bitbucket.dentsplysirona.com/mirrors/langchaingo/tools/duckduckgo"
	genaidemo "github.com/example/genai-foundation-demo"
	"github.com/example/genai-foundation-demo/pkg/llm"
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

	// Let LLM decide whether to use tools automatically
	return s.processWithLLMTools(ctx, messages, temperature, maxTokens, startTime)
}

func (s *chatService) processWithLLMTools(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32, startTime time.Time) (*ChatResult, error) {
	log.Printf("🔧 [processWithLLMTools] Starting LLM tool processing...")

	// Create tool definitions for LLM
	tools := s.createLLMTools()

	// Convert messages to langchain format
	llmMessages := llm.ConvertToLangchainMessages(messages)

	// Prepare call options with tools
	callOptions := []llms.CallOption{
		llms.WithTools(tools),
	}

	if temperature != nil {
		callOptions = append(callOptions, llms.WithTemperature(float64(*temperature)))
	}
	if maxTokens != nil {
		callOptions = append(callOptions, llms.WithMaxTokens(int(*maxTokens)))
	}

	// Call LLM with tools
	response, err := s.vertexClient.client.GenerateContent(ctx, llmMessages, callOptions...)
	if err != nil {
		log.Printf("❌ [processWithLLMTools] LLM call failed: %v", err)
		return nil, status.Error(codes.Internal, "LLM tool processing failed")
	}

	// Process tool calls if any
	content, err := s.processToolCalls(ctx, response)
	if err != nil {
		log.Printf("❌ [processWithLLMTools] Tool call processing failed: %v", err)
		return nil, status.Error(codes.Internal, "tool call processing failed")
	}

	enhancedContent := fmt.Sprintf("[Tool Mode] %s", content)

	// Estimate token usage (not available in ContentChoice, set to zero)
	tokenUsage := &TokenUsageInfo{
		InputTokens:  0,
		OutputTokens: 0,
		TotalTokens:  0,
	}

	log.Printf("✅ [processWithLLMTools] Completed in %v", time.Since(startTime))

	return &ChatResult{
		Content:    enhancedContent,
		TokenUsage: tokenUsage,
	}, nil
}

func (s *chatService) createLLMTools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "search_web",
				Description: "Search the web for current information, news, weather, facts, etc.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "The search query to find information on the web",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "calculate",
				Description: "Perform basic arithmetic calculations (addition, subtraction, multiplication, division)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"expression": map[string]interface{}{
							"type":        "string",
							"description": "The mathematical expression to calculate (e.g., '5+3', '10*2', '15/3')",
						},
					},
					"required": []string{"expression"},
				},
			},
		},
	}
}

func (s *chatService) processToolCalls(ctx context.Context, response *llms.ContentResponse) (string, error) {
	// Check if there are tool calls in the response
	if len(response.Choices) == 0 {
		return "No response from LLM", nil
	}

	choice := response.Choices[0]

	// If there are no tool calls, return the text content
	if len(choice.ToolCalls) == 0 {
		return choice.Content, nil
	}

	// Process tool calls
	var results []string
	for _, toolCall := range choice.ToolCalls {
		result, err := s.executeToolCall(ctx, toolCall)
		if err != nil {
			log.Printf("❌ [processToolCalls] Tool call failed: %v", err)
			results = append(results, fmt.Sprintf("Tool call failed: %v", err))
		} else {
			results = append(results, result)
		}
	}

	// Combine text content and tool results
	finalContent := choice.Content
	if len(results) > 0 {
		finalContent += "\n\nTool Results:\n" + strings.Join(results, "\n")
	}

	return finalContent, nil
}

func (s *chatService) executeToolCall(ctx context.Context, toolCall llms.ToolCall) (string, error) {
	switch toolCall.FunctionCall.Name {
	case "search_web":
		return s.executeSearchTool(ctx, toolCall.FunctionCall.Arguments)
	case "calculate":
		return s.executeCalculatorTool(toolCall.FunctionCall.Arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolCall.FunctionCall.Name)
	}
}

func (s *chatService) executeSearchTool(ctx context.Context, arguments string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse search arguments: %v", err)
	}

	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid query parameter")
	}

	log.Printf("🔍 [executeSearchTool] Performing search for: %s", query)

	duckduckgoTool, err := duckduckgo.New(5, "Mozilla/5.0 (compatible; GenAI-Service/1.0)")
	if err != nil {
		return "", fmt.Errorf("failed to initialize DuckDuckGo tool: %v", err)
	}

	result, err := duckduckgoTool.Call(ctx, query)
	if err != nil {
		return "", fmt.Errorf("search failed: %v", err)
	}

	log.Printf("✅ [executeSearchTool] Search completed successfully")
	return result, nil
}

func (s *chatService) executeCalculatorTool(arguments string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse calculator arguments: %v", err)
	}

	expression, ok := args["expression"].(string)
	if !ok {
		return "", fmt.Errorf("missing or invalid expression parameter")
	}

	log.Printf("🧮 [executeCalculatorTool] Calculating: %s", expression)

	// Parse and calculate the expression
	result, err := s.evaluateExpression(expression)
	if err != nil {
		return "", fmt.Errorf("calculation failed: %v", err)
	}

	log.Printf("✅ [executeCalculatorTool] Calculation completed: %s", result)
	return result, nil
}

func (s *chatService) evaluateExpression(expression string) (string, error) {
	// Clean the expression
	expr := strings.ReplaceAll(expression, " ", "")

	// Simple expression parser for basic operations
	for _, op := range []string{"+", "-", "*", "/"} {
		if strings.Contains(expr, op) {
			parts := strings.Split(expr, op)
			if len(parts) == 2 {
				left, err := strconv.ParseFloat(parts[0], 64)
				if err != nil {
					return "", fmt.Errorf("invalid left operand: %s", parts[0])
				}

				right, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					return "", fmt.Errorf("invalid right operand: %s", parts[1])
				}

				var result float64
				switch op {
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
				}

				if result == float64(int64(result)) {
					return fmt.Sprintf("%s = %.0f", expression, result), nil
				}
				return fmt.Sprintf("%s = %.2f", expression, result), nil
			}
		}
	}

	return "", fmt.Errorf("unsupported expression format: %s", expression)
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
