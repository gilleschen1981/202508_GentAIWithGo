#!/usr/bin/env python3
"""
PDF Embedding Script for ChromaDB
Processes all PDF files in the source directory and embeds them into ChromaDB
"""

import os
import sys
from pathlib import Path
import PyPDF2
import chromadb
from chromadb.config import Settings
from sentence_transformers import SentenceTransformer
import argparse
from typing import List, Dict

class PDFEmbedder:
    def __init__(self, source_dir: str = "source", db_path: str = "./chroma_db"):
        self.source_dir = Path(source_dir)
        self.db_path = Path(db_path)
        self.client = None
        self.collection = None
        self.model = None
        
    def initialize_chromadb(self):
        """Initialize ChromaDB client and collection"""
        print("Initializing ChromaDB...")
        self.client = chromadb.PersistentClient(path=str(self.db_path))
        
        # Create or get collection
        self.collection = self.client.get_or_create_collection(
            name="pdf_documents",
            metadata={"description": "PDF document embeddings"}
        )
        print(f"ChromaDB initialized at: {self.db_path}")
        
    def initialize_model(self):
        """Initialize sentence transformer model"""
        print("Loading embedding model...")
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        print("Embedding model loaded successfully")
        
    def extract_text_from_pdf(self, pdf_path: Path) -> str:
        """Extract text content from PDF file"""
        text = ""
        try:
            with open(pdf_path, 'rb') as file:
                pdf_reader = PyPDF2.PdfReader(file)
                for page_num, page in enumerate(pdf_reader.pages):
                    try:
                        page_text = page.extract_text()
                        if page_text.strip():
                            text += f"\n--- Page {page_num + 1} ---\n"
                            text += page_text
                    except Exception as e:
                        print(f"Warning: Could not extract text from page {page_num + 1} of {pdf_path}: {e}")
                        continue
        except Exception as e:
            print(f"Error reading PDF {pdf_path}: {e}")
            return ""
        
        return text.strip()
    
    def chunk_text(self, text: str, chunk_size: int = 1000, overlap: int = 100) -> List[str]:
        """Split text into overlapping chunks"""
        if len(text) <= chunk_size:
            return [text]
        
        chunks = []
        start = 0
        
        while start < len(text):
            end = start + chunk_size
            
            # Try to find a good breaking point (sentence or paragraph end)
            if end < len(text):
                # Look for sentence endings within the last 200 characters
                search_start = max(end - 200, start)
                for i in range(end, search_start, -1):
                    if text[i] in '.!?\n':
                        end = i + 1
                        break
            
            chunk = text[start:end].strip()
            if chunk:
                chunks.append(chunk)
            
            start = end - overlap
            if start >= len(text):
                break
                
        return chunks
    
    def embed_pdf(self, pdf_path: Path) -> bool:
        """Process and embed a single PDF file"""
        print(f"Processing: {pdf_path.name}")
        
        # Extract text
        text = self.extract_text_from_pdf(pdf_path)
        if not text:
            print(f"Warning: No text extracted from {pdf_path.name}")
            return False
        
        # Split into chunks
        chunks = self.chunk_text(text)
        print(f"  Split into {len(chunks)} chunks")
        
        # Generate embeddings
        embeddings = self.model.encode(chunks)
        
        # Prepare metadata and IDs
        ids = [f"{pdf_path.stem}_chunk_{i}" for i in range(len(chunks))]
        metadatas = [
            {
                "source": str(pdf_path),
                "filename": pdf_path.name,
                "chunk_index": i,
                "total_chunks": len(chunks)
            }
            for i in range(len(chunks))
        ]
        
        # Add to ChromaDB
        try:
            self.collection.add(
                embeddings=embeddings.tolist(),
                documents=chunks,
                metadatas=metadatas,
                ids=ids
            )
            print(f"  Successfully embedded {len(chunks)} chunks")
            return True
        except Exception as e:
            print(f"Error adding to ChromaDB: {e}")
            return False
    
    def process_all_pdfs(self) -> Dict[str, int]:
        """Process all PDF files in the source directory"""
        if not self.source_dir.exists():
            print(f"Source directory {self.source_dir} does not exist!")
            return {"processed": 0, "failed": 0}
        
        pdf_files = list(self.source_dir.glob("*.pdf"))
        if not pdf_files:
            print(f"No PDF files found in {self.source_dir}")
            return {"processed": 0, "failed": 0}
        
        print(f"Found {len(pdf_files)} PDF files to process")
        
        processed = 0
        failed = 0
        
        for pdf_file in pdf_files:
            try:
                if self.embed_pdf(pdf_file):
                    processed += 1
                else:
                    failed += 1
            except Exception as e:
                print(f"Error processing {pdf_file.name}: {e}")
                failed += 1
        
        return {"processed": processed, "failed": failed}
    
    def get_collection_stats(self) -> Dict:
        """Get statistics about the collection"""
        if not self.collection:
            return {}
        
        try:
            count = self.collection.count()
            return {"total_documents": count}
        except Exception as e:
            print(f"Error getting collection stats: {e}")
            return {}

def main():
    parser = argparse.ArgumentParser(description="Embed PDF files into ChromaDB")
    parser.add_argument("--source", "-s", default="source", 
                       help="Source directory containing PDF files (default: source)")
    parser.add_argument("--db-path", "-d", default="./chroma_db",
                       help="ChromaDB storage path (default: ./chroma_db)")
    parser.add_argument("--reset", "-r", action="store_true",
                       help="Reset the database before processing")
    
    args = parser.parse_args()
    
    # Initialize embedder
    embedder = PDFEmbedder(args.source, args.db_path)
    
    try:
        # Initialize components
        embedder.initialize_chromadb()
        embedder.initialize_model()
        
        # Reset database if requested
        if args.reset:
            print("Resetting database...")
            embedder.collection.delete()
            embedder.collection = embedder.client.get_or_create_collection(
                name="pdf_documents",
                metadata={"description": "PDF document embeddings"}
            )
            print("Database reset complete")
        
        # Process PDFs
        results = embedder.process_all_pdfs()
        
        # Print results
        print("\n" + "="*50)
        print("PROCESSING COMPLETE")
        print("="*50)
        print(f"Successfully processed: {results['processed']} PDFs")
        print(f"Failed to process: {results['failed']} PDFs")
        
        # Get collection stats
        stats = embedder.get_collection_stats()
        if stats:
            print(f"Total document chunks in database: {stats['total_documents']}")
        
        print(f"ChromaDB location: {embedder.db_path}")
        
    except Exception as e:
        print(f"Fatal error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()