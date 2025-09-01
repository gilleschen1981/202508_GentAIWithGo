# 快速启动指南

## 🚀 5分钟快速运行

### 1. 确保已安装并认证 gcloud
```bash
# 安装 gcloud (如果未安装)
# https://cloud.google.com/sdk/docs/install

# 认证
gcloud auth login
gcloud auth application-default login

# 设置项目 (替换为你的项目ID)
gcloud config set project YOUR_PROJECT_ID
```

### 2. 启用 VertexAI API
```bash
gcloud services enable aiplatform.googleapis.com
```

### 3. 修改配置
编辑 `service/config.go` 文件，将第7行修改为你的项目ID：
```go
DefaultProjectID = "你的实际项目ID"  // 替换这里
```

### 4. 运行服务
```bash
# 安装依赖
go mod tidy

# 启动服务 (推荐使用脚本)
./run.sh

# 或手动启动
cd service && go build -o ../genai-service . && cd .. && ./genai-service
```

### 5. 测试服务
打开浏览器，访问 `frontend/index.html` 即可开始聊天！

**🎉 新功能：4种聊天模式**
- **🔵 Chat**: 基础对话
- **🟢 Tool**: 工具增强 (搜索、计算等)
- **🟠 Agent**: 智能代理 (任务规划、协调)
- **🟣 Doc**: 文档分析 (研究、总结)

## 📋 配置说明

### 项目结构
```
├── service/
│   ├── config.go      # 🔧 主要配置文件 - 在这里修改你的设置
│   ├── main.go        # 🚀 程序入口
│   ├── service_chat.go # 🤖 LLM 交互逻辑 (4种模式)
│   ├── handler.go     # 📡 gRPC 处理器 (4个接口)
│   └── client.go      # 🔗 VertexAI 客户端封装
├── internal.proto     # 🔌 API 定义 (4个专业化接口)
├── pkg/llm/           # 📦 LLM 处理抽象层
└── frontend/          # 🌐 多模式聊天界面
```

### 重要配置项
在 `service/config.go` 中修改：

```go
const (
    // 🏷️ 你的 GCP 项目 ID (必须修改)
    DefaultProjectID = "your-gcp-project-id"
    
    // 🌍 服务区域 (可选)
    // 选项: us-central1, us-east1, europe-west1, asia-southeast1
    DefaultLocation = "us-central1"
    
    // 🤖 AI 模型 (可选)  
    // 推荐: gemini-1.5-flash (快速经济)
    // 高级: gemini-1.5-pro (功能强大)
    DefaultModelName = "gemini-1.5-flash"
)
```

## 🔧 常见问题

### Q: 出现认证错误怎么办？
```bash
# 重新认证
gcloud auth application-default login
```

### Q: 项目ID在哪里找？
```bash
# 查看当前项目
gcloud config get-value project

# 列出所有项目
gcloud projects list
```

### Q: 想使用更强大的模型？
修改 `service/config.go` 中的 `DefaultModelName`：
- `gemini-1.5-flash` - 快速经济 ⚡
- `gemini-1.5-pro` - 功能强大 🧠
- `gemini-1.0-pro` - 稳定可靠 🛡️

### Q: 想部署到其他区域？
修改 `service/config.go` 中的 `DefaultLocation`：
- 美国: `us-central1`, `us-east1`, `us-west1`  
- 欧洲: `europe-west1`, `europe-west4`
- 亚洲: `asia-southeast1`, `asia-northeast1`

## 🎯 下一步

1. **测试4个gRPC接口**: 使用 grpcurl 或 Postman 测试不同模式
2. **集成到你的应用**: 通过 gRPC 客户端调用专业化接口
3. **自定义响应逻辑**: 修改 `service/service_chat.go` 中的不同模式实现
4. **扩展工具功能**: 在 Tool 模式中添加真实的工具集成
5. **实现智能代理**: 在 Agent 模式中添加任务规划和执行
6. **文档分析功能**: 在 Doc 模式中集成文档处理能力
7. **添加流式响应**: 实现 Server-side Streaming
8. **部署生产环境**: 添加监控、日志、认证等