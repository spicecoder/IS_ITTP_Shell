# 📡 IPTP WiFi Hotspot Manager

## What This Does

The IPTP WiFi Hotspot Manager automatically enables WiFi hotspot on your machine **only if you're not already connected to WiFi**. This makes it perfect for:

- 🏡 Sharing your internet connection at home
- 📱 Monitoring family devices via DNS
- 🔒 Network security analysis
- 🌐 Creating instant DNS monitoring setups

## The Problem It Solves

**Traditional approach:**
1. Manually check if you're on WiFi
2. Go to System Settings
3. Enable Internet Sharing
4. Configure SSID and password
5. Remember to disable later

**IPTP approach:**
```bash
[IPTP-1] ~$ hotspot auto
🔌 Checking WiFi connection status...
✓ Already connected to WiFi
  Hotspot not needed
```

Or if not on WiFi:
```bash
[IPTP-1] ~$ hotspot auto
🔌 Checking WiFi connection status...
✗ Not connected to WiFi
📱 Auto-enabling hotspot...
✓ Hotspot auto-enabled!
🌐 Auto-starting DNS router...
✓ DNS router started on 192.168.2.1
✨ Your network is ready!
```

## Commands

### Quick Reference

| Command | What It Does |
|---------|--------------|
| `hotspot auto` | Enable hotspot ONLY if not on WiFi |
| `hotspot enable` | Force enable hotspot |
| `hotspot disable` | Disable hotspot |
| `hotspot status` | Show current status |

### Detailed Usage

#### 1. Auto Mode (Recommended)

```bash
# Smart mode - only enables if needed
[IPTP-1] ~$ hotspot auto

# With custom settings
[IPTP-1] ~$ hotspot auto --ssid "MyNetwork" --password "secure123"
```

**When to use:** Every time! This is the safest option.

#### 2. Manual Enable

```bash
# Enable with defaults
[IPTP-1] ~$ hotspot enable

# Custom SSID and password
[IPTP-1] ~$ hotspot enable --ssid "IPTP-Net" --password "mypassword"

# Short form
[IPTP-1] ~$ hotspot enable -s "MyWiFi" -p "pass1234"
```

**When to use:** When you want to force hotspot even if on WiFi.

#### 3. Status Check

```bash
[IPTP-1] ~$ hotspot status
Status: ✓ ENABLED
  IP Address: 192.168.2.1

Devices should use this IP as their DNS server
  DNS Router: ✓ RUNNING
  DNS Queries: 247
```

#### 4. Disable

```bash
[IPTP-1] ~$ hotspot disable
📱 Disabling WiFi hotspot...
✓ Hotspot disabled
```

## Platform-Specific Behavior

### macOS

**What happens:**
- Uses AppleScript to enable Internet Sharing
- May require accessibility permissions first time
- Falls back to manual instructions if automatic fails

**First time setup:**
1. System Settings → Privacy & Security → Accessibility
2. Grant Terminal (or your terminal app) accessibility access

**Example output:**
```bash
[IPTP-1] ~$ hotspot enable
📱 Enabling WiFi Hotspot on macOS...

To enable Internet Sharing (WiFi Hotspot) on macOS:
1. Open System Settings
2. Go to General → Sharing
3. Enable 'Internet Sharing'
...
✓ Internet Sharing enabled!
```

### Linux

**What happens:**
- Uses NetworkManager (nmcli) to create hotspot
- Automatically configures SSID and password
- Works on most modern Linux distributions

**Requirements:**
- NetworkManager installed
- WiFi adapter supports AP mode

**Example output:**
```bash
[IPTP-1] ~$ hotspot enable
📱 Enabling WiFi Hotspot on Linux...
✓ WiFi hotspot enabled!
  SSID: IPTP-Hotspot
  Password: iptp123456
```

### Windows

**What happens:**
- Uses `netsh` to configure hosted network
- Creates virtual WiFi adapter
- Starts the hotspot service

**Requirements:**
- Administrator privileges
- WiFi adapter supports hosted network

**Example output:**
```bash
[IPTP-1] ~$ hotspot enable
📱 Enabling WiFi Hotspot on Windows...
✓ WiFi hotspot enabled!
  SSID: IPTP-Hotspot
  Password: iptp123456
```

## Complete Workflow Example

### Use Case: Monitor Home Network

```bash
# 1. Start IPTP (needs sudo for DNS)
$ sudo ./iptp

# 2. Set intention
[IPTP-1] ~$ name "monitoring home network"

# 3. Auto-enable hotspot + DNS in one command
[monitoring_home_network] ~$ hotspot auto
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

# 4. Connect devices to "IPTP-Hotspot"

# 5. Monitor queries
[monitoring_home_network] ~$ dns logs 20
=== Last 20 DNS Queries ===
[14:23:45] 192.168.2.10 -> youtube.com. (A) = 142.250.185.110
[14:23:46] 192.168.2.11 -> facebook.com. (A) = 157.240.22.35
...

# 6. Check overall status
[monitoring_home_network] ~$ hotspot status
Status: ✓ ENABLED
  IP Address: 192.168.2.1
  DNS Router: ✓ RUNNING
  DNS Queries: 247

# 7. Save state
[monitoring_home_network] ~$ save
✓ Saved state for: monitoring_home_network

# 8. Later, resume
[IPTP-1] ~$ jump monitoring_home_network
[monitoring_home_network] ~$ hotspot status
```

## Integration with DNS Router

The hotspot and DNS router work together seamlessly:

```bash
# Option 1: Manual steps
[IPTP-1] ~$ hotspot enable
[IPTP-1] ~$ dns start

# Option 2: Auto mode (does both!)
[IPTP-1] ~$ hotspot auto
# Automatically starts DNS router too!
```

