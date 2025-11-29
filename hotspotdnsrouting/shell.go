package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// Global counter for unnamed shells
var iptpCounter int32 = 0

// Shell represents the interactive shell
type Shell struct {
	state          *State
	currentProcess string
	displayName    string // For prompt display
	reader         *bufio.Reader
	running        bool
	dnsRouter      *DNSRouter      // DNS router instance
	hotspot        *HotspotManager // Hotspot manager instance
}

// NewShell creates a new interactive shell
func NewShell(state *State) *Shell {
	processName := fmt.Sprintf("shell_%d", os.Getpid())
	
	// Generate IPTP-n name for display
	n := atomic.AddInt32(&iptpCounter, 1)
	displayName := fmt.Sprintf("IPTP-%d", n)
	
	// Initialize DNS router (default config, not started)
	dnsRouter := NewDNSRouter("0.0.0.0:53", "8.8.8.8:53")
	
	// Initialize hotspot manager
	hotspot := NewHotspotManager()
	
	return &Shell{
		state:          state,
		currentProcess: processName,
		displayName:    displayName,
		reader:         bufio.NewReader(os.Stdin),
		running:        true,
		dnsRouter:      dnsRouter,
		hotspot:        hotspot,
	}
}

// Run starts the REPL loop
func (sh *Shell) Run() {
	// Initialize current process
	currentDir, _ := os.Getwd()
	sh.state.SetProcess(sh.currentProcess, "Working in "+sh.currentProcess, currentDir)
	sh.state.Save()

	for sh.running {
		// Show prompt with just current directory name
		fmt.Printf("[%s] %s$ ", sh.displayName, sh.getCurrentDirName())

		// Read input
		line, err := sh.reader.ReadString('\n')
		if err != nil {
			break
		}

		// Parse and execute command
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		sh.executeCommand(line)
	}

	fmt.Println("\nGoodbye!")
}

// getCurrentDirName returns just the current directory name (not full path)
func (sh *Shell) getCurrentDirName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "?"
	}

	// Check if we're at home directory
	home, _ := os.UserHomeDir()
	if dir == home {
		return "~"
	}

	// Return just the directory name
	return filepath.Base(dir)
}

// executeCommand parses and executes a command
func (sh *Shell) executeCommand(line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "cd":
		// IMPORTANT: cd must be handled as a builtin
		sh.cmdGoto(args)
	case "name":
		sh.cmdName(args)
	case "goto":
		sh.cmdGoto(args)
	case "getmethere":
		sh.cmdGetMeThere()
	case "save":
		sh.cmdSave()
	case "list":
		sh.cmdList()
	case "jump":
		sh.cmdJump(args)
	case "back":
		sh.cmdBack()
	case "state":
		sh.cmdState()
	case "dns":
		sh.cmdDNS(args)
	case "hotspot":
		sh.cmdHotspot(args)
	case "help":
		sh.cmdHelp()
	case "pwd":
		// Show current directory
		dir, _ := os.Getwd()
		fmt.Println(dir)
	case "exit", "quit":
		sh.running = false
	default:
		// Try to execute as external command/script
		sh.cmdExec(parts)
	}
}

// cmdName handles the 'name' command
func (sh *Shell) cmdName(args []string) {
	if len(args) == 0 {
		fmt.Printf("Current process: %s\n", sh.displayName)
		if proc, ok := sh.state.GetProcess(sh.currentProcess); ok {
			fmt.Printf("Intention: %s\n", proc.Intention)
		}
		return
	}

	intention := strings.Join(args, " ")
	processName := ParseIntention(intention)

	// Update both internal and display names
	sh.currentProcess = processName
	sh.displayName = processName

	currentDir, _ := os.Getwd()
	sh.state.SetProcess(processName, intention, currentDir)
	sh.state.Save()

	fmt.Printf("✓ Shell named: %s\n", processName)
	fmt.Printf("  Intention: %s\n", intention)
	if intention != processName {
		fmt.Printf("  (parsed process name: %s)\n", processName)
	}
}

