# 🚀 IPTP DNS Router - Quick Setup

## For Your Mac (30 Seconds)

### Step 1: Build
```bash
cd ~/Downloads/iptp
chmod +x build.sh
./build.sh
```

### Step 2: Run with DNS Router
```bash
# DNS needs root for port 53
sudo ./dist/iptp-darwin-arm64
```

### Step 3: Start DNS Router
```bash
[IPTP-1] ~$ dns start
✓ DNS Router starting on 0.0.0.0:53
  Upstream: 8.8.8.8:53
  Log file: /tmp/iptp_dns_queries.log
✓ DNS router is now running
```

### Step 4: Configure WiFi Hotspot

1. **Enable WiFi Hotspot** on your Mac:
   - System Settings → General → Sharing
   - Enable "Internet Sharing"

2. **Get your Mac's IP**: 
   ```bash
   ifconfig en0 | grep "inet "
   # Example: 192.168.2.1
   ```

3. **Devices connecting to your hotspot will automatically use your DNS!**

## Monitor DNS Queries

```bash
# View last 20 queries
[IPTP-1] ~$ dns logs 20

# View statistics
[IPTP-1] ~$ dns stats

# Watch in real-time (tail the log file)
[IPTP-1] ~$ exit
$ tail -f /tmp/iptp_dns_queries.log | jq .
```

## Common Commands

```bash
dns start                    # Start DNS router
dns stop                     # Stop DNS router  
dns status                   # Check status
dns logs [N]                 # Show last N queries
dns stats                    # Show statistics
dns install                  # Install as service
```

## Use Cases

### 1. Monitor Your Kids' Devices
```bash
[IPTP-1] ~$ name "monitoring family devices"
[monitoring_family_devices] ~$ dns start
# Kids connect phones to your hotspot
[monitoring_family_devices] ~$ dns logs 50
# See what sites they're visiting
```

### 2. Debug Your IoT Devices
```bash
[IPTP-1] ~$ dns start
# Connect smart home devices to your hotspot
[IPTP-1] ~$ dns logs
# See what cloud services they're connecting to
```

### 3. Network Traffic Analysis
```bash
[IPTP-1] ~$ dns start
[IPTP-1] ~$ dns logs 100 > queries.txt
# Analyze patterns, find bandwidth hogs
```

## Advanced: Run as Background Service

```bash
# View installation instructions
[IPTP-1] ~$ dns install

# On macOS, this creates a launchd service
# On Linux, this creates a systemd service
# On Windows, this creates a Windows service
```

## Troubleshooting

**Q: "Permission denied" error**
A: DNS needs port 53, run with sudo:
```bash
sudo ./iptp
```

**Q: Can't see DNS queries**
A: Make sure devices are using your Mac's IP as DNS:
```bash
# On the device, check DNS settings
# Should show your Mac's IP (e.g., 192.168.2.1)
```

**Q: DNS router won't start**
A: Check if another DNS is running:
```bash
sudo lsof -i :53
# If something is there, stop it first
```

## What Gets Logged

Every DNS query creates a log entry:
```json
{
  "timestamp": "2025-11-29T14:23:45+11:00",
  "client_ip": "192.168.2.50:54321",
  "domain": "youtube.com.",
  "query_type": "A",
  "response": "142.250.185.110",
  "upstream": "8.8.8.8:53"
}
```

You can see:
- What device made the request (client_ip)
- What domain they wanted (domain)
- When they requested it (timestamp)
- What IP it resolved to (response)

## Integration with IPTP Philosophy

This DNS router is a perfect example of IPTP's approach:

**Traditional approach**: Program calls DNS, waits for response
**IPTP approach**: DNS queries create pulses in the Field

```bash
# Create an intention
[IPTP-1] ~$ name "network security analysis"

# Start the DNS router (creates a DN - Design Node)
[network_security_analysis] ~$ dns start

# Queries create pulses
# Each query updates the Field (state)
# Other processes can read from the Field
# No direct process-to-process calls needed!
```

## Files Included

- **dns_router.go** - DNS router implementation
- **shell.go** - Enhanced shell with DNS commands
- **commands.go** - Command handlers
- **state.go** - State management (The Field)
- **utils.go** - Utility functions
- **main.go** - Entry point
- **build.sh** - Build all platforms
- **README.md** - Complete documentation
- **DNS_ROUTER.md** - DNS router deep dive

## Next Steps

1. ✅ Build the project: `./build.sh`
2. ✅ Start with DNS: `sudo ./iptp`
3. ✅ Start DNS router: `dns start`
4. ✅ Share WiFi hotspot
5. ✅ Monitor queries: `dns logs`

Enjoy your IPTP shell with DNS routing! 🎉

---

**IntentixLab Keybyte Systems**
Melbourne, Australia
