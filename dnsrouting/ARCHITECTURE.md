# IPTP Architecture with DNS Router

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         IPTP Shell                              │
│                     (main.go + shell.go)                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ├──────────────┬──────────────┐
                              │              │              │
                              ▼              ▼              ▼
                    ┌─────────────┐  ┌──────────┐  ┌─────────────┐
                    │   Commands  │  │  State   │  │ DNS Router  │
                    │  (REPL/CLI) │  │ (Field)  │  │  (Service)  │
                    └─────────────┘  └──────────┘  └─────────────┘
                                            │              │
                                            ▼              ▼
                                    ┌──────────┐  ┌──────────────┐
                                    │   JSON   │  │  DNS Queries │
                                    │  State   │  │     Log      │
                                    └──────────┘  └──────────────┘
```

## DNS Router Flow

```
┌──────────────┐
│ Your Device  │
│ (Phone/PC)   │
└──────┬───────┘
       │ DNS Query: "youtube.com"?
       │
       ▼
┌─────────────────────────────────────────┐
│     IPTP DNS Router (Port 53)           │
│                                         │
│  1. Receive Query                       │
│  2. Log: timestamp, client, domain      │
│  3. Forward to Upstream (8.8.8.8)       │
│  4. Receive Response                    │
│  5. Log: response IP                    │
│  6. Return to Client                    │
└─────────────┬───────────────────────────┘
              │
              ├──────────┬────────────────┐
              │          │                │
              ▼          ▼                ▼
       ┌──────────┐  ┌─────────┐  ┌──────────┐
       │ Upstream │  │   Log   │  │  Stats   │
       │   DNS    │  │  File   │  │ (Memory) │
       │ 8.8.8.8  │  │  JSON   │  │          │
       └──────────┘  └─────────┘  └──────────┘
```

## State Management (The Field)

```
┌──────────────────────────────────────────────────────┐
│                  The Field                           │
│              (/tmp/iptp_state.json)                  │
├──────────────────────────────────────────────────────┤
│                                                      │
│  Processes: {                                        │
│    "authentication_module": {                        │
│      "intention": "working on auth",                 │
│      "current_dir": "/home/user/code/auth",          │
│      "history": ["/home/user", "/home/user/code"],   │
│      "pulses": [                                     │
│        {"name": "process named", "TV": "Y"},         │
│        {"name": "directory saved", "TV": "Y"}        │
│      ]                                               │
│    }                                                 │
│  }                                                   │
└──────────────────────────────────────────────────────┘
```

## DNS Query Log Structure

```
┌──────────────────────────────────────────────────────┐
│           DNS Queries Log                            │
│      (/tmp/iptp_dns_queries.log)                     │
├──────────────────────────────────────────────────────┤
│                                                      │
│  {"timestamp": "2025-11-29T14:23:45+11:00",          │
│   "client_ip": "192.168.1.50:54321",                 │
│   "domain": "youtube.com.",                          │
│   "query_type": "A",                                 │
│   "response": "142.250.185.110",                     │
│   "upstream": "8.8.8.8:53"}                          │
│                                                      │
│  {"timestamp": "2025-11-29T14:23:46+11:00",          │
│   "client_ip": "192.168.1.50:54322",                 │
│   "domain": "google.com.",                           │
│   "query_type": "A",                                 │
│   "response": "142.250.185.46",                      │
│   "upstream": "8.8.8.8:53"}                          │
└──────────────────────────────────────────────────────┘
```

## Process Lifecycle

```
┌─────────────────────────────────────────────────────────┐
│                  1. Shell Starts                        │
│               [IPTP-1] ~$                               │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             2. User Names Process                       │
│     [IPTP-1] ~$ name "working on auth"                  │
│     → Process: authentication_module                    │
│     → Intention: "working on auth"                      │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             3. User Navigates                           │
│     [authentication_module] ~$ goto ~/code/auth         │
│     → Current Dir: /home/user/code/auth                 │
│     → State Auto-saved to Field                         │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             4. User Starts DNS Router                   │
│     [authentication_module] auth$ dns start             │
│     → DNS Router running on 0.0.0.0:53                  │
│     → Logging enabled                                   │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             5. User Monitors DNS                        │
│     [authentication_module] auth$ dns logs 20           │
│     → Shows last 20 DNS queries                         │
│     → From devices using this DNS                       │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             6. User Saves and Exits                     │
│     [authentication_module] auth$ save                  │
│     [authentication_module] auth$ exit                  │
│     → State persisted to /tmp/iptp_state.json           │
└─────────────────┬───────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────┐
│             7. Resume Later                             │
│     [IPTP-1] ~$ jump authentication_module              │
│     [authentication_module] auth$                       │
│     → Back in same directory                            │
│     → DNS router can be restarted                       │
└─────────────────────────────────────────────────────────┘
```

## WiFi Hotspot Setup

```
┌──────────────────────────────────────────────────┐
│              Your Mac/Linux PC                   │
│                                                  │
│  ┌────────────────────────────────┐              │
│  │     IPTP Shell + DNS Router    │              │
│  │       (Port 53 - UDP)          │              │
│  └────────────────┬───────────────┘              │
│                   │                              │
│  ┌────────────────▼───────────────┐              │
│  │     WiFi Hotspot Enabled       │              │
│  │      SSID: "MyHotspot"         │              │
│  │      IP: 192.168.2.1           │              │
│  │      DNS: 192.168.2.1 (self)   │              │
│  └────────────────┬───────────────┘              │
└───────────────────┼──────────────────────────────┘
                    │
        ┌───────────┴──────────────────┐
        │                              │
        ▼                              ▼
