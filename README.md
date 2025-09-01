# GenAI Foundation Demo

A Go microservice that provides LLM interaction capabilities through gRPC, built with langchain-go.

## Features

- gRPC server with 4 specialized chat interfaces
- LLM integration using langchain-go and VertexAI
- Interactive web frontend with multiple chat modes
- Following the same architecture as the reference project

## Project Structure

```
├── internal.proto          # gRPC interface definition (4 chat interfaces)
├── internal.pb.go          # Generated protobuf code
├── internal_grpc.pb.go     # Generated gRPC code
├── service/                # Service implementation directory
│   ├── main.go            # Main entry point and gRPC server
│   ├── handler.go         # gRPC handler implementation (4 interfaces)
│   ├── service_chat.go    # VertexAI interaction service
│   ├── client.go          # VertexAI client wrapper
│   └── config.go          # Configuration constants
├── pkg/llm/
│   └── processor.go       # LLM processing abstraction
├── go.mod                 # Go module dependencies
├── run.sh                 # One-click startup script
├── frontend/
│   ├── index.html         # Interactive chat UI with 4 modes
│   └── chat.js            # Frontend JavaScript with mode support
├── QUICKSTART.md          # Quick setup guide
├── CLAUDE.md              # Claude Code development guide
└── README.md              # This file
```

## Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- Google Cloud Project with VertexAI enabled
- gcloud CLI installed and authenticated

## Setup

1. **Authenticate with Google Cloud:**
   ```bash
   gcloud auth login
   gcloud auth application-default login
   ```

2. **Set your Google Cloud Project:**
   ```bash
   gcloud config set project YOUR_PROJECT_ID
   ```

3. **Enable VertexAI API:**
   ```bash
   gcloud services enable aiplatform.googleapis.com
   ```

4. **Clone or download the project**

5. **Install dependencies:**
   ```bash
   go mod tidy
   ```

6. **Generate protobuf files (if needed):**
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal.proto
   ```

7. **配置 VertexAI 参数:**

   **方法一: 直接修改代码中的配置 (推荐)**
   ```go
   // 编辑 service/config.go 文件，修改以下常量:
   const (
       DefaultProjectID = "your-actual-gcp-project-id"  // 替换为你的项目ID
       DefaultLocation  = "us-central1"                 // 可选择其他区域
       DefaultModelName = "gemini-1.5-flash"            // 可选择其他模型
   )
   ```

   **方法二: 使用环境变量 (可选)**
   ```bash
   export GCP_PROJECT_ID="your-gcp-project-id"
   export VERTEX_AI_LOCATION="us-central1"
   export VERTEX_AI_MODEL="gemini-1.5-flash-001"
   ```

## Running the Service

### 快速启动 (推荐)
```bash
# 一键启动服务
./run.sh
```

### 手动启动
```bash
# 编译并启动服务
cd service && go build -o ../genai-service . && cd .. && ./genai-service

# 服务将在 50051 端口启动
```

### 测试服务
```bash
# 方法1: 打开前端页面 (推荐)
open frontend/index.html

# 方法2: 使用grpcurl测试基础接口
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"Hello!"}]}' localhost:50051 genaidemo.ChatService/Chat

# 方法3: 测试工具模式接口
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"What is the weather today?"}]}' localhost:50051 genaidemo.ChatService/ChatWithTool
```

## ✅ 项目状态

- ✅ **4个专业化聊天接口** - Chat、ChatWithTool、ChatWithAgent、ChatWithDoc
- ✅ **真实VertexAI模式运行** - 使用Gemini-1.5-flash模型
- ✅ **gRPC服务** - 完全功能的微服务架构
- ✅ **完整Token统计** - 真实的输入/输出token计数
- ✅ **多模式前端界面** - 4个按钮对应不同聊天模式
- ✅ **一键启动** - ./run.sh 脚本自动化构建和启动
- ✅ **模块化架构** - 清晰的分层设计和组件分离

## API Usage

The service exposes 4 specialized gRPC methods:

```protobuf
service ChatService {
  rpc Chat (ChatRequest) returns (ChatResponse) {}           // Basic chat
  rpc ChatWithTool(ChatRequest) returns (ChatResponse) {}    // Tool-enhanced chat
  rpc ChatWithAgent(ChatRequest) returns (ChatResponse) {}   // Agent-powered chat
  rpc ChatWithDoc(ChatRequest) returns (ChatResponse) {}     // Document-aware chat
}
```

### Interface Descriptions

- **Chat**: Basic LLM conversation interface
- **ChatWithTool**: Enhanced with external tools (web search, calculator, etc.)
- **ChatWithAgent**: Intelligent agent capabilities for complex task coordination
- **ChatWithDoc**: Document analysis and research-oriented responses

### ChatRequest

```protobuf
message ChatRequest {
  repeated Message messages = 1;
  optional float temperature = 2;
  optional int32 max_tokens = 3;
}
```

### ChatResponse

```protobuf
message ChatResponse {
  string content = 1;
  TokenUsage token_usage = 2;
}
```

## Implementation Details

- **service/main.go**: Sets up the gRPC server and initializes the service
- **service/handler.go**: Implements all 4 gRPC interfaces and handles request validation
- **service/service_chat.go**: Contains the actual LLM interaction logic using langchain-go
- **service/client.go**: VertexAI client wrapper with langchain-go integration
- **pkg/llm/processor.go**: LLM processing abstraction layer
- **internal.proto**: Defines 4 specialized gRPC interfaces

## Configuration

The service uses the following environment variables:

- `GCP_PROJECT_ID`: Your Google Cloud Project ID
- `VERTEX_AI_LOCATION`: VertexAI service location (default: us-central1)
- `VERTEX_AI_MODEL`: Model name to use (default: gemini-1.5-flash)

### Available VertexAI Models:
- `gemini-1.5-pro` - Most capable model
- `gemini-1.5-flash` - Fast and efficient (recommended)
- `gemini-1.0-pro` - Stable version

### Available Locations:
- `us-central1`, `us-east1`, `us-west1`
- `europe-west1`, `europe-west4`
- `asia-southeast1`, `asia-northeast1`

The service automatically uses gcloud authentication, so no API keys are needed.

## Frontend

The frontend provides an interactive chat interface with 4 specialized modes:

- **🔵 Chat**: Basic conversation mode
- **🟢 Tool**: Tool-enhanced responses with search and calculation capabilities  
- **🟠 Agent**: Intelligent agent mode for complex task coordination
- **🟣 Doc**: Document analysis and research mode

Since browsers cannot directly make gRPC calls, the frontend currently simulates different response modes. In a real implementation, you would typically:

1. Create a REST API gateway that translates HTTP to gRPC
2. Use gRPC-Web for direct browser gRPC communication
3. Or use server-sent events for streaming responses

## Testing

You can test the gRPC service using tools like:

- grpcurl
- BloomRPC
- Postman (with gRPC support)

Examples with grpcurl:
```bash
# Test basic Chat
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"Hello!"}]}' localhost:50051 genaidemo.ChatService/Chat

# Test ChatWithTool
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"What is the weather?"}]}' localhost:50051 genaidemo.ChatService/ChatWithTool

# Test ChatWithAgent  
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"Plan a trip to Tokyo"}]}' localhost:50051 genaidemo.ChatService/ChatWithAgent

# Test ChatWithDoc
grpcurl -plaintext -d '{"messages":[{"role":"ROLE_USER","content":"Analyze this document"}]}' localhost:50051 genaidemo.ChatService/ChatWithDoc
```