#!/usr/bin/env python3
"""
Test script for ChromaDB setup
Tests the embedding and service functionality with sample data
"""

import os
import sys
import time
import requests
import subprocess
from pathlib import Path
import chromadb
from sentence_transformers import SentenceTransformer

def test_chromadb_installation():
    """Test if ChromaDB is properly installed"""
    print("🔧 Testing ChromaDB installation...")
    try:
        import chromadb
        print("✅ ChromaDB imported successfully")
        
        # Test basic functionality
        client = chromadb.Client()
        collection = client.create_collection("test_collection")
        collection.add(
            documents=["This is a test document"],
            metadatas=[{"source": "test"}],
            ids=["test1"]
        )
        results = collection.query(query_texts=["test"], n_results=1)
        print("✅ ChromaDB basic functionality works")
        return True
    except Exception as e:
        print(f"❌ ChromaDB test failed: {e}")
        return False

def test_sentence_transformers():
    """Test if sentence transformers work"""
    print("\n🤖 Testing sentence transformers...")
    try:
        model = SentenceTransformer('all-MiniLM-L6-v2')
        embeddings = model.encode(["This is a test sentence"])
        print(f"✅ Sentence transformers work. Embedding shape: {embeddings.shape}")
        return True
    except Exception as e:
        print(f"❌ Sentence transformers test failed: {e}")
        return False

def create_test_data():
    """Create some test data for embedding"""
    print("\n📄 Creating test data...")
    
    os.makedirs("source", exist_ok=True)
    
    # Create multiple test documents
    test_docs = {
        "machine_learning.txt": """
Machine Learning Overview
Machine learning is a subset of artificial intelligence that enables computers to learn from data.

Key Types:
1. Supervised Learning - learns from labeled examples
2. Unsupervised Learning - finds patterns in unlabeled data  
3. Reinforcement Learning - learns through trial and error

Popular algorithms include decision trees, neural networks, and support vector machines.
        """,
        
        "data_science.txt": """
Data Science Fundamentals
Data science combines statistics, programming, and domain expertise to extract insights from data.

Process includes:
- Data collection and cleaning
- Exploratory data analysis
- Statistical modeling
- Machine learning application
- Results interpretation and communication

Tools commonly used: Python, R, SQL, Tableau, Jupyter notebooks.
        """,
        
        "ai_applications.txt": """
AI Applications in Industry
Artificial intelligence is transforming various industries:

Healthcare: Medical diagnosis, drug discovery, personalized treatment
Finance: Fraud detection, algorithmic trading, risk assessment
Transportation: Autonomous vehicles, route optimization
Retail: Recommendation systems, inventory management
Entertainment: Content recommendation, game AI

The impact continues to grow as technology advances.
        """
    }
    
    for filename, content in test_docs.items():
        filepath = Path("source") / filename
        with open(filepath, 'w') as f:
            f.write(content.strip())
        print(f"✅ Created: {filepath}")
    
    return len(test_docs)

def test_embedding_script():
    """Test the PDF embedder script with text files"""
    print("\n📊 Testing embedding script...")
    
    try:
        # Run the embedder (it should work with any files)
        result = subprocess.run([
            sys.executable, "pdf_embedder.py", 
            "--source", "source",
            "--db-path", "./test_chroma_db"
        ], capture_output=True, text=True, timeout=60)
        
        if result.returncode == 0:
            print("✅ Embedding script completed successfully")
            print("Output:", result.stdout.split('\n')[-3:])  # Last few lines
            return True
        else:
            print(f"❌ Embedding script failed: {result.stderr}")
            return False
    except subprocess.TimeoutExpired:
        print("❌ Embedding script timed out")
        return False
    except Exception as e:
        print(f"❌ Error running embedding script: {e}")
        return False

