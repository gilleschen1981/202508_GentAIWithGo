package main

import (
	"context"
	"fmt"

	"bitbucket.dentsplysirona.com/mirrors/langchaingo/llms"
	"bitbucket.dentsplysirona.com/mirrors/langchaingo/llms/googleai"
	"bitbucket.dentsplysirona.com/mirrors/langchaingo/llms/googleai/vertex"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VertexAIModelParams 定义创建 VertexAI 模型的参数
type VertexAIModelParams struct {
	Project            string
	LLMName            string
	EmbeddingModelName string
	Location           string
}

// VertexAIChatParams 定义聊天参数
type VertexAIChatParams struct {
	Temperature float64
	MaxToken    int
}

const (
	GlobalRegion   = "global"
	GlobalEndpoint = "aiplatform.googleapis.com:443"
)

// IVertexAI 定义 VertexAI 接口
type IVertexAI interface {
	GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error)
}

// VertexAIClient VertexAI 客户端包装器
type VertexAIClient struct {
	client IVertexAI
}

// withGlobalEndPoint 设置全局端点选项
func withGlobalEndPoint(endpoint string) googleai.Option {
	return func(opts *googleai.Options) {
		opts.ClientOptions = append(opts.ClientOptions, option.WithEndpoint(endpoint))
	}
}

// NewVertexAIClient 创建新的 VertexAI 客户端
func NewVertexAIClient(modelParams VertexAIModelParams, chatParams VertexAIChatParams) (*VertexAIClient, error) {
	ctx := context.Background()

	// 构建 VertexAI 选项
	opts := []googleai.Option{
		googleai.WithCloudProject(modelParams.Project),
		googleai.WithCloudLocation(modelParams.Location),
		googleai.WithDefaultModel(modelParams.LLMName),
		googleai.WithDefaultEmbeddingModel(modelParams.EmbeddingModelName),
		googleai.WithDefaultTemperature(chatParams.Temperature),
		googleai.WithDefaultMaxTokens(chatParams.MaxToken),
	}

	// 如果是全局区域，添加全局端点
	if modelParams.Location == GlobalRegion {
		opts = append(opts, withGlobalEndPoint(GlobalEndpoint))
	}

	// 创建 VertexAI 客户端
	client, err := vertex.New(ctx, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Vertex AI client creation failed: %v", err)
	}

	return &VertexAIClient{client: client}, nil
}

// GenerateContent 生成内容
func (v *VertexAIClient) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
	content, err := v.client.GenerateContent(ctx, messages, options...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Vertex AI generate content failed: %v", err)
	}

	if content == nil || len(content.Choices) < 1 {
		return nil, status.Error(codes.Internal, "No content generated from Vertex AI")
	}

	return content, nil
}

// Call 调用 VertexAI 进行简单文本生成
func (v *VertexAIClient) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	response, err := v.client.Call(ctx, prompt, options...)
	if err != nil {
		return "", status.Errorf(codes.Internal, "Vertex AI call failed: %v", err)
	}

	return response, nil
}

// CreateEmbedding 创建文本嵌入
func (v *VertexAIClient) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings, err := v.client.CreateEmbedding(ctx, texts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Vertex AI create embedding failed: %v", err)
	}

	return embeddings, nil
}

// NewVertexAIClientFromConfig 从配置创建 VertexAI 客户端
func NewVertexAIClientFromConfig(cfg *serviceConfig) (*VertexAIClient, error) {
	modelParams := VertexAIModelParams{
		Project:            cfg.projectID,
		LLMName:            cfg.modelName,
		EmbeddingModelName: "textembedding-gecko@latest", // 默认嵌入模型
		Location:           cfg.location,
	}

	chatParams := VertexAIChatParams{
		Temperature: 0.7,  // 默认温度
		MaxToken:    2048, // 默认最大token数
	}

	return NewVertexAIClient(modelParams, chatParams)
}

// UpdateWithVertexAI 更新聊天服务以使用 VertexAI 客户端
func (s *chatService) UpdateWithVertexAI(vertexClient *VertexAIClient) {
	s.vertexClient = vertexClient
}

// GetVertexAIStats 获取 VertexAI 客户端统计信息
func (v *VertexAIClient) GetVertexAIStats() map[string]interface{} {
	return map[string]interface{}{
		"client_type": "VertexAI",
		"status":      "connected",
	}
}

// Close 关闭 VertexAI 客户端连接
func (v *VertexAIClient) Close() error {
	// VertexAI 客户端通常不需要显式关闭
	fmt.Println("VertexAI client closed")
	return nil
}
