# 🚀 IPTP - Intention Pulse Transfer Protocol Shell

**IntentixLab Keybyte Systems** - Melbourne, Australia

A next-generation shell and process coordination system based on Intention Space theory, implementing IPTP (Intention Pulse Transfer Protocol) alongside CPUX (Coordinated Process Under eXecution) concepts.

## What Makes IPTP Different

Traditional shells have a fundamental limitation: **processes can't communicate state changes to their parent shells**. When you run a script that changes directories, you're still in the same place when it exits because the child process can't modify the parent's state.

IPTP solves this through **external state management** - processes coordinate through shared knowledge (The Field) rather than direct calls.

## Features

### 🧭 Navigation that Actually Works
```bash
[IPTP-1] ~$ getmethere
Hello, pronab! Where would you like to work today?
Enter the directory name: projects
Searching...

Found 3 matching directories:
  1) ./work/projects
  2) ./personal/projects
  3) ./archive/old-projects

Enter number (or 'q' to quit): 1
✓ Changed to: /home/pronab/work/projects
```

### 💭 Intention-Based Process Naming
```bash
[IPTP-1] ~$ name "working on authentication module"
✓ Shell named: authentication_module
  Intention: working on authentication module

[authentication_module] ~$ goto ~/code/auth
[authentication_module] auth$ save
✓ Saved state for: authentication_module @ /home/pronab/code/auth
```

### 📡 Smart WiFi Hotspot (NEW!)
```bash
# One command to share internet + monitor devices
[IPTP-1] ~$ hotspot auto
🔌 Checking WiFi connection status...
✗ Not connected to WiFi
📱 Auto-enabling hotspot...
✓ Hotspot auto-enabled!
  SSID: IPTP-Hotspot
  Your IP: 192.168.2.1

🌐 Auto-starting DNS router...
✓ DNS router started on 192.168.2.1

✨ Your network is ready!
   Devices can now connect and their DNS queries will be logged
```

### 🌐 Built-in DNS Router
```bash
# Start DNS router with logging
sudo ./iptp
[IPTP-1] ~$ dns start
✓ DNS Router starting on 0.0.0.0:53
  Upstream: 8.8.8.8:53
  Log file: /tmp/iptp_dns_queries.log

# Monitor DNS queries from connected devices
[IPTP-1] ~$ dns logs 20
=== Last 20 DNS Queries ===
[14:23:45] 192.168.2.10 -> google.com. (A) = 142.250.185.46
[14:23:46] 192.168.2.11 -> facebook.com. (A) = 157.240.22.35
...

# View statistics
[IPTP-1] ~$ dns stats
=== DNS Router Statistics ===
Total queries: 1247
Unique domains: 342
Running: true
```

See [DNS_ROUTER.md](DNS_ROUTER.md) for complete DNS router documentation.
See [HOTSPOT.md](HOTSPOT.md) for complete hotspot documentation.

### 🔄 Process State Persistence
```bash
# Work in multiple sessions
[IPTP-1] ~$ name "debugging login issue"
[debugging_login_issue] ~$ goto ~/app/auth
[debugging_login_issue] auth$ save

# Later, in a new shell
[IPTP-2] ~$ list
=== Available Processes ===
  → debugging_login_issue: /home/pronab/app/auth (PID: 12345)
  → authentication_module: /home/pronab/code/auth (PID: 12340)

[IPTP-2] ~$ jump debugging_login_issue
✓ Jumped to debugging_login_issue @ /home/pronab/app/auth
```

## Installation

### Quick Start (macOS Apple Silicon)

```bash
# Download the repository
cd ~/Downloads/iptp

# Make executable
chmod +x build.sh

# Build
./build.sh

# Run it
sudo ./iptp  # sudo needed for DNS + hotspot

# Or install system-wide
sudo cp dist/iptp-darwin-arm64 /usr/local/bin/iptp
sudo iptp
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/intentixlab/iptp.git
cd iptp

# Build for all platforms
./build.sh

# Or build for your platform only
go build -o iptp .

# Run
sudo ./iptp
```

### System Requirements

- **Go 1.21+** (for building from source)
- **sudo/admin access** (for DNS router port 53 and hotspot)
- **Supported Platforms**: macOS (ARM64/AMD64), Linux (AMD64/ARM64), Windows (AMD64)

## Quick Command Reference

### Navigation
```bash
cd [PATH]              # Standard directory change
goto PATH              # Change directory with auto-save
goto '*pattern*'       # Fuzzy find and navigate
getmethere             # Interactive directory search
back                   # Navigate backward in history
pwd                    # Show current directory
```

### Process Management
```bash
name [INTENTION]       # Name current process with intention
save                   # Save current state
list                   # List all saved processes
jump PROCESS           # Jump to saved process location
state                  # Show current state (IPTP format)
```

### WiFi Hotspot
```bash
hotspot auto           # Smart enable (only if not on WiFi)
hotspot enable         # Force enable hotspot
hotspot disable        # Disable hotspot
hotspot status         # Show hotspot status
```

### DNS Router
```bash
dns start              # Start DNS router service
dns stop               # Stop DNS router service
dns status             # Show DNS router status
dns logs [N]           # Show last N DNS queries (default 10)
dns stats              # Show DNS statistics
dns install            # Show service installation instructions
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                 IPTP Shell (main.go)                │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────┐  ┌────────┐ │
│  │  Shell   │  │  State   │  │ DNS  │  │Hotspot │ │
│  │ (REPL)   │  │ (Field)  │  │Router│  │Manager │ │
│  │          │  │          │  │      │  │        │ │
│  │-Commands │  │-Processes│  │-Query│  │-Auto   │ │
│  │-Navigate │  │-History  │  │ Log  │  │ Enable │ │
│  │-Intents  │  │-Pulses   │  │-Stats│  │-Status │ │
│  └──────────┘  └──────────┘  └──────┘  └────────┘ │
│                                                     │
├─────────────────────────────────────────────────────┤
│              External State (JSON)                  │
│         /tmp/iptp_state.json                        │
│         /tmp/iptp_dns_queries.log                   │
└─────────────────────────────────────────────────────┘
```

