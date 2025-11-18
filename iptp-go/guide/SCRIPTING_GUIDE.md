# 📜 Writing Custom Scripts for gobash

## Overview

gobash scripts are regular bash/shell scripts that can **interact with gobash's state and context**. When you run a script inside gobash, it has access to special environment variables and can use gobash commands.

## 🌟 The gobash Environment

When a script runs inside gobash, these environment variables are automatically available:

```bash
GOBASH_SHELL=1              # Indicates running inside gobash
GOBASH_CWD=/current/path    # Current working directory
GOBASH_PROCESS=process_name # Current process name (if set)
```

## 📝 Example 1: Simple Status Script

Create `status.sh`:

```bash
#!/bin/bash
# status.sh - Show gobash context

if [ "$GOBASH_SHELL" = "1" ]; then
    echo "✓ Running inside gobash"
    echo "  Process: ${GOBASH_PROCESS:-unnamed}"
    echo "  Directory: $GOBASH_CWD"
    echo "  Current files:"
    ls -lh
else
    echo "✗ Not running in gobash context"
    echo "  Start gobash and run: ./status.sh"
fi
```

Usage:
```bash
[IPTP-1] myproject$ chmod +x status.sh
[IPTP-1] myproject$ ./status.sh
✓ Running inside gobash
  Process: IPTP-1
  Directory: /Users/you/myproject
  Current files:
  ...
```

## 📝 Example 2: Project Setup Script

Create `setup-project.sh`:

```bash
#!/bin/bash
# setup-project.sh - Initialize a new project

PROJECT_NAME=$1

if [ -z "$PROJECT_NAME" ]; then
    echo "Usage: ./setup-project.sh PROJECT_NAME"
    exit 1
fi

echo "🚀 Setting up project: $PROJECT_NAME"

# Create directory structure
mkdir -p "$PROJECT_NAME"/{src,tests,docs}
cd "$PROJECT_NAME"

# Create initial files
touch README.md
echo "# $PROJECT_NAME" > README.md
touch src/main.go
touch tests/.gitkeep

# If in gobash, save this location
if [ "$GOBASH_SHELL" = "1" ]; then
    echo "✓ Created project structure"
    echo "  You're now in: $(pwd)"
    echo ""
    echo "💡 Tip: Use 'name \"working on $PROJECT_NAME\"' to name this session"
else
    echo "✓ Project created at: $(pwd)"
fi

ls -la
```

Usage:
```bash
[IPTP-1] ~$ ./setup-project.sh myapp
🚀 Setting up project: myapp
✓ Created project structure
  You're now in: /Users/you/myapp
💡 Tip: Use 'name "working on myapp"' to name this session

[IPTP-1] myapp$ name "working on myapp"
✓ Shell named: myapp
[myapp] myapp$ 
```

## 📝 Example 3: Smart Search & Jump

Create `find-and-go.sh`:

```bash
#!/bin/bash
# find-and-go.sh - Find a directory and navigate to it

SEARCH_TERM=$1

if [ -z "$SEARCH_TERM" ]; then
    echo "Usage: ./find-and-go.sh SEARCH_TERM"
    exit 1
fi

echo "🔍 Searching for: $SEARCH_TERM"

# Search from current directory first
FOUND=$(find . -maxdepth 3 -type d -iname "*$SEARCH_TERM*" 2>/dev/null | head -1)

if [ -z "$FOUND" ]; then
    # Search from home if nothing found locally
    echo "  Nothing local, searching from ~..."
    FOUND=$(find ~ -maxdepth 4 -type d -iname "*$SEARCH_TERM*" 2>/dev/null | head -1)
fi

if [ -n "$FOUND" ]; then
    echo "✓ Found: $FOUND"
    
    if [ "$GOBASH_SHELL" = "1" ]; then
        # In gobash: just cd (gobash handles the rest)
        cd "$FOUND"
        echo "  Changed to: $(pwd)"
    else
        echo "  Run in gobash to automatically navigate"
    fi
else
    echo "✗ No directory matching '$SEARCH_TERM' found"
    exit 1
fi
```

Usage:
```bash
[IPTP-1] ~$ ./find-and-go.sh golang
🔍 Searching for: golang
✓ Found: ./projects/golang-api
  Changed to: /Users/you/projects/golang-api
[IPTP-1] golang-api$ 
```

