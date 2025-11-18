# 🎯 Quick Guide: Creating Your Own gobash Scripts

## TL;DR

gobash scripts are **just bash scripts** that can access gobash context through environment variables!

## 🚀 The Simplest Way

### Method 1: Use the Template Generator

```bash
# Inside gobash or regular shell:
./create-script.sh myscript

# Edit it
vim myscript.sh

# Run it in gobash
./gobash
[IPTP-1] ~$ ./myscript.sh
```

### Method 2: Start from Scratch

```bash
#!/bin/bash
# myscript.sh - My custom script

if [ "$GOBASH_SHELL" = "1" ]; then
    echo "Running in gobash!"
    echo "Current dir: $GOBASH_CWD"
    echo "Process: $GOBASH_PROCESS"
    
    # Your logic here
    cd ~/projects/myapp
    # gobash automatically saves state!
else
    echo "Run inside gobash for full features"
fi
```

Make it executable:
```bash
chmod +x myscript.sh
```

## 🎓 The Magic: Environment Variables

When your script runs inside gobash, these are automatically available:

```bash
$GOBASH_SHELL    # = "1" when in gobash
$GOBASH_CWD      # Current directory
$GOBASH_PROCESS  # Current process name
```

## 💡 Common Patterns

### Pattern: Navigation Script
```bash
#!/bin/bash
cd ~/my/favorite/place
echo "Now in: $(pwd)"
# That's it! gobash handles the rest
```

### Pattern: Interactive Selector
```bash
#!/bin/bash
echo "Where to go?"
echo "1) Projects"
echo "2) Documents"
read -p "Select: " choice

case $choice in
    1) cd ~/projects ;;
    2) cd ~/Documents ;;
esac
```

### Pattern: Conditional Behavior
```bash
#!/bin/bash
if [ "$GOBASH_SHELL" = "1" ]; then
    # Inside gobash - can navigate
    cd ~/myproject
else
    # Outside gobash - just inform
    echo "Run in gobash to auto-navigate"
    echo "Or manually: cd ~/myproject"
fi
```

## 📚 Ready-Made Examples

We've created 5 example scripts for you:

1. **workspace.sh** - Quick workspace switcher
2. **newproject.sh** - Project initializer (Go, Node, Python, Rust)
3. **gitstatus.sh** - Multi-repo git status
4. **qf.sh** - Quick directory finder
5. **create-script.sh** - Template generator

Download them all from the `examples/` folder!

## 🎨 Make It Pretty

Add colors and icons:

```bash
#!/bin/bash
GREEN='\033[0;32m'
NC='\033[0m'  # No Color

echo "${GREEN}✓${NC} Success!"
echo "🚀 Starting..."
echo "💡 Tip: Use gobash commands"
```

## 🔧 Advanced: Call gobash Commands

You **can't directly call** gobash commands like `name` or `save` from scripts (they're internal to gobash). But you can:

### Work With The Environment
```bash
#!/bin/bash
# Navigate (this works!)
cd ~/projects/myapp

# The directory change happens in gobash
# No need to call 'save' - it's automatic!
```

### Print Commands
```bash
#!/bin/bash
# Generate commands for the user
echo "# Run these in gobash:"
echo "name 'my project'"
echo "save"
```

## ✅ Complete Example

Here's a real, working script:

```bash
#!/bin/bash
# goto-project.sh - Navigate to a project

PROJECT=$1

if [ -z "$PROJECT" ]; then
    echo "Usage: ./goto-project.sh PROJECT_NAME"
    exit 1
fi

PROJECT_DIR="$HOME/projects/$PROJECT"

if [ ! -d "$PROJECT_DIR" ]; then
    echo "✗ Project not found: $PROJECT"
    exit 1
fi

echo "🚀 Going to: $PROJECT"

if [ "$GOBASH_SHELL" = "1" ]; then
    cd "$PROJECT_DIR"
    echo "✓ Current: $(pwd)"
    echo ""
    echo "💡 Tip: Use 'name \"$PROJECT\"' to name this session"
else
    echo "📍 Location: $PROJECT_DIR"
    echo "⚠️  Run inside gobash to auto-navigate"
fi
```

Usage:
```bash
[IPTP-1] ~$ ./goto-project.sh myapp
🚀 Going to: myapp
✓ Current: /Users/you/projects/myapp
💡 Tip: Use 'name "myapp"' to name this session

[IPTP-1] myapp$ name "myapp"
✓ Shell named: myapp
[myapp] myapp$ 
```

## 📖 Learn More

- **[SCRIPTING_GUIDE.md](SCRIPTING_GUIDE.md)** - Full documentation
- **[examples/README.md](examples/README.md)** - Example scripts explained

## 🎯 Key Takeaway

**gobash scripts = bash scripts + gobash awareness**

That's it! No special API, no complex integration. Just check `$GOBASH_SHELL` and use normal bash commands. gobash takes care of the rest! 🚀