### Key Concepts

**Design Nodes (DNs)**: Developer-implemented logic
**Objects**: Shared knowledge/state components
**Intentions**: Communication layer carrying Signals
**The Field**: Runtime state stored in database
**Pulses**: Trivalent truth values (Y/N/U)

## IPTP Philosophy

IPTP embodies a different approach to process coordination:

1. **Intentions over Commands**: Express what you want, not how to do it
2. **Shared Knowledge**: The Field stores state accessible to all processes
3. **Semantic Progression**: Progress based on meaning, not timing
4. **Trivalent Logic**: Yes/No/Unknown truth values for real-world uncertainty

### Example: Traditional vs IPTP

**Traditional Shell:**
```bash
./change_directory.sh /some/path
pwd  # Still in the same place! 😡
```

**IPTP Shell:**
```bash
[IPTP-1] ~$ goto /some/path
[IPTP-1] path$ # Actually here! 🎉
```

**Traditional Network Monitoring:**
```bash
# Check WiFi status
# Go to System Settings
# Enable Internet Sharing
# Configure SSID/password
# Start DNS server
# Check IP address
# Configure devices
```

**IPTP Network Monitoring:**
```bash
[IPTP-1] ~$ hotspot auto
# Done! Hotspot + DNS + monitoring ready!
```

## Use Cases

### 1. Instant Network Monitoring
```bash
sudo iptp
[IPTP-1] ~$ name "monitoring home network"
[monitoring_home_network] ~$ hotspot auto
# Family connects devices
[monitoring_home_network] ~$ dns logs 50
```

### 2. Development Session Management
```bash
# Morning session
[IPTP-1] ~$ name "working on API endpoints"
[API_endpoints] ~$ goto ~/projects/api
[API_endpoints] api$ save

# Lunch break - close terminal

# Afternoon session
iptp
[IPTP-1] ~$ jump API_endpoints
[API_endpoints] api$ # Back where you were!
```

### 3. IoT Device Analysis
```bash
[IPTP-1] ~$ hotspot auto
# Connect smart home device
[IPTP-1] ~$ dns logs 100
# See what cloud services it contacts
```

## File Structure

```
iptp/
├── main.go              # Entry point
├── shell.go             # Interactive REPL and commands
├── state.go             # State management (The Field)
├── dns_router.go        # DNS router implementation
├── hotspot.go           # WiFi hotspot manager (NEW!)
├── commands.go          # Command handlers
├── utils.go             # Utility functions
├── go.mod               # Go module dependencies
├── build.sh             # Build script
├── README.md            # This file
├── DNS_ROUTER.md        # DNS router documentation
├── HOTSPOT.md           # Hotspot documentation (NEW!)
├── QUICK_REFERENCE.md   # Quick reference card (NEW!)
└── dist/                # Compiled binaries
    ├── iptp-darwin-arm64
    ├── iptp-darwin-amd64
    ├── iptp-linux-amd64
    ├── iptp-linux-arm64
    └── iptp-windows-amd64.exe
```

## Development

### Adding New Commands

1. Add case to `executeCommand()` in `shell.go`
2. Implement command function (e.g., `cmdYourCommand()`)
3. Update help text in `cmdHelp()`

Example:
```go
case "mycommand":
    sh.cmdMyCommand(args)

func (sh *Shell) cmdMyCommand(args []string) {
    // Implementation
}
```

### Running Tests

```bash
go test ./...
```

### Dependencies

- `github.com/miekg/dns` - DNS library for router functionality

## Roadmap

### Phase 1: Core Shell (✓ Complete)
- [x] Basic REPL
- [x] Navigation commands
- [x] State persistence
- [x] Intention parsing

### Phase 2: Network Features (✓ Complete)
- [x] DNS query forwarding
- [x] Query logging
- [x] Service management
- [x] WiFi hotspot automation
- [x] Smart WiFi detection
- [x] Cross-platform support

### Phase 3: Enhanced Coordination (In Progress)
- [ ] Progressor loops
- [ ] Field-gated progression
- [ ] Signal payloads
- [ ] DN (Design Node) framework

### Phase 4: Ecosystem (Planned)
- [ ] `iptp serve` - Domain server
- [ ] `iptp auth` - Auth server
- [ ] Shared script dialect
- [ ] Built-in functions library
- [ ] DNS filtering/blacklisting
- [ ] Query analytics dashboard
- [ ] MAC address filtering for hotspot
- [ ] Bandwidth monitoring per device

## Contributing

This is a personal research project by Pronab at IntentixLab Keybyte Systems. Feedback and ideas welcome!

## Background

IPTP emerged from real-world frustrations with shell scripting limitations and draws on:
- NASA/Boeing RTOS engineering experience
- Unix design philosophy
- Real-time systems concepts
- Intention Space theory

The goal is to create infrastructure for coordinated processes that communicate through shared knowledge rather than direct calls - similar to how Express.js provides HTTP infrastructure for web development.

## License

MIT License - See LICENSE file

## Contact

**IntentixLab Keybyte Systems**
Melbourne, Victoria, Australia

---

*"Processes should coordinate through shared knowledge in the Intention Space, not through temporal sequences of direct calls."*
