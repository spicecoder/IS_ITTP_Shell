# 🚀 IPTP Quick Reference - Hotspot + DNS

## One-Line Network Monitoring Setup

```bash
sudo ./iptp
[IPTP-1] ~$ hotspot auto
```

**That's it!** Your machine now:
- ✅ Shares internet via WiFi hotspot (if not already on WiFi)
- ✅ Logs all DNS queries from connected devices
- ✅ Ready for monitoring

## Essential Commands

### Setup (Do Once)
```bash
# 1. Build
./build.sh

# 2. Run with sudo (for DNS port 53)
sudo ./iptp
```

### Quick Start (Every Session)
```bash
# Smart auto-mode (recommended)
[IPTP-1] ~$ hotspot auto

# OR step-by-step
[IPTP-1] ~$ hotspot enable
[IPTP-1] ~$ dns start
```

### Monitoring
```bash
# View recent queries
dns logs 20

# Check status
hotspot status
dns status

# See statistics
dns stats
```

### Cleanup
```bash
# Disable everything
hotspot disable
dns stop
exit
```

## Common Workflows

### 1. Monitor Home Network
```bash
sudo ./iptp
[IPTP-1] ~$ name "home monitoring"
[home_monitoring] ~$ hotspot auto
[home_monitoring] ~$ dns logs
```

### 2. Check What Kids Are Browsing
```bash
[IPTP-1] ~$ hotspot auto
# Kids connect phones
[IPTP-1] ~$ dns logs 50 | grep "192.168.2.10"
```

### 3. Track IoT Device Behavior
```bash
[IPTP-1] ~$ hotspot auto
# Connect smart device
[IPTP-1] ~$ dns logs 100
# See what servers it contacts
```

## Hotspot Commands

| Command | Action |
|---------|--------|
| `hotspot auto` | Enable only if not on WiFi (smart) |
| `hotspot enable` | Force enable |
| `hotspot disable` | Turn off |
| `hotspot status` | Check current state |

### Custom Settings
```bash
hotspot enable --ssid "MyNetwork" --password "secure123"
hotspot auto -s "FamilyNet" -p "pass456"
```

## DNS Commands

| Command | Action |
|---------|--------|
| `dns start` | Start DNS router |
| `dns stop` | Stop DNS router |
| `dns logs [N]` | Show last N queries |
| `dns stats` | Show statistics |
| `dns status` | Check if running |

## Default Settings

**Hotspot:**
- SSID: `IPTP-Hotspot`
- Password: `iptp123456`
- Your IP: Usually `192.168.2.1`

**DNS:**
- Listen: `0.0.0.0:53`
- Upstream: `8.8.8.8:53`
- Log: `/tmp/iptp_dns_queries.log`

## Platform Notes

### macOS
- First time: Grant accessibility permissions
- May need manual Internet Sharing setup
- Works on Apple Silicon & Intel

### Linux  
- Requires NetworkManager (nmcli)
- WiFi adapter must support AP mode
- Works on most distributions

### Windows
- Run as Administrator
- Uses hosted network (netsh)
- Adapter must support hosted network

## Troubleshooting

**"Permission denied"**
→ Run with `sudo ./iptp`

**Hotspot won't enable**
→ Check if WiFi adapter supports AP/hotspot mode

**DNS not logging**
→ Verify with `dns status`, check if started

**Devices can't connect**
→ Check password, verify hotspot is enabled

**No internet on devices**
→ Enable Internet Sharing (macOS) or IP forwarding (Linux)

## Quick Tips

💡 **Use auto mode:** `hotspot auto` is safest - only enables if needed

💡 **One command setup:** `hotspot auto` enables hotspot AND starts DNS

💡 **Check before enabling:** Auto mode checks WiFi status first

💡 **Custom names:** Use `--ssid` and `--password` for custom config

💡 **Monitor specific device:** `dns logs | grep "DEVICE_IP"`

## Example Session

```bash
$ sudo ./iptp
🚀 iptp- IPTP Shell Process Manager
   Type 'help' for commands, 'exit' to quit

[IPTP-1] ~$ name "monitoring family devices"
✓ Shell named: monitoring_family_devices

[monitoring_family_devices] ~$ hotspot auto
🔌 Checking WiFi connection status...
✗ Not connected to WiFi
📱 Auto-enabling hotspot...
✓ Hotspot auto-enabled!
  SSID: IPTP-Hotspot
  Your IP: 192.168.2.1

🌐 Auto-starting DNS router...
✓ DNS router started on 192.168.2.1

✨ Your network is ready!

# Family members connect their devices...

[monitoring_family_devices] ~$ dns logs 10
=== Last 10 DNS Queries ===
[14:23:45] 192.168.2.10 -> youtube.com. (A) = 142.250.185.110
[14:23:46] 192.168.2.11 -> facebook.com. (A) = 157.240.22.35
[14:23:47] 192.168.2.10 -> instagram.com. (A) = 157.240.22.174

[monitoring_family_devices] ~$ hotspot status
Status: ✓ ENABLED
  IP Address: 192.168.2.1
  DNS Router: ✓ RUNNING
  DNS Queries: 247

[monitoring_family_devices] ~$ save
✓ Saved state for: monitoring_family_devices

[monitoring_family_devices] ~$ exit
Goodbye!
```

## IPTP Philosophy

**Traditional Approach:**
1. Check if on WiFi manually
2. Go to System Settings
3. Enable Internet Sharing
4. Configure SSID/password
5. Start DNS server separately
6. Check IP address
7. Configure DNS on devices

**IPTP Approach:**
```bash
hotspot auto
```

That's the power of **intention-driven design** + **smart automation**.

---

**Need Help?**
- Full docs: [HOTSPOT.md](HOTSPOT.md)
- DNS guide: [DNS_ROUTER.md](DNS_ROUTER.md)
- General help: `help` command

**IntentixLab Keybyte Systems**
Melbourne, Victoria, Australia