## 📝 Example 4: Build & Test Runner

Create `dev-cycle.sh`:

```bash
#!/bin/bash
# dev-cycle.sh - Build, test, and run development cycle

set -e  # Exit on error

PROJECT_TYPE=${1:-go}  # Default to Go

echo "🔨 Development Cycle ($PROJECT_TYPE)"
echo "================================"

case $PROJECT_TYPE in
    go)
        echo "→ Building Go project..."
        go build -o bin/app
        echo "✓ Build successful"
        
        echo ""
        echo "→ Running tests..."
        go test ./... -v
        echo "✓ Tests passed"
        
        echo ""
        echo "→ Running application..."
        ./bin/app
        ;;
        
    node)
        echo "→ Installing dependencies..."
        npm install
        echo "✓ Dependencies installed"
        
        echo ""
        echo "→ Running tests..."
        npm test
        echo "✓ Tests passed"
        
        echo ""
        echo "→ Starting dev server..."
        npm run dev
        ;;
        
    python)
        echo "→ Setting up virtual environment..."
        python3 -m venv venv
        source venv/bin/activate
        pip install -r requirements.txt
        echo "✓ Environment ready"
        
        echo ""
        echo "→ Running tests..."
        pytest
        echo "✓ Tests passed"
        ;;
        
    *)
        echo "✗ Unknown project type: $PROJECT_TYPE"
        echo "  Supported: go, node, python"
        exit 1
        ;;
esac

if [ "$GOBASH_SHELL" = "1" ]; then
    echo ""
    echo "💾 Tip: Use 'save' to save this location in gobash"
fi
```

Usage:
```bash
[myapp] myapp$ ./dev-cycle.sh go
🔨 Development Cycle (go)
================================
→ Building Go project...
✓ Build successful
...
```

## 📝 Example 5: Interactive Workspace Selector

Create `select-workspace.sh`:

```bash
#!/bin/bash
# select-workspace.sh - Interactive workspace selection

WORKSPACES=(
    "$HOME/projects/auth-service:Authentication microservice"
    "$HOME/projects/web-frontend:React frontend"
    "$HOME/projects/data-pipeline:ETL pipeline"
    "$HOME/personal/blog:Personal blog"
)

echo "🗂️  Available Workspaces"
echo "======================="
echo ""

# Display menu
for i in "${!WORKSPACES[@]}"; do
    IFS=':' read -r path desc <<< "${WORKSPACES[$i]}"
    dirname=$(basename "$path")
    printf "%d) %-20s - %s\n" $((i+1)) "$dirname" "$desc"
done

echo ""
read -p "Select workspace (1-${#WORKSPACES[@]}): " choice

# Validate input
if ! [[ "$choice" =~ ^[0-9]+$ ]] || [ "$choice" -lt 1 ] || [ "$choice" -gt "${#WORKSPACES[@]}" ]; then
    echo "✗ Invalid selection"
    exit 1
fi

# Get selected workspace
selected="${WORKSPACES[$((choice-1))]}"
IFS=':' read -r path desc <<< "$selected"

if [ ! -d "$path" ]; then
    echo "✗ Directory doesn't exist: $path"
    exit 1
fi

echo ""
echo "✓ Selected: $desc"
echo "  Path: $path"

if [ "$GOBASH_SHELL" = "1" ]; then
    cd "$path"
    echo ""
    echo "💡 Tip: Use 'name \"$desc\"' to name this session"
else
    echo ""
    echo "⚠️  Run inside gobash to automatically navigate"
    echo "  Or manually: cd $path"
fi
```

Usage:
```bash
[IPTP-1] ~$ ./select-workspace.sh
🗂️  Available Workspaces
=======================

1) auth-service      - Authentication microservice
2) web-frontend      - React frontend
3) data-pipeline     - ETL pipeline
4) blog              - Personal blog

Select workspace (1-4): 2

✓ Selected: React frontend
  Path: /Users/you/projects/web-frontend
💡 Tip: Use 'name "React frontend"' to name this session

[IPTP-1] web-frontend$ 
```

## 📝 Example 6: Context-Aware Git Helper

Create `git-status-all.sh`:

