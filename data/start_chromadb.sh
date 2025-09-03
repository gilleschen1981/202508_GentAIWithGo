#!/bin/bash
# ChromaDB Service Startup Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== ChromaDB Service Startup ===${NC}"

# Check if virtual environment exists
if [ ! -d ".venv" ]; then
    echo -e "${RED}Error: Virtual environment not found. Please run setup first.${NC}"
    exit 1
fi

# Activate virtual environment
echo -e "${YELLOW}Activating virtual environment...${NC}"
source .venv/bin/activate

# Default values
DB_PATH="./chroma_db"
COLLECTION="pdf_documents"
HOST="0.0.0.0"
PORT="8000"
RELOAD_FLAG=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --db-path|-d)
            DB_PATH="$2"
            shift 2
            ;;
        --collection|-c)
            COLLECTION="$2"
            shift 2
            ;;
        --host)
            HOST="$2"
            shift 2
            ;;
        --port|-p)
            PORT="$2"
            shift 2
            ;;
        --reload)
            RELOAD_FLAG="--reload"
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --db-path, -d PATH     ChromaDB storage path (default: ./chroma_db)"
            echo "  --collection, -c NAME  Collection name (default: pdf_documents)"
            echo "  --host HOST           Host to bind to (default: 0.0.0.0)"
            echo "  --port, -p PORT       Port to bind to (default: 8000)"
            echo "  --reload              Enable auto-reload for development"
            echo "  --help, -h            Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Check if ChromaDB exists
if [ ! -d "$DB_PATH" ]; then
    echo -e "${YELLOW}Warning: ChromaDB not found at $DB_PATH${NC}"
    echo -e "${YELLOW}Make sure to embed some PDFs first using: python pdf_embedder.py${NC}"
fi

echo -e "${GREEN}Starting ChromaDB service...${NC}"
echo "Database path: $DB_PATH"
echo "Collection: $COLLECTION"
echo "Host: $HOST"
echo "Port: $PORT"
echo ""
echo -e "${BLUE}API will be available at: http://$HOST:$PORT${NC}"
echo -e "${BLUE}API documentation: http://$HOST:$PORT/docs${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the service${NC}"
echo ""

# Start the service
python chromadb_service.py \
    --db-path "$DB_PATH" \
    --collection "$COLLECTION" \
    --host "$HOST" \
    --port "$PORT" \
    $RELOAD_FLAG