When `hotspot auto` enables the hotspot, it also:
1. Starts the DNS router automatically
2. Configures it with your hotspot IP
3. Ready to log queries immediately

## Smart WiFi Detection

The hotspot manager intelligently checks if you're already on WiFi:

```bash
# Scenario 1: At home on WiFi
[IPTP-1] ~$ hotspot auto
🔌 Checking WiFi connection status...
✓ Already connected to WiFi
  Hotspot not needed
# No changes made!

# Scenario 2: Not on WiFi
[IPTP-1] ~$ hotspot auto
🔌 Checking WiFi connection status...
✗ Not connected to WiFi
📱 Auto-enabling hotspot...
# Hotspot enabled!

# Scenario 3: Force enable even on WiFi
[IPTP-1] ~$ hotspot enable
🔌 Checking WiFi connection status...
⚠️  You are currently connected to a WiFi network
   Enabling hotspot may disconnect you from the network

Continue anyway? (y/N): y
# Continues after confirmation
```

## Default Settings

When using `hotspot enable` or `hotspot auto` without options:

- **SSID:** `IPTP-Hotspot`
- **Password:** `iptp123456`
- **DNS:** Automatically set to your machine's IP
- **Upstream DNS:** 8.8.8.8 (Google DNS)

## Custom Configuration

```bash
# Custom SSID only
hotspot enable --ssid "FamilyNetwork"

# Custom password only
hotspot enable --password "strongpass123"

# Both custom
hotspot enable --ssid "MyNet" --password "secure456"

# Works with auto mode too
hotspot auto -s "AutoNet" -p "pass123"
```

## Troubleshooting

### "Permission denied" or "Failed to enable"

**Solution:** Run with sudo/administrator
```bash
# macOS/Linux
sudo ./iptp

# Windows
# Run as Administrator
```

### Hotspot enables but no internet on connected devices

**Check:**
1. Is internet sharing enabled?
2. Is the primary internet connection active?
3. Firewall blocking?

**macOS specific:**
```bash
# Verify Internet Sharing is on
# System Settings → Sharing → Internet Sharing
```

### Can't connect to hotspot

**Check:**
1. Correct password
2. WiFi adapter supports AP mode (Linux)
3. Hosted network supported (Windows)

### DNS not working for connected devices

**Check:**
1. DNS router is running: `dns status`
2. Devices using correct DNS (your machine's IP)
3. Port 53 not blocked by firewall

## Advanced: Service Integration

You can create a startup script that automatically:
1. Checks WiFi status
2. Enables hotspot if needed
3. Starts DNS router
4. Logs to a file

```bash
#!/bin/bash
# auto-network.sh

sudo ./iptp << EOF
name "auto network monitoring"
hotspot auto --ssid "HomeNet" --password "secure123"
save
exit
EOF
```

## IPTP Philosophy Integration

The hotspot manager follows IPTP principles:

**Intention-Based:**
```bash
name "sharing internet with guests"
hotspot auto
# The system understands your intention
```

**Smart Decisions:**
- Only enables hotspot when needed
- Asks confirmation if already on WiFi
- Auto-starts DNS for monitoring

**State Awareness:**
- Checks current WiFi status
- Integrates with DNS router state
- Saves configuration in The Field

**Coordination Through Knowledge:**
- Hotspot state visible to DNS router
- Both services coordinate through shared IP
- No direct process-to-process calls

## Security Considerations

**Default Password:**
- Default is `iptp123456` - **CHANGE THIS** for production use
- Use strong passwords for public spaces

**Recommendations:**
```bash
# Good password
hotspot enable --password "$(openssl rand -base64 12)"

# Or custom strong password
hotspot enable -s "MyNetwork" -p "Str0ng!P@ssw0rd#2024"
```

**DNS Logging:**
- All queries are logged - be mindful of privacy
- Inform users if monitoring family/guest devices
- Rotate logs regularly

## Tips & Tricks

**1. One-command network setup:**
```bash
sudo ./iptp
[IPTP-1] ~$ hotspot auto
# Done! Hotspot + DNS in one go
```

**2. Monitor specific device:**
```bash
[IPTP-1] ~$ dns logs | grep "192.168.2.10"
```

**3. See what kids are browsing:**
```bash
[IPTP-1] ~$ dns logs 100
# Look for their device IP
```

**4. Resume monitoring session:**
```bash
[IPTP-1] ~$ list
[IPTP-1] ~$ jump monitoring_home_network
[monitoring_home_network] ~$ hotspot status
```

**5. Quick disable everything:**
```bash
[IPTP-1] ~$ hotspot disable
[IPTP-1] ~$ dns stop
```

## Next Steps

After setting up hotspot:

1. **Test connectivity:** Connect a device and browse
2. **Verify DNS logging:** Check `dns logs`
3. **Monitor usage:** Watch `dns stats`
4. **Set up filtering:** (Future feature)
5. **Create alerts:** (Future feature)

## Future Enhancements

Planned features:
- 🎯 **Auto-reconnect:** Restore hotspot after wake from sleep
- 📊 **Bandwidth monitoring:** Track data usage per device
- 🔒 **MAC filtering:** Whitelist/blacklist devices
- ⏰ **Scheduled hotspot:** Auto-enable at specific times
- 📱 **Mobile app:** Control from your phone

---

**The IPTP Way:**

Traditional: "I need to enable WiFi hotspot"
IPTP: "I want to share my network and monitor devices"

```bash
name "sharing network with monitoring"
hotspot auto
# The system handles the rest!
```

**IntentixLab Keybyte Systems**
Melbourne, Victoria, Australia