```bash
#!/bin/bash
# git-status-all.sh - Check git status in current and subdirectories

echo "📊 Git Repository Status"
echo "======================="

check_repo() {
    local dir=$1
    local name=$(basename "$dir")
    
    if [ -d "$dir/.git" ]; then
        cd "$dir"
        
        # Get branch
        branch=$(git branch --show-current 2>/dev/null)
        
        # Check for changes
        if git diff-index --quiet HEAD -- 2>/dev/null; then
            status="✓ clean"
        else
            status="⚠️  changes"
        fi
        
        # Check for unpushed commits
        unpushed=$(git log origin/$branch..HEAD 2>/dev/null | grep -c "^commit" || echo "0")
        
        printf "  %-30s [%s] %s" "$name" "$branch" "$status"
        [ "$unpushed" -gt 0 ] && printf " (%d unpushed)" "$unpushed"
        echo ""
    fi
}

# Check current directory
if [ -d ".git" ]; then
    echo ""
    echo "Current repository:"
    check_repo "."
fi

# Check subdirectories
echo ""
echo "Subdirectories:"
for dir in */; do
    [ -d "$dir" ] && check_repo "$dir"
done

if [ "$GOBASH_SHELL" = "1" ]; then
    echo ""
    echo "💡 In gobash: Use 'cd <repo>' to navigate, 'save' to bookmark"
fi
```

Usage:
```bash
[IPTP-1] projects$ ./git-status-all.sh
📊 Git Repository Status
=======================

Current repository:
  projects                    [main] ✓ clean

Subdirectories:
  auth-service               [feature/oauth] ⚠️  changes (2 unpushed)
  web-frontend               [main] ✓ clean
  data-pipeline              [dev] ⚠️  changes
```

## 🎯 Best Practices

### 1. **Always Check for gobash Context**
```bash
if [ "$GOBASH_SHELL" = "1" ]; then
    # gobash-specific behavior
else
    # Fallback behavior
fi
```

### 2. **Provide Helpful Tips**
```bash
echo "💡 Tip: Use 'name' to name this session"
echo "💾 Tip: Use 'save' to bookmark this location"
```

### 3. **Make Scripts Portable**
```bash
#!/bin/bash
# Works both inside and outside gobash
set -e  # Exit on error
```

### 4. **Use Meaningful Output**
```bash
echo "✓ Success message"
echo "✗ Error message"
echo "⚠️  Warning message"
echo "💡 Tip/suggestion"
echo "🔍 Searching..."
```

### 5. **Document Your Scripts**
```bash
#!/bin/bash
# script-name.sh - Brief description
#
# Usage: ./script-name.sh [ARGS]
# Example: ./script-name.sh myproject
```

## 🚀 Advanced: Scripts that Call gobash Commands

You can't directly call gobash commands from external scripts (since gobash commands run in the gobash process). But you can:

### Method 1: Use Environment Variables
```bash
#!/bin/bash
# The script runs in gobash context, so cd works!

cd ~/projects/myapp
# gobash will save the state automatically
```

### Method 2: Print Instructions
```bash
#!/bin/bash
# generate-commands.sh

echo "# Run these commands in gobash:"
echo "cd ~/projects/myapp"
echo "name \"working on myapp\""
echo "save"
```

Usage:
```bash
[IPTP-1] ~$ ./generate-commands.sh
# Run these commands in gobash:
cd ~/projects/myapp
name "working on myapp"
save

[IPTP-1] ~$ cd ~/projects/myapp
[IPTP-1] myapp$ name "working on myapp"
✓ Shell named: myapp
[myapp] myapp$ save
✓ Saved state for: myapp @ /Users/you/projects/myapp
```

## 📚 Script Library Ideas

Here are some useful scripts you might want to create:

1. **project-switcher.sh** - Quick project switching
2. **backup-state.sh** - Backup gobash state
3. **workspace-setup.sh** - Set up development environment
4. **git-batch.sh** - Batch git operations
5. **docker-dev.sh** - Docker development helpers
6. **test-runner.sh** - Automated testing workflows
7. **deploy.sh** - Deployment automation
8. **cleanup.sh** - Clean up build artifacts

## 🎓 Learning Path

1. **Start Simple**: Begin with status/info scripts
2. **Add Interactivity**: Use read to get user input
3. **Check Context**: Use GOBASH_SHELL to adapt behavior
4. **Combine Tools**: Use gobash + git + docker + your tools
5. **Share**: Create a library of useful scripts

---

Remember: gobash scripts are just regular bash scripts that can leverage gobash's context awareness! 🚀
