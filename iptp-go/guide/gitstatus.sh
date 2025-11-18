#!/bin/bash
# gitstatus.sh - Check git status across multiple repositories

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "📊 Git Repository Status"
echo "========================"
echo ""

check_repo() {
    local dir=$1
    local name=$(basename "$dir")
    
    if [ ! -d "$dir/.git" ]; then
        return
    fi
    
    cd "$dir" || return
    
    # Get branch
    branch=$(git branch --show-current 2>/dev/null)
    [ -z "$branch" ] && branch="(detached)"
    
    # Check for changes
    local status_symbol
    local status_text
    if git diff-index --quiet HEAD -- 2>/dev/null; then
        status_symbol="${GREEN}✓${NC}"
        status_text="clean"
    else
        status_symbol="${YELLOW}⚠${NC}"
        status_text="modified"
        
        # Count changes
        modified=$(git diff --name-only | wc -l | tr -d ' ')
        [ "$modified" -gt 0 ] && status_text="$status_text ($modified files)"
    fi
    
    # Check for untracked files
    untracked=$(git ls-files --others --exclude-standard | wc -l | tr -d ' ')
    [ "$untracked" -gt 0 ] && status_text="$status_text, ${YELLOW}$untracked untracked${NC}"
    
    # Check for unpushed commits
    if git rev-parse --abbrev-ref @{u} >/dev/null 2>&1; then
        unpushed=$(git log @{u}..HEAD 2>/dev/null | grep -c "^commit" || echo "0")
        [ "$unpushed" -gt 0 ] && status_text="$status_text, ${BLUE}$unpushed unpushed${NC}"
    fi
    
    # Print status
    printf "  %s %-30s ${BLUE}%-15s${NC} %s\n" "$status_symbol" "$name" "[$branch]" "$status_text"
}

# Check current directory
if [ -d ".git" ]; then
    echo "Current directory:"
    check_repo "$(pwd)"
    echo ""
fi

# Check subdirectories
found_any=false
echo "Subdirectories:"
for dir in */; do
    if [ -d "$dir/.git" ]; then
        check_repo "$dir"
        found_any=true
    fi
done

if [ "$found_any" = false ]; then
    echo "  No git repositories found in subdirectories"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Summary
total_repos=0
dirty_repos=0

for dir in */ .; do
    [ ! -d "$dir/.git" ] && continue
    total_repos=$((total_repos + 1))
    
    cd "$dir" || continue
    if ! git diff-index --quiet HEAD -- 2>/dev/null; then
        dirty_repos=$((dirty_repos + 1))
    fi
    cd - >/dev/null || exit
done

echo "Summary: $total_repos repositories, $dirty_repos with changes"

if [ "$GOBASH_SHELL" = "1" ]; then
    echo ""
    echo "💡 Tip: Use 'cd <repo>' to navigate to a repository"
fi
