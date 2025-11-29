# 🌐 IPTP DNS Router

## What This Does

The IPTP DNS Router is a built-in DNS server that runs on your machine and logs all DNS queries. Perfect for:
- 📱 Monitoring devices on your WiFi hotspot
- 🔍 Understanding what domains your network accesses
- 🛡️ Building custom DNS filtering (future feature)
- 📊 Network traffic analysis

## Quick Start

### 1. Start the DNS Router

```bash
# Start iptp shell (requires sudo for port 53)
sudo ./iptp

[IPTP-1] ~$ dns start
✓ DNS Router starting on 0.0.0.0:53
  Upstream: 8.8.8.8:53
  Log file: /tmp/iptp_dns_queries.log
✓ DNS router is now running
```

### 2. Configure Your Devices

**Option A: Configure Individual Device**
1. Go to WiFi settings on your phone/laptop
2. Find your machine's IP address (e.g., `192.168.1.100`)
3. Set DNS server to your machine's IP

**Option B: Configure Hotspot**
1. Enable WiFi hotspot on your machine
2. Devices connecting will automatically use your DNS

### 3. View DNS Queries

```bash
[IPTP-1] ~$ dns logs
=== Last 10 DNS Queries ===
[14:23:45] 192.168.1.50 -> google.com. (A) = 142.250.185.46
[14:23:46] 192.168.1.50 -> facebook.com. (A) = 157.240.22.35
[14:23:47] 192.168.1.51 -> youtube.com. (A) = 142.250.185.110

[IPTP-1] ~$ dns logs 50
# Shows last 50 queries

[IPTP-1] ~$ dns stats
=== DNS Router Statistics ===
Running: true
Total queries: 1247
Unique domains: 342
Listen address: 0.0.0.0:53
Upstream DNS: 8.8.8.8:53
Log file: /tmp/iptp_dns_queries.log
```

## Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `dns start` | Start DNS router | `dns start` |
| `dns start --listen ADDR` | Start with custom listen address | `dns start --listen 127.0.0.1:5353` |
| `dns start --upstream DNS` | Use custom upstream DNS | `dns start --upstream 1.1.1.1:53` |
| `dns stop` | Stop DNS router | `dns stop` |
| `dns status` | Show running status | `dns status` |
| `dns logs [N]` | Show last N queries (default 10) | `dns logs 100` |
| `dns stats` | Show statistics | `dns stats` |
| `dns install` | Show service installation | `dns install` |

## Advanced Usage

### Custom Upstream DNS

Use Cloudflare DNS instead of Google:
```bash
[IPTP-1] ~$ dns start --upstream 1.1.1.1:53
```

Use your ISP's DNS:
```bash
[IPTP-1] ~$ dns start --upstream 192.168.1.1:53
```

### Run on Different Port (No sudo required)

```bash
./iptp  # No sudo
[IPTP-1] ~$ dns start --listen 0.0.0.0:5353
```

Then configure devices to use port 5353 (not standard, most won't work).

### Install as System Service

The DNS router can run as a background service that starts on boot:

**Linux (systemd):**
```bash
[IPTP-1] ~$ dns install
# Follow the instructions shown
sudo systemctl start iptp-dns-router
sudo systemctl enable iptp-dns-router
```

**macOS (launchd):**
```bash
[IPTP-1] ~$ dns install
# Follow the instructions shown
sudo launchctl load /Library/LaunchDaemons/com.intentixlab.iptp.dnsrouter.plist
```

**Windows:**
```bash
[IPTP-1] ~$ dns install
# Follow the instructions shown
sc create IPTPDNSRouter binPath= "C:\path\to\iptp.exe dns start --service"
sc start IPTPDNSRouter
```

## Log File Format

The DNS router logs queries to `/tmp/iptp_dns_queries.log` (or `%TEMP%\iptp_dns_queries.log` on Windows) in JSON format:

```json
{"timestamp":"2025-11-29T14:23:45+11:00","client_ip":"192.168.1.50:54321","domain":"google.com.","query_type":"A","response":"142.250.185.46","upstream":"8.8.8.8:53"}
{"timestamp":"2025-11-29T14:23:46+11:00","client_ip":"192.168.1.50:54322","domain":"youtube.com.","query_type":"A","response":"142.250.185.110","upstream":"8.8.8.8:53"}
```

You can parse this with `jq` or any JSON tool:

```bash
# Count queries by domain
cat /tmp/iptp_dns_queries.log | jq -r '.domain' | sort | uniq -c | sort -rn

# Show queries from specific IP
cat /tmp/iptp_dns_queries.log | jq 'select(.client_ip | startswith("192.168.1.50"))'
```

## How It Works

```
[Your Device] → [IPTP DNS Router] → [Upstream DNS (8.8.8.8)]
                       ↓
              [Log to file + memory]
```

1. Device sends DNS query to IPTP DNS Router
2. Router logs: timestamp, client IP, domain, query type
3. Router forwards to upstream DNS (e.g., Google's 8.8.8.8)
4. Upstream DNS resolves and returns IP
5. Router logs the response
6. Router returns IP to device

## Troubleshooting

### "Permission denied" error

DNS runs on port 53, which requires root privileges:
```bash
sudo ./iptp
```

### DNS router won't start

Check if another DNS server is running:
```bash
# Linux/macOS
sudo lsof -i :53

# Stop systemd-resolved on Linux
sudo systemctl stop systemd-resolved
```

### Devices can't connect

1. Check firewall allows UDP port 53
2. Verify your machine's IP address
3. Test DNS resolution:
```bash
# From another device
nslookup google.com 192.168.1.100  # Your machine's IP
```

## Integration with IPTP Intentions

The DNS router follows IPTP philosophy:

```bash
[IPTP-1] ~$ name "monitoring network traffic"
✓ Shell named: monitoring_network_traffic

[monitoring_network_traffic] ~$ dns start
✓ DNS router is now running

[monitoring_network_traffic] ~$ save
✓ Saved state for: monitoring_network_traffic @ /home/user

# Later, resume the same context
[IPTP-2] ~$ jump monitoring_network_traffic
✓ Jumped to monitoring_network_traffic @ /home/user

[monitoring_network_traffic] ~$ dns status
Status: ✓ RUNNING
```

## Future Features (Roadmap)

- 🎯 Domain filtering / blacklisting
- 📈 Real-time query dashboard
- 🔒 DNSSEC validation
- 🌍 Custom DNS records (local DNS server)
- 📊 Query analytics and visualization
- 🔔 Alert on suspicious domains

## Why This Matters for IPTP

The DNS router is a perfect example of IPTP's "Intention Space" philosophy:

1. **Intention**: "I want to monitor my network"
2. **Pulse**: DNS queries create natural pulses in the system
3. **State**: Query logs are shared knowledge in the Field
4. **Coordination**: Devices coordinate through DNS resolution

This aligns with your vision of processes coordinating through shared knowledge rather than direct calls.
