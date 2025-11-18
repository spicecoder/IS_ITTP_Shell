# 🚀 iptp - Quick Start for Your Mac

## What This Is

**A REAL SHELL BINARY** written in Go - not a bash wrapper!
- ✅ No sourcing needed
- ✅ Actually changes directories  
- ✅ Runs scripts in iptp context
- ✅ Cross-platform (Mac, Linux, Windows)

## Installation (30 seconds)

### For Your Mac (Apple Silicon M1/M2/M3)

```bash
# 1. Download the binary (already in your Downloads)
cd ~/Downloads/iptp-go/dist

# 2. Make it executable
chmod +x iptp-darwin-arm64

# 3. Run it!
./iptp-darwin-arm64
```

### Or Install System-Wide

```bash
# Copy to your PATH
sudo cp iptp-darwin-arm64 /usr/local/bin/iptp

# Now run from anywhere
iptp
```

## Your First Session

```bash
$ ./iptp-darwin-arm64
🚀 iptp - IPTP Shell Process Manager
   Type 'help' for commands, 'exit' to quit

[shell_12345] ~$ name "working on my project"
✓ Shell named: my_project
  Intention: working on my project

[my_project] ~$ goto ~/Documents
✓ Changed to: /Users/pronab/Documents

[my_project] ~/Documents$ save
✓ Saved state for: my_project @ /Users/pronab/Documents

[my_project] ~/Documents$ exit
Goodbye!
```

## Why This Solves Your Problem

### The OLD Problem:
```bash
./getmethere.sh
# Select directory...
# BUT YOU'RE STILL IN THE SAME PLACE! 😡
```

### The NEW Solution:
```bash
./iptp-darwin-arm64
[shell]$ getmethere
# Select directory...
# YOU'RE ACTUALLY THERE! 🎉
```

Download the entire `iptp-go` folder and you have everything!
