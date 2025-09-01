package main

// VertexAI 配置常量
// 请根据你的实际情况修改这些配置
const (
	// 你的 Google Cloud 项目 ID
	DefaultProjectID = "prj-dscore-d2-qeby"

	// VertexAI 服务区域
	// 可选项: us-central1, us-east1, us-west1, europe-west1, europe-west4, asia-southeast1
	DefaultLocation = "us-central1"

	// VertexAI 模型名称
	// 推荐选项:
	// - "gemini-1.5-flash"  (快速、经济、适合大多数场景)
	// - "gemini-1.5-pro"    (功能最强、但成本较高)  
	// - "gemini-1.0-pro"    (稳定版本)
	DefaultModelName = "gemini-1.5-flash"
)

// 模型配置说明
// gemini-1.5-flash:
//   - 速度最快
//   - 成本最低
//   - 适合日常对话、文本生成
//   - 支持多模态输入
//
// gemini-1.5-pro:
//   - 功能最强大
//   - 推理能力最好
//   - 适合复杂任务、代码生成
//   - 成本较高
//
// gemini-1.0-pro:
//   - 稳定可靠
//   - 平衡的性能和成本
//   - 适合生产环境