// cmdGoto handles the 'goto' and 'cd' commands
func (sh *Shell) cmdGoto(args []string) {
	if len(args) == 0 {
		// No arguments - cd to home directory (standard bash behavior)
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("✗ Cannot get home directory: %v\n", err)
			return
		}
		args = []string{home}
	}

	path := args[0]
	
	// Expand home directory
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			if path == "~" {
				path = home
			} else {
				path = filepath.Join(home, path[2:]) // Skip "~/"
			}
		}
	}
	
	// Handle wildcard matching
	if strings.Contains(path, "*") {
		matched, err := FindDirectoryFuzzy(path)
		if err != nil {
			fmt.Printf("✗ Error: %v\n", err)
			return
		}
		if matched == "" {
			fmt.Printf("✗ No match found for: %s\n", path)
			return
		}
		path = matched
		fmt.Printf("✓ Matched: %s\n", path)
	}

	oldDir, _ := os.Getwd()

	if err := os.Chdir(path); err != nil {
		fmt.Printf("✗ Cannot change directory: %v\n", err)
		return
	}

	newDir, _ := os.Getwd()
	sh.state.UpdateDirectory(sh.currentProcess, newDir, oldDir)
	sh.state.Save()

	// Silent like bash cd (no output on success)
}

// cmdGetMeThere handles interactive directory finding
func (sh *Shell) cmdGetMeThere() {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows
	}

	fmt.Printf("Hello, %s! Where would you like to work today?\n", username)
	fmt.Print("Enter the directory name: ")

	dirName, _ := sh.reader.ReadString('\n')
	dirName = strings.TrimSpace(dirName)

	if dirName == "" {
		fmt.Println("✗ No input provided")
		return
	}

	fmt.Println("Searching...")
	
	// First, search from current directory (fast)
	currentDir, _ := os.Getwd()
	dirs, err := FindDirectoriesFrom(dirName, currentDir, 3, 10)
	
	// If nothing found locally, search from home (slower)
	if len(dirs) == 0 {
		fmt.Println("Nothing found locally, searching from home directory...")
		homeDir, _ := os.UserHomeDir()
		dirs, err = FindDirectoriesFrom(dirName, homeDir, 4, 20)
	}
	
	if err != nil || len(dirs) == 0 {
		fmt.Printf("✗ No directories found matching '%s'\n", dirName)
		return
	}

	fmt.Printf("\nFound %d matching directories:\n", len(dirs))
	for i, dir := range dirs {
		// Show relative path if under current dir
		relPath, err := filepath.Rel(currentDir, dir)
		displayPath := dir
		if err == nil && !strings.HasPrefix(relPath, "..") {
			displayPath = "./" + relPath
		}
		fmt.Printf("  %d) %s\n", i+1, displayPath)
	}

	fmt.Print("\nEnter number (or 'q' to quit): ")
	choice, _ := sh.reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "q" || choice == "Q" {
		fmt.Println("Cancelled")
		return
	}

	var selectedIdx int
	if _, err := fmt.Sscanf(choice, "%d", &selectedIdx); err != nil || selectedIdx < 1 || selectedIdx > len(dirs) {
		fmt.Println("✗ Invalid selection")
		return
	}

	selectedDir := dirs[selectedIdx-1]
	oldDir, _ := os.Getwd()

	if err := os.Chdir(selectedDir); err != nil {
		fmt.Printf("✗ Cannot change directory: %v\n", err)
		return
	}

	newDir, _ := os.Getwd()
	sh.state.UpdateDirectory(sh.currentProcess, newDir, oldDir)
	sh.state.Save()

	fmt.Printf("✓ Changed to: %s\n", newDir)
}

// cmdSave saves the current state
func (sh *Shell) cmdSave() {
	currentDir, _ := os.Getwd()
	sh.state.UpdateDirectory(sh.currentProcess, currentDir, "")
	sh.state.Save()

	fmt.Printf("✓ Saved state for: %s @ %s\n", sh.displayName, currentDir)
}

// cmdList lists all saved processes
func (sh *Shell) cmdList() {
	processes := sh.state.ListProcesses()
	if len(processes) == 0 {
		fmt.Println("No saved processes")
		return
	}

	fmt.Println("=== Available Processes ===")
	for _, name := range processes {
		if proc, ok := sh.state.GetProcess(name); ok {
			fmt.Printf("  → %s: %s (PID: %d)\n", name, proc.CurrentDir, proc.PID)
		}
	}
}

