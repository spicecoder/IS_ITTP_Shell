#!/bin/bash
# qf.sh (quick find) - Fast directory finder for gobash

SEARCH_TERM=$1

if [ -z "$SEARCH_TERM" ]; then
    echo "Usage: ./qf.sh SEARCH_TERM"
    echo ""
    echo "Examples:"
    echo "  ./qf.sh golang    # Find directories containing 'golang'"
    echo "  ./qf.sh myapp     # Find directories containing 'myapp'"
    exit 1
fi

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "🔍 Searching for: ${BLUE}$SEARCH_TERM${NC}"
echo ""

# Search from current directory first (fast!)
echo "Searching locally..."
matches=()
while IFS= read -r -d '' dir; do
    matches+=("$dir")
done < <(find . -maxdepth 3 -type d -iname "*$SEARCH_TERM*" -print0 2>/dev/null)

# If nothing found locally, search from home (slower)
if [ ${#matches[@]} -eq 0 ]; then
    echo "Nothing local, searching from home directory..."
    while IFS= read -r -d '' dir; do
        matches+=("$dir")
    done < <(find ~ -maxdepth 5 -type d -iname "*$SEARCH_TERM*" -print0 2>/dev/null | head -z -20)
fi

if [ ${#matches[@]} -eq 0 ]; then
    echo "${YELLOW}✗${NC} No directories found matching '$SEARCH_TERM'"
    exit 1
fi

echo ""
echo "Found ${GREEN}${#matches[@]}${NC} matching directories:"
echo ""

# Display results
for i in "${!matches[@]}"; do
    dir="${matches[$i]}"
    # Make path relative if possible
    rel_path=$(realpath --relative-to="." "$dir" 2>/dev/null || echo "$dir")
    printf "%2d) %s\n" $((i+1)) "$rel_path"
done

# Auto-select if only one match
if [ ${#matches[@]} -eq 1 ]; then
    if [ "$GOBASH_SHELL" = "1" ]; then
        echo ""
        echo "${GREEN}✓${NC} Only one match, navigating..."
        cd "${matches[0]}" || exit 1
        echo "  Current: $(pwd)"
        exit 0
    fi
fi

# Interactive selection
echo ""
read -p "Select (1-${#matches[@]}, or 'q' to quit): " choice

if [ "$choice" = "q" ] || [ "$choice" = "Q" ]; then
    echo "Cancelled"
    exit 0
fi

# Validate input
if ! [[ "$choice" =~ ^[0-9]+$ ]] || [ "$choice" -lt 1 ] || [ "$choice" -gt "${#matches[@]}" ]; then
    echo "${YELLOW}✗${NC} Invalid selection"
    exit 1
fi

selected="${matches[$((choice-1))]}"

echo ""
echo "${GREEN}✓${NC} Selected: $selected"

if [ "$GOBASH_SHELL" = "1" ]; then
    cd "$selected" || exit 1
    echo "  Current: $(pwd)"
    echo ""
    echo "💡 Tip: Use ${BLUE}save${NC} to bookmark this location"
else
    echo ""
    echo "⚠️  Run inside gobash to navigate automatically"
    echo "  Or manually: cd $selected"
fi
