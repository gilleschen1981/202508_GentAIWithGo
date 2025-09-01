package llm

import (
	"context"

	"bitbucket.dentsplysirona.com/mirrors/langchaingo/llms"
	"bitbucket.dentsplysirona.com/mirrors/langchaingo/prompts"
	genaidemo "github.com/example/genai-foundation-demo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Processor 封装 LLM 处理逻辑
type Processor struct {
	client Client
}

// Client 定义 LLM 客户端接口
type Client interface {
	GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error)
}

// NewProcessor 创建新的 LLM 处理器
func NewProcessor(client Client) *Processor {
	return &Processor{
		client: client,
	}
}

// ProcessResult LLM 处理结果
type ProcessResult struct {
	Content    string
	TokenUsage *TokenUsage
}

// ProcessMessages 处理消息并生成响应
func (p *Processor) ProcessMessages(ctx context.Context, messages []*genaidemo.Message, temperature *float32, maxTokens *int32) (*ProcessResult, error) {
	// 构建聊天提示模板
	chatPrompt := p.buildChatPrompt(messages)
	
	// 准备调用选项
	var options []llms.CallOption
	if temperature != nil {
		options = append(options, llms.WithTemperature(float64(*temperature)))
	}
	if maxTokens != nil {
		options = append(options, llms.WithMaxTokens(int(*maxTokens)))
	}

	// 使用 prompts 格式化和调用 LLM
	result, err := chatPrompt.FormatPrompt(map[string]any{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to format prompt: %v", err)
	}

	// 转换为 MessageContent 格式
	chatMessages := result.Messages()
	var llmMessages []llms.MessageContent
	for _, chatMsg := range chatMessages {
		llmMessages = append(llmMessages, llms.MessageContent{
			Role: chatMsg.GetType(),
			Parts: []llms.ContentPart{
				llms.TextPart(chatMsg.GetContent()),
			},
		})
	}

	// 调用 LLM
	resp, err := p.client.GenerateContent(ctx, llmMessages, options...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LLM call failed: %v", err)
	}

	// 提取响应
	if len(resp.Choices) == 0 {
		return nil, status.Error(codes.Internal, "no response from LLM")
	}

	choice := resp.Choices[0]
	if choice.Content == "" {
		return nil, status.Error(codes.Internal, "empty response from LLM")
	}

	// 估算 token 使用情况
	tokenUsage := EstimateTokenUsage(messages, choice.Content)

	return &ProcessResult{
		Content:    choice.Content,
		TokenUsage: tokenUsage,
	}, nil
}

// buildChatPrompt 构建使用 prompts 包装的聊天提示
func (p *Processor) buildChatPrompt(messages []*genaidemo.Message) prompts.ChatPromptTemplate {
	var promptMessages []prompts.MessageFormatter
	
	for _, msg := range messages {
		switch msg.Role {
		case genaidemo.Role_ROLE_SYSTEM:
			promptMessages = append(promptMessages, prompts.NewSystemMessagePromptTemplate(msg.Content, nil))
		case genaidemo.Role_ROLE_USER:
			promptMessages = append(promptMessages, prompts.NewHumanMessagePromptTemplate(msg.Content, nil))
		case genaidemo.Role_ROLE_ASSISTANT:
			promptMessages = append(promptMessages, prompts.NewAIMessagePromptTemplate(msg.Content, nil))
		}
	}
	
	return prompts.NewChatPromptTemplate(promptMessages)
}

// TokenUsage token 使用情况统计
type TokenUsage struct {
	InputTokens  int32
	OutputTokens int32
	TotalTokens  int32
}

// EstimateTokens 估算消息的 token 数量
func EstimateTokens(messages []*genaidemo.Message) int {
	totalTokens := 0
	for _, msg := range messages {
		// 简单估算: 每4个字符约等于1个token
		totalTokens += len(msg.Content) / 4
	}
	return totalTokens
}

// EstimateTokenUsage 估算 token 使用情况
func EstimateTokenUsage(messages []*genaidemo.Message, responseContent string) *TokenUsage {
	inputTokens := EstimateTokens(messages)
	outputTokens := len(responseContent) / 4
	
	return &TokenUsage{
		InputTokens:  int32(inputTokens),
		OutputTokens: int32(outputTokens),
		TotalTokens:  int32(inputTokens + outputTokens),
	}
}