// cmdJump jumps to a saved process location
func (sh *Shell) cmdJump(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: jump PROCESS")
		return
	}

	targetProcess := args[0]
	proc, ok := sh.state.GetProcess(targetProcess)
	if !ok {
		fmt.Printf("✗ Process '%s' not found\n", targetProcess)
		return
	}

	oldDir, _ := os.Getwd()

	if err := os.Chdir(proc.CurrentDir); err != nil {
		fmt.Printf("✗ Cannot change directory: %v\n", err)
		return
	}

	// Update both internal and display names
	sh.currentProcess = targetProcess
	sh.displayName = targetProcess
	
	newDir, _ := os.Getwd()
	sh.state.UpdateDirectory(sh.currentProcess, newDir, oldDir)
	sh.state.Save()

	fmt.Printf("✓ Jumped to %s @ %s\n", targetProcess, newDir)
}

// cmdBack goes back in navigation history
func (sh *Shell) cmdBack() {
	prevDir, ok := sh.state.PopHistory(sh.currentProcess)
	if !ok {
		fmt.Println("No history available")
		return
	}

	if err := os.Chdir(prevDir); err != nil {
		fmt.Printf("✗ Cannot change directory: %v\n", err)
		return
	}

	currentDir, _ := os.Getwd()
	sh.state.UpdateDirectory(sh.currentProcess, currentDir, "")
	sh.state.Save()

	fmt.Printf("✓ Back to: %s\n", currentDir)
}

// cmdState shows current process state
func (sh *Shell) cmdState() {
	proc, ok := sh.state.GetProcess(sh.currentProcess)
	if !ok {
		fmt.Println("No state for current process")
		return
	}

	currentDir, _ := os.Getwd()

	fmt.Println("=== Current Process State (IPTP Format) ===")
	fmt.Printf("Process: %s\n", sh.displayName)
	fmt.Printf("Directory: %s\n", currentDir)
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Println()
	fmt.Println("=== Intentions ===")
	fmt.Printf("  \"%s\"\n", proc.Intention)
	fmt.Println()
	fmt.Println("=== Pulses (Trivalent) ===")
	for _, pulse := range proc.Pulses {
		fmt.Printf("  {\"name\": \"%s\", \"TV\": \"%s\", \"response\": \"%s\"}\n",
			pulse.Name, pulse.TV, pulse.Response)
	}
}

