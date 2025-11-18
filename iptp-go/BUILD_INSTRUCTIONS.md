# iptp - Build Instructions

## What You Have

All Go source files for building iptp from scratch:

- `main.go` - Entry point
- `state.go` - IPTP state management  
- `shell.go` - Interactive REPL
- `commands.go` - Command-line handlers
- `utils.go` - Utilities & intention parsing
- `exec.go` - Script execution
- `go.mod` - Go module file
- `build.sh` - Build script

## Quick Build (Your Mac)

### Step 1: Create a folder and download files

```bash
# Create a folder
mkdir -p ~/iptp-source
cd ~/iptp-source

# Download all .go files from Claude outputs:
# - main.go
# - state.go
# - shell.go
# - commands.go
# - utils.go
# - exec.go
# - go.mod
# - build.sh
```

### Step 2: Build for your Mac

```bash
# Make sure you're in the folder with all .go files
cd ~/iptp-source

# Build for Apple Silicon (M1/M2/M3)
GOOS=darwin GOARCH=arm64 go build -o iptp

# Or if Intel Mac:
GOOS=darwin GOARCH=amd64 go build -o iptp
```

### Step 3: Run it!

```bash
./iptp
```

## Build for All Platforms

```bash
# Make build script executable
chmod +x build.sh

# Build for all platforms
./build.sh
```

This creates:
- `dist/iptp-darwin-arm64` - macOS Apple Silicon
- `dist/iptp-darwin-amd64` - macOS Intel
- `dist/iptp-linux-amd64` - Linux
- `dist/iptp-linux-arm64` - Linux ARM
- `dist/iptp-windows-amd64.exe` - Windows

## Requirements

- **Go 1.22 or later**
- No external dependencies!

### Install Go on Mac

```bash
# Using Homebrew
brew install go

# Verify
go version
```

## File Structure

```
iptp-source/
├── main.go           # Entry point and initialization
├── state.go          # State management (IPTP)
├── shell.go          # Interactive shell REPL
├── commands.go       # Command-line mode
├── utils.go          # Utilities and parsing
├── exec.go           # Script execution
├── go.mod            # Go module definition
└── build.sh          # Build script
```

## Testing the Build

```bash
# Test version
./iptp version

# Test help
./iptp help

# Test interactive mode
./iptp
[shell]$ help
[shell]$ name "testing"
[shell]$ goto ~/Documents
[shell]$ exit
```

## Features Included

✅ **Core Navigation**
- `goto PATH` - Navigate
- `goto '*pattern*'` - Fuzzy match
- `getmethere` - Interactive finder
- `back` - Go back in history

✅ **Process Management**
- `name "intention"` - Natural language
- `save` - Save state
- `list` - List processes
- `jump PROCESS` - Jump to location
- `state` - Show IPTP state

✅ **Script Execution**
- `./script.sh` - Run scripts
- Any command works

✅ **IPTP Compliance**
- Trivalent pulses (Y/N/U)
- External state
- Plain language intentions
- No direct process communication

✅ **Cross-Platform**
- macOS (Intel + Apple Silicon)
- Linux (amd64 + arm64)
- Windows

## Binary Sizes

- macOS: ~3.3-3.4 MB
- Linux: ~3.4 MB  
- Windows: ~3.5 MB

All are single, standalone binaries with zero dependencies!

## Troubleshooting

### "go: command not found"

Install Go:
```bash
# macOS
brew install go

# Linux (Ubuntu/Debian)
sudo apt-get install golang-go

# Or download from: https://go.dev/dl/
```

### "cannot find package"

Make sure all .go files are in the same directory:
```bash
ls -la
# Should show: main.go, state.go, shell.go, commands.go, utils.go, exec.go, go.mod
```

### Build fails

```bash
# Clean and rebuild
rm -rf dist/
go clean
go build -o iptp
```

## What Makes This Special

Unlike the bash version:
- ✅ **No sourcing** - Just run the binary
- ✅ **True shell** - iptp IS your shell
- ✅ **getmethere works** - Actually changes directory
- ✅ **Scripts work** - Full execution context
- ✅ **Cross-platform** - Windows support
- ✅ **Single file** - One binary, zero dependencies

## Example Session

```bash
$ ./iptp
🚀 iptp - IPTP Shell Process Manager

[shell_12345] ~$ name "debugging auth"
✓ Shell named: debugging_auth
  Intention: debugging auth

[debugging_auth] ~$ getmethere
Hello, pronab! Where would you like to work today?
Enter directory name: Documents

Found 1 matching directories:
  1) /Users/pronab/Documents

Enter number: 1
✓ Changed to: /Users/pronab/Documents

[debugging_auth] ~/Documents$ save
✓ Saved state for: debugging_auth @ /Users/pronab/Documents

[debugging_auth] ~/Documents$ list
=== Available Processes ===
  → debugging_auth: /Users/pronab/Documents (PID: 12345)

[debugging_auth] ~/Documents$ exit
Goodbye!
```

## State File

All state saved to:
- **macOS/Linux**: `/tmp/iptp_state.json`
- **Windows**: `%TEMP%\iptp_state.json`

View it:
```bash
cat /tmp/iptp_state.json
```

## Next Steps

1. **Build it**: `go build -o iptp`
2. **Run it**: `./iptp`
3. **Test getmethere**: It works! No more subshell trap!
4. **Use it**: Real IPTP shell, not a bash function

Enjoy your freedom from bash limitations! 🎉
