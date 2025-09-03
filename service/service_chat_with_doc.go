package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	genaidemo "github.com/example/genai-foundation-demo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChromaDBQueryRequest represents the request structure for ChromaDB queries
type ChromaDBQueryRequest struct {
	Query    string `json:"query"`
	NResults int    `json:"n_results"`
}

// ChromaDBQueryResponse represents the response structure from ChromaDB
type ChromaDBQueryResponse struct {
	Documents []string                 `json:"documents"`
	Metadatas []map[string]interface{} `json:"metadatas"`
	Distances []float64                `json:"distances"`
	IDs       []string                 `json:"ids"`
}

// queryChromaDB searches ChromaDB for relevant documents
func (s *chatService) queryChromaDB(ctx context.Context, query string, nResults int) (*ChromaDBQueryResponse, error) {
	chromaDBURL := "http://localhost:8000/query"

	reqBody := ChromaDBQueryRequest{
		Query:    query,
		NResults: nResults,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", chromaDBURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query ChromaDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ChromaDB query failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var queryResp ChromaDBQueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &queryResp, nil
}

// ChatWithDoc handles chat interactions with document capabilities using RAG
func (s *chatService) ChatWithDoc(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ChatResult, error) {
	startTime := time.Now()
	log.Printf("🚀 [ChatWithDoc] Starting RAG-enabled chat session at %s", startTime.Format("15:04:05.000"))
	if len(messages) == 0 {
		return nil, status.Error(codes.InvalidArgument, "messages cannot be empty")
	}

	// 1. Extract user query from last message
	lastMessage := messages[len(messages)-1]
	userQuery := lastMessage.Content
	log.Printf("📝 [ChatWithDoc] User query: %s", userQuery)

	// 2. Search ChromaDB for relevant documents
	log.Printf("🔍 [ChatWithDoc] Searching ChromaDB for relevant documents...")
	chromaResp, err := s.queryChromaDB(ctx, userQuery, 3)
	if err != nil {
		log.Printf("⚠️ [ChatWithDoc] ChromaDB query failed: %v", err)
		// Fallback to normal chat without RAG
		result, err := s.llmProcessor.ProcessMessages(ctx, messages, temperature, maxTokens)
		if err != nil {
			return nil, err
		}
		enhancedContent := "[Doc Mode - ChromaDB unavailable] " + result.Content
		return &ChatResult{
			Content: enhancedContent,
			TokenUsage: &TokenUsageInfo{
				InputTokens:  result.TokenUsage.InputTokens,
				OutputTokens: result.TokenUsage.OutputTokens,
				TotalTokens:  result.TokenUsage.TotalTokens,
			},
		}, nil
	}

	log.Printf("📚 [ChatWithDoc] Found %d relevant documents", len(chromaResp.Documents))

	// 3. Enhance prompt with retrieved context
	contextDocs := ""
	for i, doc := range chromaResp.Documents {
		filename := "unknown"
		if len(chromaResp.Metadatas) > i {
			if fn, ok := chromaResp.Metadatas[i]["filename"].(string); ok {
				filename = fn
			}
		}
		distance := 0.0
		if len(chromaResp.Distances) > i {
			distance = chromaResp.Distances[i]
		}
		contextDocs += fmt.Sprintf("\n\n--- Document %d (from: %s, relevance: %.3f) ---\n%s", i+1, filename, 1.0-distance, doc)
	}

	// Create enhanced messages with document context
	enhancedMessages := make([]*genaidemo.Message, 0, len(messages)+1)

	// Add system message with document context
	systemMessage := &genaidemo.Message{
		Role:    genaidemo.Role_ROLE_SYSTEM,
		Content: fmt.Sprintf("You are a helpful AI assistant with access to relevant documents. Use the following document excerpts to help answer the user's question:\n\n=== RELEVANT DOCUMENTS ===%s\n\n=== END DOCUMENTS ===\n\nWhen answering, reference specific information from the documents when relevant. If the documents don't contain information to answer the question, say so clearly.", contextDocs),
	}
	enhancedMessages = append(enhancedMessages, systemMessage)

	// Add original messages
	enhancedMessages = append(enhancedMessages, messages...)

	log.Printf("🔄 [ChatWithDoc] Processing enhanced prompt with %d total messages", len(enhancedMessages))

	// 4. Generate response using LLM with enhanced context
	result, err := s.llmProcessor.ProcessMessages(ctx, enhancedMessages, temperature, maxTokens)
	if err != nil {
		return nil, err
	}

	// Add RAG indicator to response
	enhancedContent := "[RAG-Enhanced] " + result.Content

	tokenUsage := &TokenUsageInfo{
		InputTokens:  result.TokenUsage.InputTokens,
		OutputTokens: result.TokenUsage.OutputTokens,
		TotalTokens:  result.TokenUsage.TotalTokens,
	}

	log.Printf("✅ [ChatWithDoc] RAG response generated successfully in %v", time.Since(startTime))
	return &ChatResult{
		Content:    enhancedContent,
		TokenUsage: tokenUsage,
	}, nil
}