┌───────────────┐              ┌───────────────┐
│  Phone/Tablet │              │  Laptop/PC    │
│               │              │               │
│ Connected to  │              │ Connected to  │
│ "MyHotspot"   │              │ "MyHotspot"   │
│               │              │               │
│ DNS queries → │              │ DNS queries → │
│ 192.168.2.1   │              │ 192.168.2.1   │
└───────────────┘              └───────────────┘
        │                              │
        └──────────┬───────────────────┘
                   │
                   ▼
          ┌────────────────────┐
          │  All DNS queries   │
          │  logged by IPTP    │
          └────────────────────┘
```

## Command Flow Example

```
User Input: dns start
     │
     ▼
executeCommand(line)
     │
     ├─ Parse: cmd="dns", args=["start"]
     │
     ▼
cmdDNS(args)
     │
     ├─ Switch on args[0]
     │
     ▼
dnsStart(args[1:])
     │
     ├─ Create DNSRouter instance
     ├─ Configure listen address (0.0.0.0:53)
     ├─ Configure upstream DNS (8.8.8.8:53)
     │
     ▼
DNSRouter.Start()
     │
     ├─ Register DNS handler
     ├─ Start UDP server on port 53
     ├─ Set running = true
     │
     ▼
handleDNSRequest(w, r)
     │
     ├─ Parse query
     ├─ Log query details
     ├─ Forward to upstream
     ├─ Receive response
     ├─ Log response
     ├─ Return to client
     │
     ▼
logQuery(...)
     │
     ├─ Add to memory buffer
     ├─ Append to JSON log file
     │
     ▼
Complete ✓
```

## Key Files and Their Roles

```
main.go
  └─ Entry point
  └─ Initialize state
  └─ Start shell

shell.go
  ├─ REPL loop
  ├─ Command parsing
  ├─ Navigation (cd, goto, getmethere)
  ├─ Process management (name, jump, list)
  └─ DNS commands (dns start/stop/logs/stats)

dns_router.go
  ├─ DNS server implementation
  ├─ Query logging
  ├─ Statistics tracking
  └─ Service installation

state.go
  ├─ Process state management
  ├─ JSON persistence
  ├─ History tracking
  └─ Pulse management

commands.go
  └─ Non-interactive command handlers

utils.go
  ├─ Directory search
  ├─ Intention parsing
  └─ Path manipulation
```

## IPTP Concepts Mapping

```
Traditional          IPTP
───────────         ──────────────────────
Process calls  →    Intentions
Variables      →    Objects in Field
Return values  →    Pulses (Y/N/U)
Functions      →    Design Nodes (DNs)
Program flow   →    Semantic Progression
State          →    The Field
Events         →    Signals
```

---

This architecture enables:
- **Separation of concerns**: DNS routing is independent
- **State persistence**: The Field survives shell restarts
- **Natural coordination**: Processes share knowledge
- **Intention-based**: Express what, not how
