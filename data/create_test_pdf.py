#!/usr/bin/env python3
"""
Create a simple test PDF for testing the embedding system
"""

try:
    from reportlab.pdfgen import canvas
    from reportlab.lib.pagesizes import letter
    REPORTLAB_AVAILABLE = True
except ImportError:
    REPORTLAB_AVAILABLE = False

import os

def create_test_pdf(filename="test_document.pdf"):
    """Create a simple test PDF with sample content"""
    
    if not REPORTLAB_AVAILABLE:
        raise ImportError("reportlab not available")
    
    # Ensure source directory exists
    os.makedirs("source", exist_ok=True)
    filepath = os.path.join("source", filename)
    
    # Create PDF
    c = canvas.Canvas(filepath, pagesize=letter)
    width, height = letter
    
    # Add content
    c.setFont("Helvetica-Bold", 16)
    c.drawString(50, height - 50, "Test Document for ChromaDB")
    
    c.setFont("Helvetica", 12)
    y_pos = height - 100
    
    content = [
        "This is a test document for the ChromaDB embedding system.",
        "",
        "Machine Learning Overview:",
        "Machine learning is a subset of artificial intelligence that focuses on",
        "algorithms that can learn and make decisions from data without being",
        "explicitly programmed for every specific task.",
        "",
        "Key Concepts:",
        "- Supervised Learning: Learning from labeled examples",
        "- Unsupervised Learning: Finding patterns in unlabeled data", 
        "- Deep Learning: Neural networks with multiple layers",
        "- Natural Language Processing: Understanding and generating text",
        "",
        "Applications:",
        "Machine learning is used in many areas including:",
        "• Image recognition and computer vision",
        "• Speech recognition and natural language processing",
        "• Recommendation systems",
        "• Autonomous vehicles",
        "• Medical diagnosis and drug discovery",
        "",
        "Future Trends:",
        "The field continues to evolve with advances in:",
        "- Large language models and generative AI",
        "- Federated learning for privacy-preserving ML",
        "- AutoML for automated model development",
        "- Edge AI for on-device intelligence"
    ]
    
    for line in content:
        c.drawString(50, y_pos, line)
        y_pos -= 15
        if y_pos < 50:  # Start new page if needed
            c.showPage()
            c.setFont("Helvetica", 12)
            y_pos = height - 50
    
    c.save()
    print(f"Test PDF created: {filepath}")

if __name__ == "__main__":
    try:
        create_test_pdf()
    except ImportError:
        print("reportlab not installed. Creating a simple text file instead...")
        # Create a simple text file that we can manually describe as PDF content
        os.makedirs("source", exist_ok=True)
        with open("source/README_test.txt", "w") as f:
            f.write("""
This is test content for ChromaDB embedding.

Machine Learning Overview:
Machine learning is a subset of artificial intelligence that focuses on
algorithms that can learn and make decisions from data without being
explicitly programmed for every specific task.

Key Concepts:
- Supervised Learning: Learning from labeled examples
- Unsupervised Learning: Finding patterns in unlabeled data
- Deep Learning: Neural networks with multiple layers
- Natural Language Processing: Understanding and generating text

Applications:
Machine learning is used in many areas including:
• Image recognition and computer vision
• Speech recognition and natural language processing  
• Recommendation systems
• Autonomous vehicles
• Medical diagnosis and drug discovery
""")
        print("Created test text file: source/README_test.txt")
        print("Note: For full testing, place some actual PDF files in the source/ directory")