# Quick Start for macOS (Your Machine!)

## What You Have

✅ **iptp binary for Apple Silicon** - `iptp-darwin-arm64`  
✅ **Complete Go source code** - All in `iptp-go/` folder  
✅ **IPTP-compliant real shell** - Not a bash wrapper!  

## Installation (30 seconds)

```bash
# 1. Download the binary for your Mac
# (Use iptp-darwin-arm64 for M1/M2/M3, or iptp-darwin-amd64 for Intel)

# 2. Make it executable
chmod +x iptp-darwin-arm64

# 3. Move to your PATH (optional but recommended)
sudo mv iptp-darwin-arm64 /usr/local/bin/iptp

# 4. Run it!
iptp
```

## Or Run Directly

```bash
# From the downloads folder
cd ~/Downloads/iptp-go/dist
chmod +x iptp-darwin-arm64
./iptp-darwin-arm64
```

## Your First Session

```bash
$ iptp
🚀 iptp - IPTP Shell Process Manager
   Type 'help' for commands, 'exit' to quit

[shell_12345] ~$ name "working on my project"
✓ Shell named: my_project
  Intention: working on my project

[my_project] ~$ goto ~/Documents
✓ Changed to: /Users/pronab/Documents

[my_project] ~/Documents$ save
✓ Saved state for: my_project @ /Users/pronab/Documents

[my_project] ~/Documents$ list
=== Available Processes ===
  → my_project: /Users/pronab/Documents (PID: 12345)

[my_project] ~/Documents$ help
... (shows all commands)

[my_project] ~/Documents$ exit
Goodbye!
```


## Test the getmethere Feature

```bash
$ iptp
[shell]$ getmethere
Hello, pronab! Where would you like to work today?
Enter the directory name: Documents

Found 1 matching directories:
  1) /Users/pronab/Documents

Enter number (or 'q' to quit): 1
✓ Changed to: /Users/pronab/Documents
```


## Command Mode

You can also use iptp like a regular command:

```bash
# Name a process
$ iptp name "debugging auth"

# Navigate
$ iptp goto /var/log

# List processes
$ iptp list

# Show version
$ iptp version
```

## Building from Source (Optional)

If you want to modify and rebuild:

```bash
cd iptp-go
go build -o iptp
./iptp
```

## Your Mac Specs

You're on an **Apple Silicon Mac** (M1/M2/M3), so use:
- **Binary**: `iptp-darwin-arm64`

If you were on Intel Mac, you'd use:
- **Binary**: `iptp-darwin-amd64`

## Next Steps

1. **Try it now**: `./iptp-darwin-arm64`
2. **Name your shell**: `name "testing iptp"`
3. **Navigate around**: `goto ~/Documents`
4. **Try getmethere**: `getmethere`
5. **Run a script**: `./any-script.sh`
6. **Save state**: `save`
7. **List processes**: `list`

## Common Questions

**Q: Do I still need the bash version?**  
A: No! This Go binary replaces it completely.

**Q: Can I use this as my default shell?**  
A: Not yet - it's designed to run INSIDE your terminal. But you can always run `iptp` to enter it.

**Q: Does state survive reboots?**  
A: On macOS, `/tmp` is cleared on reboot. To persist, edit `main.go` to use `~/.iptp_state.json` instead.

**Q: Can I run this alongside the bash version?**  
A: Yes! They use the same state file format, so they're compatible.

## Getting Help

```bash
iptp help           # Show help
iptp version        # Show version
```

Inside iptp:
```bash
[shell]$ help         # Show commands
[shell]$ state        # Show current state
```

## State File Location

macOS: `/tmp/iptp_state.json`

View it:
```bash
cat /tmp/iptp_state.json | jq '.'
```
