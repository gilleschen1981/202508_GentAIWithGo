# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Start Commands

### Build and Run
```bash
# Quick start with build script
./run.sh

# Manual build and run
cd service && go build -o ../genai-service . && cd .. && ./genai-service

# Install dependencies
go mod tidy

# Generate protobuf files (if needed)
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal.proto
```

### Testing
```bash
# Test with client
cd client && go build -o ../test-client . && cd .. && ./test-client

# Test with grpcurl
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"Hello!"}]}' localhost:50051 genaidemo.ChatService/Chat

# Open frontend
open frontend/index.html
```

## Architecture

This is a Go microservice that provides LLM interaction through gRPC using langchain-go and Google Cloud VertexAI.

### Core Components

- **service/main.go**: gRPC server entry point, initializes service with VertexAI configuration
- **service/handler.go**: gRPC interface implementation, validates requests and calls service layer
- **service/service_chat.go**: Core chat service that orchestrates LLM interactions
- **service/client.go**: VertexAI client wrapper using langchain-go
- **pkg/llm/processor.go**: LLM processing abstraction layer with prompts formatting
- **service/config.go**: Configuration constants for VertexAI (project, location, model)

### Key Patterns

1. **Layered Architecture**: Handler → Service → LLM Processor → VertexAI Client
2. **Configuration Priority**: Environment variables override config.go defaults
3. **Error Handling**: Uses gRPC status codes for structured error responses
4. **Token Estimation**: Simple character-based token counting fallback

### VertexAI Integration

- Uses langchain-go's VertexAI provider (`bitbucket.dentsplysirona.com/mirrors/langchaingo`)
- Requires GCP authentication via `gcloud auth application-default login`
- Supports multiple models: gemini-1.5-flash (default), gemini-1.5-pro, gemini-1.0-pro
- Auto-fallback to demo mode if VertexAI fails

### gRPC Interface

- **Chat**: Basic LLM conversation with message history
- **ChatToUseTool**: Extended chat with tool integration (DuckDuckGo, calculator)
- Message roles: ROLE_USER, ROLE_ASSISTANT, ROLE_SYSTEM
- Returns content + token usage statistics

## Configuration

Update `service/config.go` with your GCP project ID before running:

```go
const (
    DefaultProjectID = "your-gcp-project-id"  // Required
    DefaultLocation  = "us-central1"           // Optional
    DefaultModelName = "gemini-1.5-flash"     // Optional
)
```

Environment variables (optional):
- `GCP_PROJECT_ID`
- `VERTEX_AI_LOCATION` 
- `VERTEX_AI_MODEL`

## Development Notes

- Service runs on port 50051
- Frontend is demo-only (simulates gRPC calls via mock responses)
- Tool service (service_chat_tool.go) extends basic chat with external tools
- Uses mirrors of langchain-go from Dentsply Sirona BitBucket

## 测试配置
请勿在每次代码修改后自动运行测试，请等待用户确认后再运行测试。
测试相关的代码单独放到test目录，和业务代码隔离