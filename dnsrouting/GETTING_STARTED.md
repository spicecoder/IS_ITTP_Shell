# 🚀 Getting Started Checklist

## ✅ What You Have

**7 Source Files** (Ready to build):
- [x] main.go - Entry point
- [x] shell.go - Enhanced with DNS commands
- [x] dns_router.go - DNS server implementation
- [x] state.go - IPTP Field management
- [x] commands.go - Command execution
- [x] utils.go - Utility functions
- [x] go.mod - Dependencies

**6 Documentation Files**:
- [x] README.md - Complete project guide
- [x] DNS_ROUTER.md - DNS router documentation
- [x] QUICKSTART.md - 30-second setup
- [x] ARCHITECTURE.md - System diagrams
- [x] IMPLEMENTATION_SUMMARY.md - Implementation details
- [x] PROJECT_STRUCTURE.txt - File overview

**1 Build Script**:
- [x] build.sh - Cross-platform build

## 📋 Setup Steps (5 Minutes)

### Step 1: Download Files
```bash
# All files are in /mnt/user-data/outputs/
# Download them to your Mac
```

### Step 2: Prepare Environment
```bash
# Create project directory
mkdir ~/iptp-project
cd ~/iptp-project

# Copy all files there
# (drag and drop or use scp)
```

### Step 3: Install Go (if needed)
```bash
# Check if Go is installed
go version

# If not, download from: https://go.dev/dl/
# For Mac: brew install go
```

### Step 4: Build
```bash
# Make build script executable
chmod +x build.sh

# Build all platforms (or just your Mac)
./build.sh

# Or build for Mac only:
# go build -o iptp .
```

### Step 5: Test Basic Shell
```bash
# Run without sudo (no DNS)
./dist/iptp-darwin-arm64

# Try some commands
[IPTP-1] ~$ name "testing iptp"
[testing_iptp] ~$ goto ~/Downloads
[testing_iptp] Downloads$ pwd
[testing_iptp] Downloads$ exit
```

### Step 6: Test DNS Router
```bash
# Run with sudo (for DNS)
sudo ./dist/iptp-darwin-arm64

# Start DNS router
[IPTP-1] ~$ dns start
✓ DNS Router starting on 0.0.0.0:53

# Check status
[IPTP-1] ~$ dns status

# Test DNS resolution (from another terminal)
# nslookup google.com 127.0.0.1

# View logs
[IPTP-1] ~$ dns logs 10
```

## 🎯 First Real Use Case: Monitor Your Network

### Scenario: Track DNS queries from your WiFi hotspot

```bash
# 1. Start IPTP with DNS
sudo ./iptp

# 2. Set intention
[IPTP-1] ~$ name "monitoring home network"

# 3. Start DNS router
[monitoring_home_network] ~$ dns start
✓ DNS Router starting on 0.0.0.0:53

# 4. Enable WiFi hotspot on your Mac
#    System Settings → General → Sharing → Internet Sharing

# 5. Get your Mac's IP
[monitoring_home_network] ~$ exit
$ ifconfig en0 | grep "inet "
inet 192.168.2.1 netmask 0xffffff00 broadcast 192.168.2.255

# 6. Connect devices to your hotspot
#    They will automatically use your DNS (192.168.2.1)

# 7. Monitor queries
sudo ./iptp
[IPTP-1] ~$ jump monitoring_home_network
[monitoring_home_network] ~$ dns logs 50

# 8. View statistics
[monitoring_home_network] ~$ dns stats

# 9. Save state
[monitoring_home_network] ~$ save
```

## 🔍 What to Look For

**Success Indicators**:
- ✅ `dns start` shows "DNS Router starting"
- ✅ `dns status` shows "Status: ✓ RUNNING"
- ✅ `dns logs` shows query entries
- ✅ Log file exists: `/tmp/iptp_dns_queries.log`

**Troubleshooting**:
- ❌ "Permission denied" → Need sudo
- ❌ "Address already in use" → Port 53 busy (stop systemd-resolved)
- ❌ No queries showing → Check device DNS settings

## 📊 Understanding the Output

