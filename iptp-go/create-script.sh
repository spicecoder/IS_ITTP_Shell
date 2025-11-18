#!/bin/bash
# create-script.sh - Generate a new gobash script template

SCRIPT_NAME=$1

if [ -z "$SCRIPT_NAME" ]; then
    echo "Usage: ./create-script.sh SCRIPT_NAME"
    echo ""
    echo "Example: ./create-script.sh myscript"
    echo "  Creates: myscript.sh"
    exit 1
fi

# Add .sh if not present
[[ "$SCRIPT_NAME" != *.sh ]] && SCRIPT_NAME="${SCRIPT_NAME}.sh"

if [ -f "$SCRIPT_NAME" ]; then
    echo "✗ File already exists: $SCRIPT_NAME"
    exit 1
fi

cat > "$SCRIPT_NAME" << 'TEMPLATE'
#!/bin/bash
# SCRIPT_NAME - Brief description of what this script does

# Exit on error
set -e

# Colors for nice output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running in gobash
in_gobash() {
    [ "$GOBASH_SHELL" = "1" ]
}

# Main function
main() {
    echo "🚀 SCRIPT_NAME"
    echo "=================="
    echo ""
    
    # Your script logic here
    
    if in_gobash; then
        echo ""
        echo "${GREEN}✓${NC} Running in gobash context"
        echo "  Process: ${GOBASH_PROCESS:-unnamed}"
        echo "  Directory: ${GOBASH_CWD}"
    else
        echo ""
        echo "${YELLOW}ℹ${NC} Tip: Run inside gobash for enhanced features"
    fi
    
    echo ""
    echo "${GREEN}✓${NC} Done!"
}

# Run main function
main "$@"
TEMPLATE

# Replace SCRIPT_NAME in template
sed -i.bak "s/SCRIPT_NAME/$SCRIPT_NAME/g" "$SCRIPT_NAME" && rm "${SCRIPT_NAME}.bak"

# Make executable
chmod +x "$SCRIPT_NAME"

echo "✓ Created script: $SCRIPT_NAME"
echo ""
echo "Template includes:"
echo "  • Proper shebang and error handling"
echo "  • Color output support"
echo "  • gobash context detection"
echo "  • Basic structure"
echo ""
echo "Next steps:"
echo "  1. Edit $SCRIPT_NAME"
echo "  2. Add your logic to the main() function"
echo "  3. Test: ./$SCRIPT_NAME"
echo ""
echo "Try it now:"
echo "  ./gobash"
echo "  ./$SCRIPT_NAME"
