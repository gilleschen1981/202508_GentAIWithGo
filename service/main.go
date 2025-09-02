package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/example/genai-foundation-demo"
)

const (
	serviceName = "genai-chat-service"
	httpPort    = "8080"
)

type serviceConfig struct {
	projectID string
	location  string
	modelName string
}

func main() {
	log.Printf("starting service %s", serviceName)

	cfg, err := getConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to get service config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler, err := createHandler(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}
	defer func() { _ = handler.Close() }()

	// Start HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/chat", createHTTPHandler(handler, "Chat"))
	mux.HandleFunc("/api/chat-with-tool", createHTTPHandler(handler, "ChatWithTool"))
	mux.HandleFunc("/api/chat-with-agent", createHTTPHandler(handler, "ChatWithAgent"))
	mux.HandleFunc("/api/chat-with-doc", createHTTPHandler(handler, "ChatWithDoc"))
	mux.HandleFunc("/api/health", healthHandler)
	
	log.Printf("🌐 HTTP server starting on port %s", httpPort)
	log.Printf("📍 API endpoints:")
	log.Printf("   - POST /api/chat")
	log.Printf("   - POST /api/chat-with-tool")
	log.Printf("   - POST /api/chat-with-agent")
	log.Printf("   - POST /api/chat-with-doc")
	log.Printf("   - GET  /api/health")
	
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}

	log.Printf("stopped %s service", serviceName)
}

func getConfigFromEnv() (*serviceConfig, error) {
	// 使用默认配置 (在 config.go 中定义)
	config := &serviceConfig{
		projectID: DefaultProjectID,
		location:  DefaultLocation,
		modelName: DefaultModelName,
	}
	
	// 如果设置了环境变量，优先使用环境变量
	if envProjectID := os.Getenv("GCP_PROJECT_ID"); envProjectID != "" {
		config.projectID = envProjectID
		log.Printf("Using project ID from environment: %s", envProjectID)
	}
	if envLocation := os.Getenv("VERTEX_AI_LOCATION"); envLocation != "" {
		config.location = envLocation
		log.Printf("Using location from environment: %s", envLocation)
	}
	if envModel := os.Getenv("VERTEX_AI_MODEL"); envModel != "" {
		config.modelName = envModel
		log.Printf("Using model from environment: %s", envModel)
	}
	
	log.Printf("VertexAI Config - Project: %s, Location: %s, Model: %s", 
		config.projectID, config.location, config.modelName)
	
	// 验证配置
	if config.projectID == "your-gcp-project-id" {
		log.Printf("⚠️  Warning: Please update DefaultProjectID in config.go with your actual GCP project ID")
	}
	
	return config, nil
}

func createHandler(ctx context.Context, cfg *serviceConfig) (*Handler, error) {
	service, err := newService(ctx, cfg)
	if err != nil {
		return nil, err
	}

	handler, err := newHandler(service)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// HTTP Handler types
type HTTPMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type HTTPChatRequest struct {
	Messages    []HTTPMessage `json:"messages"`
	Temperature *float32      `json:"temperature,omitempty"`
	MaxTokens   *int32        `json:"max_tokens,omitempty"`
}

type HTTPChatResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// Create HTTP handler for gRPC service methods
func createHTTPHandler(handler *Handler, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req HTTPChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendErrorResponse(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Convert HTTP messages to gRPC messages
		grpcMessages := make([]*genaidemo.Message, len(req.Messages))
		for i, msg := range req.Messages {
			grpcMessages[i] = &genaidemo.Message{
				Role:    parseRole(msg.Role),
				Content: msg.Content,
			}
		}

		// Create gRPC request
		grpcReq := &genaidemo.ChatRequest{
			Messages:    grpcMessages,
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
		}

		// Call appropriate gRPC method
		var grpcResp *genaidemo.ChatResponse
		var err error
		
		switch method {
		case "Chat":
			grpcResp, err = handler.Chat(r.Context(), grpcReq)
		case "ChatWithTool":
			grpcResp, err = handler.ChatWithTool(r.Context(), grpcReq)
		case "ChatWithAgent":
			grpcResp, err = handler.ChatWithAgent(r.Context(), grpcReq)
		case "ChatWithDoc":
			grpcResp, err = handler.ChatWithDoc(r.Context(), grpcReq)
		default:
			sendErrorResponse(w, "Unknown method", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Printf("❌ gRPC call failed: %v", err)
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send response
		response := HTTPChatResponse{
			Content: grpcResp.Content,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func parseRole(role string) genaidemo.Role {
	switch role {
	case "ROLE_USER":
		return genaidemo.Role_ROLE_USER
	case "ROLE_ASSISTANT":
		return genaidemo.Role_ROLE_ASSISTANT
	case "ROLE_SYSTEM":
		return genaidemo.Role_ROLE_SYSTEM
	default:
		return genaidemo.Role_ROLE_USER
	}
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := HTTPChatResponse{
		Error: message,
	}
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "healthy",
		"service": "genai-foundation-demo",
	}
	json.NewEncoder(w).Encode(response)
}