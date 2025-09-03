#!/bin/bash
# PDF Embedding Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== PDF Embedding Script ===${NC}"

# Check if virtual environment exists
if [ ! -d ".venv" ]; then
    echo -e "${RED}Error: Virtual environment not found. Please run setup first.${NC}"
    exit 1
fi

# Activate virtual environment
echo -e "${YELLOW}Activating virtual environment...${NC}"
source .venv/bin/activate

# Default values
SOURCE_DIR="source"
DB_PATH="./chroma_db"
RESET_FLAG=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --source|-s)
            SOURCE_DIR="$2"
            shift 2
            ;;
        --db-path|-d)
            DB_PATH="$2"
            shift 2
            ;;
        --reset|-r)
            RESET_FLAG="--reset"
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo "Options:"
            echo "  --source, -s DIR      Source directory containing PDF files (default: source)"
            echo "  --db-path, -d PATH    ChromaDB storage path (default: ./chroma_db)"
            echo "  --reset, -r           Reset the database before processing"
            echo "  --help, -h            Show this help message"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Check if source directory exists
if [ ! -d "$SOURCE_DIR" ]; then
    echo -e "${YELLOW}Source directory $SOURCE_DIR does not exist. Creating it...${NC}"
    mkdir -p "$SOURCE_DIR"
    echo -e "${YELLOW}Please place your PDF files in the $SOURCE_DIR directory and run this script again.${NC}"
    exit 0
fi

# Count PDF files
PDF_COUNT=$(find "$SOURCE_DIR" -name "*.pdf" | wc -l)
if [ "$PDF_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}No PDF files found in $SOURCE_DIR directory.${NC}"
    echo -e "${YELLOW}Please place some PDF files there and try again.${NC}"
    exit 0
fi

echo -e "${GREEN}Found $PDF_COUNT PDF file(s) in $SOURCE_DIR${NC}"

# Show reset warning
if [ ! -z "$RESET_FLAG" ]; then
    echo -e "${RED}WARNING: This will reset the existing database!${NC}"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 0
    fi
fi

echo -e "${BLUE}Starting PDF embedding process...${NC}"
echo "Source directory: $SOURCE_DIR"
echo "Database path: $DB_PATH"
echo ""

# Run the embedding script
python pdf_embedder.py \
    --source "$SOURCE_DIR" \
    --db-path "$DB_PATH" \
    $RESET_FLAG

echo ""
echo -e "${GREEN}PDF embedding completed!${NC}"
echo -e "${BLUE}You can now start the ChromaDB service with: ./start_chromadb.sh${NC}"