// Build disabled to prevent package conflicts
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewChatServiceClient(conn)

	// Test the Chat method
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &ChatRequest{
		Messages: []*Message{
			{
				Role:    Role_ROLE_USER,
				Content: "Hello, how are you?",
			},
		},
		Temperature: floatPtr(0.7),
		MaxTokens:   int32Ptr(100),
	}

	resp, err := client.Chat(ctx, req)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Content)
	if resp.TokenUsage != nil {
		fmt.Printf("Token Usage - Input: %d, Output: %d, Total: %d\n",
			resp.TokenUsage.InputTokenNum,
			resp.TokenUsage.OutputTokenNum,
			resp.TokenUsage.TotalTokenNum)
	}

	// Test another message
	req2 := &ChatRequest{
		Messages: []*Message{
			{
				Role:    Role_ROLE_USER,
				Content: "Tell me a joke",
			},
		},
	}

	resp2, err := client.Chat(ctx, req2)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Printf("\nSecond Response: %s\n", resp2.Content)
	if resp2.TokenUsage != nil {
		fmt.Printf("Token Usage - Input: %d, Output: %d, Total: %d\n",
			resp2.TokenUsage.InputTokenNum,
			resp2.TokenUsage.OutputTokenNum,
			resp2.TokenUsage.TotalTokenNum)
	}
}

// Helper functions to create pointers
func floatPtr(f float32) *float32 {
	return &f
}

func int32Ptr(i int32) *int32 {
	return &i
}