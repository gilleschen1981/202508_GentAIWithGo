package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/example/genai-foundation-demo"
	"google.golang.org/grpc"
)

const (
	serviceName = "genai-chat-service"
	grpcPort    = "50051"
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

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	genaidemo.RegisterChatServiceServer(grpcServer, handler)

	log.Printf("🚀 gRPC server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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