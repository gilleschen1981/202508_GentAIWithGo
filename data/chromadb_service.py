#!/usr/bin/env python3
"""
ChromaDB Service Server
Starts ChromaDB as a service with HTTP API access
"""

import os
import sys
import argparse
import signal
from pathlib import Path
import chromadb
from chromadb.config import Settings
import uvicorn
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional, Dict, Any
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class QueryRequest(BaseModel):
    query: str
    n_results: int = 5
    include_metadata: bool = True

class QueryResponse(BaseModel):
    documents: List[str]
    metadatas: Optional[List[Dict[str, Any]]] = None
    distances: Optional[List[float]] = None
    ids: List[str]

class ChromaDBService:
    def __init__(self, db_path: str = "./chroma_db", collection_name: str = "pdf_documents"):
        self.db_path = Path(db_path)
        self.collection_name = collection_name
        self.client = None
        self.collection = None
        
    def initialize(self):
        """Initialize ChromaDB client and collection"""
        try:
            logger.info(f"Initializing ChromaDB from: {self.db_path}")
            self.client = chromadb.PersistentClient(path=str(self.db_path))
            
            # Get existing collection
            try:
                self.collection = self.client.get_collection(name=self.collection_name)
                count = self.collection.count()
                logger.info(f"Connected to collection '{self.collection_name}' with {count} documents")
            except Exception as e:
                logger.warning(f"Collection '{self.collection_name}' not found: {e}")
                self.collection = None
                
        except Exception as e:
            logger.error(f"Failed to initialize ChromaDB: {e}")
            raise
    
    def query_documents(self, query: str, n_results: int = 5, include_metadata: bool = True) -> Dict:
        """Query the document collection"""
        if not self.collection:
            raise HTTPException(status_code=404, detail="No collection available. Please embed some documents first.")
        
        try:
            include = ["documents", "distances", "metadatas"] if include_metadata else ["documents", "distances"]
            
            results = self.collection.query(
                query_texts=[query],
                n_results=n_results,
                include=include
            )
            
            return {
                "documents": results["documents"][0] if results["documents"] else [],
                "metadatas": results["metadatas"][0] if include_metadata and results.get("metadatas") else None,
                "distances": results["distances"][0] if results["distances"] else [],
                "ids": results["ids"][0] if results["ids"] else []
            }
        except Exception as e:
            logger.error(f"Query failed: {e}")
            raise HTTPException(status_code=500, detail=f"Query failed: {str(e)}")
    
    def get_stats(self) -> Dict:
        """Get collection statistics"""
        if not self.collection:
            return {"status": "no_collection", "message": "No collection available"}
        
        try:
            count = self.collection.count()
            return {
                "status": "ready",
                "collection_name": self.collection_name,
                "document_count": count,
                "db_path": str(self.db_path)
            }
        except Exception as e:
            return {"status": "error", "message": str(e)}

# Global service instance
service = None

# FastAPI app
app = FastAPI(
    title="ChromaDB Document Search API",
    description="Query embedded PDF documents using ChromaDB",
    version="1.0.0"
)

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
async def startup_event():
    global service
    logger.info("Starting ChromaDB service...")
    # Initialize service here
    if service:
        logger.info("Service already initialized")
    
@app.get("/")
async def root():
    return {"message": "ChromaDB Document Search API", "status": "running"}

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    global service
    if service:
        stats = service.get_stats()
        return {"health": "ok", "service": stats}
    return {"health": "ok", "service": {"status": "not_initialized"}}

@app.get("/stats")
async def get_stats():
    """Get collection statistics"""
    global service
    if not service:
        raise HTTPException(status_code=503, detail="Service not initialized")
    
    return service.get_stats()

@app.post("/query", response_model=QueryResponse)
async def query_documents(request: QueryRequest):
    """Query documents in the collection"""
    global service
    if not service:
        raise HTTPException(status_code=503, detail="Service not initialized")
    
    result = service.query_documents(
        query=request.query,
        n_results=request.n_results,
        include_metadata=request.include_metadata
    )
    
    return QueryResponse(**result)

@app.get("/collections")
async def list_collections():
    """List available collections"""
    global service
    if not service or not service.client:
        raise HTTPException(status_code=503, detail="Service not initialized")
    
    try:
        collections = service.client.list_collections()
        return {"collections": [{"name": col.name, "metadata": col.metadata} for col in collections]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Failed to list collections: {str(e)}")

def signal_handler(sig, frame):
    """Handle shutdown signals"""
    logger.info("Shutting down ChromaDB service...")
    sys.exit(0)

def main():
    parser = argparse.ArgumentParser(description="Start ChromaDB service")
    parser.add_argument("--db-path", "-d", default="./chroma_db",
                       help="ChromaDB storage path (default: ./chroma_db)")
    parser.add_argument("--collection", "-c", default="pdf_documents",
                       help="Collection name (default: pdf_documents)")
    parser.add_argument("--host", default="0.0.0.0",
                       help="Host to bind to (default: 0.0.0.0)")
    parser.add_argument("--port", "-p", type=int, default=8000,
                       help="Port to bind to (default: 8000)")
    parser.add_argument("--reload", action="store_true",
                       help="Enable auto-reload for development")
    
    args = parser.parse_args()
    
    # Set up signal handlers
    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)
    
    # Initialize service
    global service
    service = ChromaDBService(args.db_path, args.collection)
    
    try:
        service.initialize()
        logger.info("ChromaDB service initialized successfully")
    except Exception as e:
        logger.error(f"Failed to initialize service: {e}")
        sys.exit(1)
    
    # Start the server
    logger.info(f"Starting server on {args.host}:{args.port}")
    logger.info(f"API documentation available at: http://{args.host}:{args.port}/docs")
    
    # Create the app module name dynamically to avoid reload issues
    module_name = __name__ + ":app" if __name__ == "__main__" else "chromadb_service:app"
    
    uvicorn.run(
        module_name,
        host=args.host,
        port=args.port,
        reload=args.reload,
        log_level="info"
    )

if __name__ == "__main__":
    main()