### DNS Logs Format
```
[14:23:45] 192.168.2.10 -> youtube.com. (A) = 142.250.185.110
   │           │              │         │         │
   │           │              │         │         └─ Resolved IP
   │           │              │         └─ Query Type (A=IPv4)
   │           │              └─ Domain requested
   │           └─ Client IP (device on your network)
   └─ Time (HH:MM:SS)
```

### DNS Stats Meaning
```
Total queries: 1247      # Total DNS requests handled
Unique domains: 342      # Different domains accessed
Running: true            # DNS router is active
```

## 🎓 Learning the Commands

### Essential Commands
```bash
# Process naming (sets intention)
name "your intention here"

# Navigation
goto ~/path/to/directory
getmethere               # Interactive search

# DNS router
dns start                # Start router
dns logs 20              # View queries
dns stats                # Statistics
dns stop                 # Stop router

# State management
save                     # Save current state
list                     # List saved processes
jump process_name        # Resume saved process
```

## 🔄 Next Session Workflow

```bash
# Start IPTP
sudo ./iptp

# Resume previous work
[IPTP-1] ~$ list
=== Available Processes ===
  → monitoring_home_network: /home/pronab (PID: 12345)

[IPTP-1] ~$ jump monitoring_home_network
[monitoring_home_network] ~$ 

# Check if DNS is running
[monitoring_home_network] ~$ dns status

# Start if needed
[monitoring_home_network] ~$ dns start

# Continue monitoring
[monitoring_home_network] ~$ dns logs 20
```

## 📖 Documentation Reading Order

1. **QUICKSTART.md** (5 min) - Start here
2. **README.md** (15 min) - Overview and features
3. **DNS_ROUTER.md** (10 min) - DNS router deep dive
4. **ARCHITECTURE.md** (10 min) - System design
5. **IMPLEMENTATION_SUMMARY.md** (10 min) - Technical details

## 🎁 Bonus: Useful Commands

### Analyze DNS logs with jq
```bash
# Most queried domains
cat /tmp/iptp_dns_queries.log | jq -r '.domain' | sort | uniq -c | sort -rn | head -10

# Queries from specific device
cat /tmp/iptp_dns_queries.log | jq 'select(.client_ip | startswith("192.168.2.10"))'

# Queries in last hour
cat /tmp/iptp_dns_queries.log | jq 'select(.timestamp | . > "'$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%S')'")'
```

### Real-time monitoring
```bash
# Watch logs live
tail -f /tmp/iptp_dns_queries.log | jq .

# Filter for specific domain
tail -f /tmp/iptp_dns_queries.log | jq 'select(.domain | contains("youtube"))'
```

## ✨ IPTP Philosophy in Action

This DNS router demonstrates IPTP concepts:

**Intention** → `name "monitoring home network"`
**The Field** → `/tmp/iptp_dns_queries.log`
**Pulses** → Each DNS query
**Design Node** → DNS router service
**Coordination** → Through shared logs, not direct calls

## 🚦 Status Check

After setup, you should be able to:
- [ ] Build the project successfully
- [ ] Run IPTP shell (without sudo)
- [ ] Start DNS router (with sudo)
- [ ] See DNS queries in logs
- [ ] Resume saved processes
- [ ] Monitor your network devices

## 🆘 Getting Help

If something doesn't work:
1. Check the troubleshooting section in DNS_ROUTER.md
2. Review IMPLEMENTATION_SUMMARY.md for technical details
3. Examine log files: `/tmp/iptp_state.json`, `/tmp/iptp_dns_queries.log`
4. Verify Go version: `go version` (need 1.21+)
5. Check port 53: `sudo lsof -i :53`

## 🎉 Success!

You now have:
✅ A working IPTP shell with DNS router
✅ Cross-platform binaries
✅ Complete documentation
✅ Ready for real-world use

**Next**: Start monitoring your network and explore the IPTP philosophy!

---

**Built with the IPTP Philosophy**
Small chunks, step-by-step, relaxed approach
NASA/Boeing RTOS thinking meets intention-driven design

**IntentixLab Keybyte Systems**
Melbourne, Victoria, Australia
