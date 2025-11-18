#!/bin/bash
# workspace.sh - Quick workspace selector for gobash

# Define your workspaces here (path:description)
WORKSPACES=(
    "$HOME/projects:Main projects"
    "$HOME/personal:Personal projects"
    "$HOME/work:Work projects"
    "$HOME/Documents:Documents"
    "$HOME/Downloads:Downloads"
)

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🗂️  Quick Workspace Selector"
echo "============================"
echo ""

# Show current location
if [ "$GOBASH_SHELL" = "1" ]; then
    echo "Current: ${BLUE}$(pwd)${NC}"
    echo ""
fi

# Display menu
for i in "${!WORKSPACES[@]}"; do
    IFS=':' read -r path desc <<< "${WORKSPACES[$i]}"
    if [ -d "$path" ]; then
        dirname=$(basename "$path")
        printf "%d) ${GREEN}%-15s${NC} - %s\n" $((i+1)) "$dirname" "$desc"
    fi
done

echo ""
echo "Or enter a custom path..."
echo ""
read -p "Select (1-${#WORKSPACES[@]}, or path): " choice

# Handle custom path
if [[ "$choice" =~ ^[~/] ]]; then
    # Expand ~ if present
    path="${choice/#\~/$HOME}"
    
    if [ ! -d "$path" ]; then
        echo "${RED}✗${NC} Directory doesn't exist: $path"
        exit 1
    fi
    
    echo ""
    echo "${GREEN}✓${NC} Going to: $path"
    
    if [ "$GOBASH_SHELL" = "1" ]; then
        cd "$path" || exit 1
        echo "  Current: $(pwd)"
    else
        echo "  Run inside gobash to navigate automatically"
    fi
    exit 0
fi

# Validate numeric input
if ! [[ "$choice" =~ ^[0-9]+$ ]] || [ "$choice" -lt 1 ] || [ "$choice" -gt "${#WORKSPACES[@]}" ]; then
    echo "${RED}✗${NC} Invalid selection"
    exit 1
fi

# Get selected workspace
selected="${WORKSPACES[$((choice-1))]}"
IFS=':' read -r path desc <<< "$selected"

if [ ! -d "$path" ]; then
    echo "${RED}✗${NC} Directory doesn't exist: $path"
    exit 1
fi

echo ""
echo "${GREEN}✓${NC} Selected: $desc"

if [ "$GOBASH_SHELL" = "1" ]; then
    cd "$path" || exit 1
    echo "  Current: $(pwd)"
    echo ""
    echo "💡 Tip: Use ${BLUE}name \"$desc\"${NC} to name this session"
else
    echo "  Path: $path"
    echo ""
    echo "⚠️  Run inside gobash to navigate automatically"
fi
