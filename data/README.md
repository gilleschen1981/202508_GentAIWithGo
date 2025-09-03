# ChromaDB PDF Document Service

This directory contains a complete setup for embedding PDF documents into ChromaDB and serving them via a REST API.

## Quick Start

### 1. Setup (One-time)
```bash
# Virtual environment and dependencies are already installed
# Source directory is created automatically
```

### 2. Add PDF Documents
```bash
# Place your PDF files in the source/ directory
cp /path/to/your/pdfs/*.pdf source/
```

### 3. Embed Documents
```bash
# Make scripts executable
chmod +x embed_pdfs.sh start_chromadb.sh

# Embed all PDFs in source/ directory
./embed_pdfs.sh
```

### 4. Start ChromaDB Service
```bash
# Start the service on port 8000
./start_chromadb.sh
```

## Usage

### Embedding PDFs
```bash
# Basic usage
./embed_pdfs.sh

# Specify custom source directory
./embed_pdfs.sh --source /path/to/pdfs

# Reset database before embedding
./embed_pdfs.sh --reset

# Custom database location
./embed_pdfs.sh --db-path ./my_custom_db
```

### Starting the Service
```bash
# Basic usage (listens on 0.0.0.0:8000)
./start_chromadb.sh

# Custom port
./start_chromadb.sh --port 8080

# Custom host and port
./start_chromadb.sh --host localhost --port 3000

# Development mode with auto-reload
./start_chromadb.sh --reload
```

### Python Scripts (Direct Usage)
```bash
# Activate virtual environment first
source .venv/bin/activate

# Embed PDFs
python pdf_embedder.py --help
python pdf_embedder.py --source source --db-path ./chroma_db

# Start service
python chromadb_service.py --help
python chromadb_service.py --port 8000
```

## API Endpoints

Once the service is running, you can access:

- **API Documentation**: http://localhost:8000/docs
- **Health Check**: http://localhost:8000/health
- **Statistics**: http://localhost:8000/stats
- **Query Documents**: POST http://localhost:8000/query

### Example API Usage

#### Query Documents
```bash
curl -X POST "http://localhost:8000/query" \
     -H "Content-Type: application/json" \
     -d '{
       "query": "What is machine learning?",
       "n_results": 5,
       "include_metadata": true
     }'
```

#### Get Statistics
```bash
curl "http://localhost:8000/stats"
```

#### Health Check
```bash
curl "http://localhost:8000/health"
```

## Project Structure

```
data/
├── .venv/                 # Python virtual environment
├── source/                # Place PDF files here
├── chroma_db/            # ChromaDB storage (created automatically)
├── pdf_embedder.py       # PDF embedding script
├── chromadb_service.py   # ChromaDB REST API service
├── embed_pdfs.sh         # PDF embedding wrapper script
├── start_chromadb.sh     # Service startup script
└── README.md             # This file
```

## Features

### PDF Embedder (`pdf_embedder.py`)
- Extracts text from PDF files using PyPDF2
- Splits text into overlapping chunks for better search
- Uses sentence-transformers for embeddings
- Handles multiple PDFs automatically
- Provides detailed progress feedback

### ChromaDB Service (`chromadb_service.py`)
- REST API built with FastAPI
- CORS enabled for web access
- Automatic API documentation
- Health monitoring endpoints
- Query similarity search
- Collection statistics

### Shell Scripts
- **`embed_pdfs.sh`**: User-friendly PDF embedding with progress
- **`start_chromadb.sh`**: Service startup with configuration options
- Both scripts include help messages and error handling

## Configuration

### Environment Variables
You can set these environment variables to customize behavior:

```bash
export CHROMADB_PATH="./chroma_db"
export CHROMADB_PORT="8000"
export PDF_SOURCE_DIR="source"
```

### Embedding Model
The default embedding model is `all-MiniLM-L6-v2` (fast, good quality).
You can modify this in `pdf_embedder.py` if needed.

### Text Chunking
Default chunk size: 1000 characters with 100 character overlap.
Modify in `pdf_embedder.py` if you need different chunking.

## Troubleshooting

### No PDFs Found
```bash
# Make sure PDFs are in the source directory
ls -la source/
```

### Service Won't Start
```bash
# Check if port is already in use
lsof -i :8000

# Try a different port
./start_chromadb.sh --port 8080
```

### Import Errors
```bash
# Make sure virtual environment is activated
source .venv/bin/activate

# Reinstall dependencies if needed
pip install chromadb PyPDF2 sentence-transformers
```

### No Search Results
```bash
# Check if documents were embedded
curl "http://localhost:8000/stats"

# Re-embed with reset flag
./embed_pdfs.sh --reset
```

## Advanced Usage

### Custom Embedding Model
Edit `pdf_embedder.py` and change the model:
```python
self.model = SentenceTransformer('all-mpnet-base-v2')  # Better quality, slower
```

### Multiple Collections
You can create multiple collections for different document types:
```bash
./embed_pdfs.sh --db-path ./legal_docs
./start_chromadb.sh --db-path ./legal_docs --collection legal_documents --port 8001
```

### Integration with Go Service
The ChromaDB service can be easily integrated with the main Go service by making HTTP requests to the query endpoint.

## Dependencies

- Python 3.9+
- chromadb
- PyPDF2
- sentence-transformers
- fastapi
- uvicorn
- pydantic

All dependencies are automatically installed in the virtual environment.