// cmdHelp shows help information
func (sh *Shell) cmdHelp() {
	fmt.Println("iptp - IPTP Shell Process Manager")
	fmt.Println()
	fmt.Println("Navigation Commands:")
	fmt.Println("  cd [PATH]           - Change directory (standard command)")
	fmt.Println("  goto PATH           - Change directory with auto-save")
	fmt.Println("  goto '*pattern*'    - Fuzzy find and navigate")
	fmt.Println("  getmethere          - Interactive directory search")
	fmt.Println("  back                - Go back in navigation history")
	fmt.Println("  pwd                 - Show current directory")
	fmt.Println()
	fmt.Println("Process Management:")
	fmt.Println("  name [INTENTION]    - Name current process with intention")
	fmt.Println("  save                - Save current state")
	fmt.Println("  list                - List all saved processes")
	fmt.Println("  jump PROCESS        - Jump to saved process location")
	fmt.Println("  state               - Show current state (IPTP format)")
	fmt.Println()
	fmt.Println("DNS Router Commands:")
	fmt.Println("  dns start           - Start DNS router service")
	fmt.Println("  dns stop            - Stop DNS router service")
	fmt.Println("  dns status          - Show DNS router status")
	fmt.Println("  dns logs [N]        - Show last N DNS queries (default 10)")
	fmt.Println("  dns stats           - Show DNS statistics")
	fmt.Println("  dns install         - Show service installation instructions")
	fmt.Println()
	fmt.Println("WiFi Hotspot Commands:")
	fmt.Println("  hotspot enable      - Enable WiFi hotspot (with options)")
	fmt.Println("  hotspot disable     - Disable WiFi hotspot")
	fmt.Println("  hotspot status      - Show hotspot status")
	fmt.Println("  hotspot auto        - Auto-enable if not on WiFi")
	fmt.Println()
	fmt.Println("System Commands:")
	fmt.Println("  ls, mkdir, etc      - Any standard Unix command")
	fmt.Println("  ./script.sh         - Run any script in iptp context")
	fmt.Println("  help                - Show this help")
	fmt.Println("  exit                - Exit iptp")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  name \"working on authentication\"  # Set process name")
	fmt.Println("  cd wrk1                          # Standard cd works!")
	fmt.Println("  goto wrk2                        # Navigate with save")
	fmt.Println("  getmethere                       # Interactive search")
	fmt.Println("  jump test2                       # Jump to saved location")
	fmt.Println("  dns start                        # Start DNS router")
	fmt.Println("  dns logs 20                      # Show last 20 queries")
	fmt.Println("  hotspot auto                     # Auto-enable hotspot if needed")
	fmt.Println("  hotspot enable --ssid MyNet      # Custom hotspot name")
	fmt.Println()
	fmt.Println("Tips:")
	fmt.Println("  • New shells start with IPTP-1, IPTP-2, etc.")
	fmt.Println("  • Use 'name' to give your shell a meaningful name")
	fmt.Println("  • getmethere searches current dir first (fast!)")
	fmt.Println("  • cd works like normal bash/zsh")
	fmt.Println("  • goto auto-saves your location")
	fmt.Println("  • DNS router logs all queries to /tmp/iptp_dns_queries.log")
	fmt.Println("  • 'hotspot auto' enables hotspot only if not on WiFi")
	fmt.Println("  • Combine 'hotspot auto' + 'dns start' for instant monitoring")
	fmt.Println()
}

// cmdDNS handles DNS router management commands
func (sh *Shell) cmdDNS(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: dns [start|stop|status|logs|stats|install]")
		return
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "start":
		sh.dnsStart(subArgs)
	case "stop":
		sh.dnsStop()
	case "status":
		sh.dnsStatus()
	case "logs":
		sh.dnsLogs(subArgs)
	case "stats":
		sh.dnsStats()
	case "install":
		sh.dnsInstall()
	default:
		fmt.Printf("Unknown dns command: %s\n", subCmd)
		fmt.Println("Available: start, stop, status, logs, stats, install")
	}
}

// dnsStart starts the DNS router
func (sh *Shell) dnsStart(args []string) {
	if sh.dnsRouter.IsRunning() {
		fmt.Println("✗ DNS router is already running")
		fmt.Println("  Use 'dns stop' first if you want to restart")
		return
	}

	// Parse optional arguments for custom config
	listenAddr := "0.0.0.0:53"
	upstreamDNS := "8.8.8.8:53"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--listen", "-l":
			if i+1 < len(args) {
				listenAddr = args[i+1]
				i++
			}
		case "--upstream", "-u":
			if i+1 < len(args) {
				upstreamDNS = args[i+1]
				i++
			}
		}
	}

	// Update router config
	sh.dnsRouter = NewDNSRouter(listenAddr, upstreamDNS)

	if err := sh.dnsRouter.Start(); err != nil {
		fmt.Printf("✗ Failed to start DNS router: %v\n", err)
		fmt.Println("\nNote: DNS runs on port 53, which requires root/admin privileges")
		fmt.Println("Try: sudo iptp")
		return
	}

	fmt.Println("✓ DNS router is now running")
	fmt.Println("\nTo use this DNS router:")
	fmt.Println("  1. Go to WiFi/Network settings")
	fmt.Println("  2. Set DNS server to your machine's IP")
	fmt.Println("  3. Devices using your hotspot will route through this DNS")
}

// dnsStop stops the DNS router
func (sh *Shell) dnsStop() {
	if !sh.dnsRouter.IsRunning() {
		fmt.Println("DNS router is not running")
		return
	}

	if err := sh.dnsRouter.Stop(); err != nil {
		fmt.Printf("✗ Failed to stop DNS router: %v\n", err)
		return
	}

	fmt.Println("✓ DNS router stopped")
}

