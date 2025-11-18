#!/bin/bash
# newproject.sh - Create and initialize a new project in gobash

PROJECT_NAME=$1
PROJECT_TYPE=${2:-go}  # Default to Go

if [ -z "$PROJECT_NAME" ]; then
    echo "Usage: ./newproject.sh PROJECT_NAME [TYPE]"
    echo ""
    echo "Types: go, node, python, rust, simple"
    echo "Example: ./newproject.sh myapp go"
    exit 1
fi

echo "🚀 Creating project: $PROJECT_NAME ($PROJECT_TYPE)"
echo "=================================================="
echo ""

# Create base directory
mkdir -p "$PROJECT_NAME"
cd "$PROJECT_NAME" || exit 1

case $PROJECT_TYPE in
    go)
        echo "→ Setting up Go project..."
        mkdir -p cmd pkg internal
        
        cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello from $PROJECT_NAME!")
}
EOF
        
        go mod init "$PROJECT_NAME"
        echo "✓ Go module initialized"
        ;;
        
    node)
        echo "→ Setting up Node.js project..."
        mkdir -p src tests
        
        cat > package.json << EOF
{
  "name": "$PROJECT_NAME",
  "version": "1.0.0",
  "main": "src/index.js",
  "scripts": {
    "start": "node src/index.js",
    "test": "echo \"No tests yet\""
  }
}
EOF
        
        cat > src/index.js << 'EOF'
console.log('Hello from $PROJECT_NAME!');
EOF
        
        echo "✓ Node.js project initialized"
        ;;
        
    python)
        echo "→ Setting up Python project..."
        mkdir -p src tests
        
        cat > src/__init__.py << 'EOF'
"""$PROJECT_NAME package"""
__version__ = "0.1.0"
EOF
        
        cat > requirements.txt << 'EOF'
# Add your dependencies here
pytest
EOF
        
        cat > README.md << EOF
# $PROJECT_NAME

Python project created with gobash

## Setup
\`\`\`bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
\`\`\`
EOF
        
        echo "✓ Python project initialized"
        ;;
        
    rust)
        echo "→ Setting up Rust project..."
        cargo init --name "$PROJECT_NAME" .
        echo "✓ Rust project initialized"
        ;;
        
    simple)
        echo "→ Setting up simple project..."
        mkdir -p src docs tests
        ;;
        
    *)
        echo "✗ Unknown project type: $PROJECT_TYPE"
        exit 1
        ;;
esac

# Common files for all projects
cat > README.md << EOF
# $PROJECT_NAME

Created: $(date +%Y-%m-%d)
Type: $PROJECT_TYPE

## Description
Add your project description here.

## Getting Started
Add setup instructions here.
EOF

cat > .gitignore << 'EOF'
# Common
.DS_Store
*.log
*.swp

# Build outputs
/bin
/dist
/build

# Dependencies
node_modules/
venv/
target/
EOF

# Initialize git if not already in a repo
if [ ! -d .git ]; then
    git init
    git add .
    git commit -m "Initial commit: $PROJECT_NAME ($PROJECT_TYPE)"
    echo "✓ Git repository initialized"
fi

echo ""
echo "✅ Project created successfully!"
echo ""
echo "📁 Structure:"
ls -la

if [ "$GOBASH_SHELL" = "1" ]; then
    echo ""
    echo "💡 You're now in: $(pwd)"
    echo "   Try: name \"working on $PROJECT_NAME\""
else
    echo ""
    echo "💡 Navigate to: cd $PROJECT_NAME"
    echo "   Then start gobash and run: name \"working on $PROJECT_NAME\""
fi