def test_chromadb_query():
    """Test querying the embedded data directly"""
    print("\n🔍 Testing ChromaDB queries...")
    
    try:
        client = chromadb.PersistentClient(path="./test_chroma_db")
        collection = client.get_collection("pdf_documents")
        
        # Test query
        results = collection.query(
            query_texts=["machine learning algorithms"],
            n_results=3
        )
        
        print(f"✅ Query successful. Found {len(results['documents'][0])} results")
        print("Sample result:", results['documents'][0][0][:100] + "..." if results['documents'][0] else "No results")
        return True
    except Exception as e:
        print(f"❌ ChromaDB query failed: {e}")
        return False

def test_service_startup():
    """Test if the service can start (but don't leave it running)"""
    print("\n🚀 Testing service startup...")
    
    try:
        # Start the service in background
        process = subprocess.Popen([
            sys.executable, "chromadb_service.py",
            "--db-path", "./test_chroma_db",
            "--port", "8001"  # Use different port to avoid conflicts
        ], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        
        # Wait a moment for startup
        time.sleep(3)
        
        # Check if service is responding
        try:
            response = requests.get("http://localhost:8001/health", timeout=5)
            if response.status_code == 200:
                print("✅ Service started and responding")
                
                # Test a query
                query_response = requests.post("http://localhost:8001/query", 
                    json={"query": "machine learning", "n_results": 2},
                    timeout=5
                )
                if query_response.status_code == 200:
                    data = query_response.json()
                    print(f"✅ Service query successful. Found {len(data['documents'])} results")
                else:
                    print(f"⚠️  Service query failed: {query_response.status_code}")
                
                success = True
            else:
                print(f"❌ Service not responding properly: {response.status_code}")
                success = False
        except requests.exceptions.RequestException as e:
            print(f"❌ Cannot connect to service: {e}")
            success = False
        
        # Stop the service
        process.terminate()
        process.wait(timeout=5)
        print("🛑 Service stopped")
        
        return success
        
    except Exception as e:
        print(f"❌ Service test failed: {e}")
        return False

def cleanup_test_data():
    """Clean up test data"""
    print("\n🧹 Cleaning up test data...")
    
    import shutil
    
    # Remove test database
    if os.path.exists("./test_chroma_db"):
        shutil.rmtree("./test_chroma_db")
        print("✅ Removed test database")
    
    # Optionally remove test files
    # for filename in ["machine_learning.txt", "data_science.txt", "ai_applications.txt"]:
    #     filepath = Path("source") / filename
    #     if filepath.exists():
    #         filepath.unlink()
    #         print(f"✅ Removed: {filepath}")

def main():
    """Run all tests"""
    print("=" * 60)
    print("🧪 ChromaDB Setup Test Suite")
    print("=" * 60)
    
    tests = [
        ("ChromaDB Installation", test_chromadb_installation),
        ("Sentence Transformers", test_sentence_transformers),
        ("Test Data Creation", lambda: create_test_data() > 0),
        ("Embedding Script", test_embedding_script),
        ("ChromaDB Queries", test_chromadb_query),
        ("Service Startup", test_service_startup),
    ]
    
    results = []
    
    for test_name, test_func in tests:
        print(f"\n{'=' * 20} {test_name} {'=' * 20}")
        try:
            success = test_func()
            results.append((test_name, success))
        except Exception as e:
            print(f"❌ Test failed with exception: {e}")
            results.append((test_name, False))
    
    # Summary
    print("\n" + "=" * 60)
    print("📋 TEST SUMMARY")
    print("=" * 60)
    
    passed = 0
    for test_name, success in results:
        status = "✅ PASS" if success else "❌ FAIL"
        print(f"{status}: {test_name}")
        if success:
            passed += 1
    
    print(f"\nPassed: {passed}/{len(results)} tests")
    
    if passed == len(results):
        print("\n🎉 All tests passed! ChromaDB setup is working correctly.")
        print("\nNext steps:")
        print("1. Place PDF files in the source/ directory")
        print("2. Run: ./embed_pdfs.sh")
        print("3. Run: ./start_chromadb.sh")
        print("4. Access API at: http://localhost:8000/docs")
    else:
        print(f"\n⚠️  Some tests failed. Please check the errors above.")
    
    # Cleanup
    cleanup_test_data()

if __name__ == "__main__":
    main()