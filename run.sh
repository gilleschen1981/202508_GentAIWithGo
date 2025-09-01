#!/bin/bash

echo "🚀 GenAI Foundation Demo - Starting..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Build the service
echo "🔨 Building gRPC service..."
cd service && go build -o ../genai-service . && cd ..

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

# Build test client (skipped due to package conflicts)
echo "🔨 Building test client... (skipped - using frontend for testing)"
# cd client && go build -o ../test-client . && cd ..
# 
# if [ $? -ne 0 ]; then
#     echo "❌ Test client build failed"
#     exit 1
# fi

echo "✅ Build successful!"
echo ""

# Start the service
echo "🚀 Starting gRPC service on port 50051..."
echo "📝 Project configured for: $(grep DefaultProjectID config.go | cut -d'"' -f2)"
echo ""

./genai-service