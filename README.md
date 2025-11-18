# iptp - IPTP Shell Process Manager

**A real shell binary written in Go** - not a bash wrapper!

## What Makes This Different

Unlike the bash function version, this is a **standalone binary** that runs as your actual shell:

```bash
# NOT this (bash function requiring sourcing):
source iptp_init.sh
iptp goto /path

# THIS (real binary - just run it!):
./iptp
[shell]$ goto /path
[shell]$ getmethere
[shell]$ ./myscript.sh    # Scripts run in iptp context!
```

## Features

✅ **Real Shell** - No bash/zsh dependency  
✅ **Cross-Platform** - Works on macOS, Linux, Windows  
✅ **Single Binary** - No installation, just download and run  
✅ **IPTP-Compliant** - Intentions, trivalent pulses, external state  
✅ **Script Execution** - Run ANY script in iptp context  
✅ **Natural Language** - Intention parsing built-in  
✅ **Interactive & Command Mode** - Use as shell OR run single commands  

## Installation

### macOS (Your Machine!)

**Option 1: Apple Silicon (M1/M2/M3)**
```bash
# Download the binary
curl -O https://[your-server]/iptp-darwin-arm64

# Make executable
chmod +x iptp-darwin-arm64

# Move to PATH
sudo mv iptp-darwin-arm64 /usr/local/bin/iptp

# Run it!
iptp
```

**Option 2: Intel Mac**
```bash
curl -O https://[your-server]/iptp-darwin-amd64
chmod +x iptp-darwin-amd64
sudo mv iptp-darwin-amd64 /usr/local/bin/iptp
iptp
```

### Linux

```bash
# Download
curl -O https://[your-server]/iptp-linux-amd64

# Make executable
chmod +x iptp-linux-amd64

# Move to PATH
sudo mv iptp-linux-amd64 /usr/local/bin/iptp

# Run
iptp
```

### Windows

1. Download `iptp-windows-amd64.exe`
2. Rename to `iptp.exe`
3. Add to PATH or run directly
4. Open Command Prompt or PowerShell
5. Run: `iptp.exe`

## Usage

### Interactive Mode (Shell)

```bash
$ iptp
🚀 iptp - IPTP Shell Process Manager
   Type 'help' for commands, 'exit' to quit

[shell_1234] ~$ name "working on authentication"
✓ Shell named: authentication
  Intention: working on authentication
  (parsed process name: authentication)

[authentication] ~$ goto /var/www
✓ Changed to: /var/www

[authentication] /var/www$ save
✓ Saved state for: authentication @ /var/www

[authentication] /var/www$ list
=== Available Processes ===
  → authentication: /var/www (PID: 1234)

[authentication] /var/www$ exit
Goodbye!
```

### Command Mode (Single Commands)

```bash
# Name a process
$ iptp name "debugging nginx"
✓ Shell named: debugging_nginx
  Intention: debugging nginx

# Navigate
$ iptp goto /var/log/nginx
✓ Changed to: /var/log/nginx

# List processes
$ iptp list
=== Available Processes ===
  → debugging_nginx: /var/log/nginx (PID: 5678)

# Show state
$ iptp state
=== Current Process State (IPTP Format) ===
Process: debugging_nginx
Directory: /var/log/nginx
...
```

## Commands

### Built-in Commands

| Command | Description | Example |
|---------|-------------|---------|
| `name [INTENTION]` | Name current process | `name "working on auth"` |
| `goto PATH` | Navigate to directory | `goto /var/www` |
| `goto *pattern*` | Fuzzy directory match | `goto '*nginx*'` |
| `getmethere` | Interactive dir finder | `getmethere` |
| `save` | Save current state | `save` |
| `list` | List all processes | `list` |
| `jump PROCESS` | Jump to saved process | `jump webdev` |
| `back` | Go back in history | `back` |
| `state` | Show current state | `state` |
| `help` | Show help | `help` |
| `exit` | Exit iptp | `exit` |

### Script Execution

```bash
[shell]$ ./myscript.sh
[shell]$ python3 analyze.py
[shell]$ ls -la
[shell]$ any-command --args
```

**Any command or script runs in iptp context!**

## The getmethere Feature

Interactive directory finder - no more "cd doesn't work" problems!

```bash
[shell]$ getmethere
Hello, pronab! Where would you like to work today?
Enter the directory name: nginx

Found 3 matching directories:
  1) /var/log/nginx
  2) /etc/nginx
  3) /home/user/projects/nginx-config

Enter number (or 'q' to quit): 2
✓ Changed to: /etc/nginx
```

## IPTP Features

### Natural Language Intentions

```bash
name "I am working on authentication"
# → Process: authentication
# → Intention: "I am working on authentication"

name "debugging nginx configuration"
# → Process: nginx_configuration
# → Intention: "debugging nginx configuration"
```

### Trivalent Pulses

State is tracked with Y/N/U (Yes/No/Undecided):