// dnsStatus shows DNS router status
func (sh *Shell) dnsStatus() {
	if sh.dnsRouter.IsRunning() {
		fmt.Println("Status: ✓ RUNNING")
		stats := sh.dnsRouter.GetStats()
		fmt.Printf("  Listen: %s\n", stats["listen_address"])
		fmt.Printf("  Upstream: %s\n", stats["upstream_dns"])
		fmt.Printf("  Queries: %d\n", stats["total_queries"])
		fmt.Printf("  Unique domains: %d\n", stats["unique_domains"])
	} else {
		fmt.Println("Status: ✗ STOPPED")
	}
}

// dnsLogs shows recent DNS query logs
func (sh *Shell) dnsLogs(args []string) {
	count := 10
	
	// Parse count argument
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &count)
	}

	queries := sh.dnsRouter.GetRecentQueries(count)
	
	if len(queries) == 0 {
		fmt.Println("No queries logged yet")
		return
	}

	fmt.Printf("=== Last %d DNS Queries ===\n", len(queries))
	for _, q := range queries {
		fmt.Printf("[%s] %s -> %s (%s) = %s\n",
			q.Timestamp[11:19], // Just show time HH:MM:SS
			q.ClientIP,
			q.Domain,
			q.QueryType,
			q.Response)
	}
	
	stats := sh.dnsRouter.GetStats()
	fmt.Printf("\nLog file: %s\n", stats["log_file"])
}

// dnsStats shows DNS statistics
func (sh *Shell) dnsStats() {
	stats := sh.dnsRouter.GetStats()
	
	fmt.Println("=== DNS Router Statistics ===")
	fmt.Printf("Running: %v\n", stats["running"])
	fmt.Printf("Total queries: %d\n", stats["total_queries"])
	fmt.Printf("Unique domains: %d\n", stats["unique_domains"])
	fmt.Printf("Listen address: %s\n", stats["listen_address"])
	fmt.Printf("Upstream DNS: %s\n", stats["upstream_dns"])
	fmt.Printf("Log file: %s\n", stats["log_file"])
}

// dnsInstall shows service installation instructions
func (sh *Shell) dnsInstall() {
	fmt.Println("=== Install DNS Router as System Service ===\n")
	
	if err := sh.dnsRouter.InstallService(); err != nil {
		fmt.Printf("✗ Error: %v\n", err)
	}
}

// cmdHotspot handles WiFi hotspot management commands
func (sh *Shell) cmdHotspot(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: hotspot [enable|disable|status|auto]")
		return
	}

	subCmd := args[0]

	switch subCmd {
	case "enable":
		sh.hotspotEnable(args[1:])
	case "disable":
		sh.hotspotDisable()
	case "status":
		sh.hotspotStatus()
	case "auto":
		sh.hotspotAuto(args[1:])
	default:
		fmt.Printf("Unknown hotspot command: %s\n", subCmd)
		fmt.Println("Available: enable, disable, status, auto")
	}
}

// hotspotEnable enables WiFi hotspot
func (sh *Shell) hotspotEnable(args []string) {
	// Default SSID and password
	ssid := "IPTP-Hotspot"
	password := "iptp123456"

	// Parse optional arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ssid", "-s":
			if i+1 < len(args) {
				ssid = args[i+1]
				i++
			}
		case "--password", "-p":
			if i+1 < len(args) {
				password = args[i+1]
				i++
			}
		}
	}

	fmt.Println("🔌 Checking WiFi connection status...")
	
	// Check if already connected to WiFi
	connected, err := sh.hotspot.IsConnectedToWiFi()
	if err != nil {
		fmt.Printf("✗ Error checking WiFi status: %v\n", err)
		return
	}

	if connected {
		fmt.Println("⚠️  You are currently connected to a WiFi network")
		fmt.Println("   Enabling hotspot may disconnect you from the network")
		fmt.Print("\nContinue anyway? (y/N): ")
		
		response, _ := sh.reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled")
			return
		}
	}

	fmt.Println("\n📱 Enabling WiFi hotspot...")
	
	if err := sh.hotspot.EnableHotspot(ssid, password); err != nil {
		fmt.Printf("✗ Failed to enable hotspot: %v\n", err)
		return
	}

	// Try to get the IP address
	ip, err := sh.hotspot.GetIPAddress()
	if err == nil {
		fmt.Printf("\n✓ Hotspot is ready!\n")
		fmt.Printf("  Your IP: %s\n", ip)
		fmt.Println("\nTo use the DNS router with this hotspot:")
		fmt.Printf("  1. Devices connect to: %s\n", ssid)
		fmt.Printf("  2. They will use DNS: %s\n", ip)
		fmt.Println("  3. Start DNS router: dns start")
		fmt.Println("  4. Monitor queries: dns logs")
	}
}