```json
{
  "processes": {
    "authentication": {
      "intention": "I am working on authentication",
      "current_dir": "/var/www/auth",
      "pulses": [
        {"name": "process named", "TV": "Y", "response": "authentication"},
        {"name": "directory saved", "TV": "Y", "response": "/var/www/auth"}
      ]
    }
  }
}
```

### External State

All state stored in:
- **Linux/macOS**: `/tmp/iptp_state.json`
- **Windows**: `%TEMP%\iptp_state.json`

State persists across shell sessions (until reboot on Unix).

## Real-World Workflows

### Multi-Context Development

**Terminal 1:**
```bash
$ iptp
[shell]$ name frontend
[frontend]$ goto ~/projects/app/frontend
[frontend]$ save
```

**Terminal 2:**
```bash
$ iptp
[shell]$ name backend
[backend]$ goto ~/projects/app/backend
[backend]$ save
```

**Terminal 3 - Jump between contexts:**
```bash
$ iptp
[shell]$ list
  → frontend: /Users/pronab/projects/app/frontend
  → backend: /Users/pronab/projects/app/backend

[shell]$ jump frontend
✓ Jumped to frontend @ /Users/pronab/projects/app/frontend
[frontend]$ 
```

### Script Execution in Context

```bash
[webdev]$ goto ~/projects/website
[webdev]$ ./build.sh          # Runs in iptp context
[webdev]$ ./deploy.sh prod    # Has access to iptp env
[webdev]$ save                # Save after deployment
```

## Why Go Binary vs Bash Function

| Feature | Bash Function | Go Binary |
|---------|---------------|-----------|
| **Requires sourcing** | ✗ Yes | ✅ No |
| **True shell process** | ✗ No | ✅ Yes |
| **Cross-platform** | ⚠️ Unix only | ✅ Win/Mac/Linux |
| **Single binary** | ✗ No | ✅ Yes |
| **Script execution** | ⚠️ Limited | ✅ Full |
| **State management** | ✅ Yes | ✅ Yes |
| **IPTP compliant** | ✅ Yes | ✅ Yes |

## Technical Details

### Binary Sizes

- **macOS (Apple Silicon)**: ~3.3 MB
- **macOS (Intel)**: ~3.4 MB
- **Linux (amd64)**: ~3.4 MB
- **Windows**: ~3.4 MB

### Dependencies

**Zero runtime dependencies!** Everything is statically compiled into the binary.

### State File Format

```json
{
  "processes": {
    "process_name": {
      "intention": "Human readable intention",
      "current_dir": "/absolute/path",
      "history": ["/previous/paths"],
      "pid": 12345,
      "timestamp": "2025-11-13T10:30:00Z",
      "pulses": [
        {"name": "pulse_name", "TV": "Y", "response": "value"}
      ]
    }
  }
}
```

## Building from Source

```bash
# Clone/download source
cd iptp-go

# Build for your platform
go build -o iptp

# Or build for all platforms
./build.sh
```

### Build Requirements

- Go 1.22 or later
- No external dependencies

## Troubleshooting

### "iptp: command not found"

**Solution**: Add to PATH or use full path

```bash
# Add to PATH (macOS/Linux)
sudo mv iptp /usr/local/bin/

# Or use full path
./iptp
```

### State file permissions

If you get permission errors:

```bash
# Linux/macOS
chmod 644 /tmp/iptp_state.json

# Windows
# Run as administrator if needed
```

### getmethere is slow

The search depth is limited to 5 levels from home. If still slow:

```bash
# Start search from a subdirectory
[shell]$ goto ~/projects
[shell]$ getmethere
# Now searching from ~/projects instead of ~
```

## Comparison with Bash Version

Both versions support the same IPTP concepts, but the Go binary offers:

✅ **No sourcing required** - Just run the binary  
✅ **True shell process** - iptp IS your shell  
✅ **Cross-platform** - Windows support  
✅ **Single file deployment** - No dependencies  
✅ **Better script execution** - Scripts run in iptp context  
✅ **Faster** - Compiled Go vs interpreted bash  

The bash version is great for quick prototyping and understanding the concepts. The Go version is production-ready.

## Future Enhancements

- [ ] Tab completion
- [ ] Command history (up/down arrows)
- [ ] Persistent state (survive reboots)
- [ ] Process groups
- [ ] Remote state sync
- [ ] LLM integration for natural language commands
- [ ] Syntax highlighting
- [ ] Customizable prompts

## Philosophy

> **iptp embodies IPTP: Processes communicate through intentions and state, not direct calls.**

This creates:
- **Loose coupling**: Processes are independent
- **Semantic clarity**: Every action has clear meaning
- **Persistent coordination**: State survives process lifetime
- **Human-readable**: Plain language throughout

## License

Open source - modify and extend as needed.

## Author

Based on Intention Pulse Transfer Protocol (IPTP) by Pronab Pal  
Implementation: IntentixLab Keybyte Systems, Melbourne, Australia  
Go version: v1.0.0
# IS_ITTP_Shell