// hotspotDisable disables WiFi hotspot
func (sh *Shell) hotspotDisable() {
	fmt.Println("📱 Disabling WiFi hotspot...")
	
	if err := sh.hotspot.DisableHotspot(); err != nil {
		fmt.Printf("✗ Failed to disable hotspot: %v\n", err)
		return
	}

	fmt.Println("✓ Hotspot disabled")
}

// hotspotStatus shows hotspot status
func (sh *Shell) hotspotStatus() {
	enabled, err := sh.hotspot.GetHotspotStatus()
	if err != nil {
		fmt.Printf("✗ Error checking hotspot status: %v\n", err)
		return
	}

	if enabled {
		fmt.Println("Status: ✓ ENABLED")
		
		ip, err := sh.hotspot.GetIPAddress()
		if err == nil {
			fmt.Printf("  IP Address: %s\n", ip)
			fmt.Println("\nDevices should use this IP as their DNS server")
			
			// Check if DNS router is running
			if sh.dnsRouter.IsRunning() {
				fmt.Println("  DNS Router: ✓ RUNNING")
				stats := sh.dnsRouter.GetStats()
				fmt.Printf("  DNS Queries: %d\n", stats["total_queries"])
			} else {
				fmt.Println("  DNS Router: ✗ NOT RUNNING")
				fmt.Println("\nTip: Start DNS router with 'dns start'")
			}
		}
	} else {
		fmt.Println("Status: ✗ DISABLED")
		fmt.Println("\nTo enable: hotspot enable")
	}
}

// hotspotAuto automatically enables hotspot if not connected to WiFi
func (sh *Shell) hotspotAuto(args []string) {
	// Default SSID and password
	ssid := "IPTP-Hotspot"
	password := "iptp123456"

	// Parse optional arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--ssid", "-s":
			if i+1 < len(args) {
				ssid = args[i+1]
				i++
			}
		case "--password", "-p":
			if i+1 < len(args) {
				password = args[i+1]
				i++
			}
		}
	}

	fmt.Println("🔌 Checking WiFi connection status...")
	
	connected, err := sh.hotspot.IsConnectedToWiFi()
	if err != nil {
		fmt.Printf("✗ Error checking WiFi status: %v\n", err)
		return
	}

	if connected {
		fmt.Println("✓ Already connected to WiFi")
		fmt.Println("  Hotspot not needed")
		return
	}

	fmt.Println("✗ Not connected to WiFi")
	fmt.Println("📱 Auto-enabling hotspot...\n")
	
	if err := sh.hotspot.EnableHotspot(ssid, password); err != nil {
		fmt.Printf("✗ Failed to enable hotspot: %v\n", err)
		return
	}

	// Try to get the IP address
	ip, err := sh.hotspot.GetIPAddress()
	if err == nil {
		fmt.Printf("\n✓ Hotspot auto-enabled!\n")
		fmt.Printf("  SSID: %s\n", ssid)
		fmt.Printf("  Your IP: %s\n", ip)
		
		// Auto-start DNS router too
		fmt.Println("\n🌐 Auto-starting DNS router...")
		if err := sh.dnsRouter.Start(); err != nil {
			fmt.Printf("⚠️  DNS router failed: %v\n", err)
			fmt.Println("   (You may need to run with sudo)")
		} else {
			fmt.Printf("✓ DNS router started on %s\n", ip)
			fmt.Println("\n✨ Your network is ready!")
			fmt.Println("   Devices can now connect and their DNS queries will be logged")
		}
	}
}

// cmdExec executes external commands/scripts
func (sh *Shell) cmdExec(parts []string) {
	ExecuteScript(parts)
}