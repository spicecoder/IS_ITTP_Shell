package main

import (
	"encoding/json"
	"fmt"
	"log"
	//"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// DNSQuery represents a logged DNS query
type DNSQuery struct {
	Timestamp  string `json:"timestamp"`
	ClientIP   string `json:"client_ip"`
	Domain     string `json:"domain"`
	QueryType  string `json:"query_type"`
	Response   string `json:"response"`
	Upstream   string `json:"upstream"`
}

// DNSRouter manages the DNS routing service
type DNSRouter struct {
	server       *dns.Server
	logFile      string
	logMutex     sync.Mutex
	queries      []DNSQuery
	running      bool
	statusMutex  sync.RWMutex
	upstreamDNS  string
	listenAddr   string
}

// NewDNSRouter creates a new DNS router instance
func NewDNSRouter(listenAddr, upstreamDNS string) *DNSRouter {
	logPath := getDNSLogPath()
	
	return &DNSRouter{
		logFile:     logPath,
		queries:     make([]DNSQuery, 0),
		running:     false,
		upstreamDNS: upstreamDNS,
		listenAddr:  listenAddr,
	}
}

// getDNSLogPath returns platform-specific log file path
func getDNSLogPath() string {
	var logDir string
	
	switch runtime.GOOS {
	case "windows":
		logDir = os.Getenv("TEMP")
	case "darwin", "linux":
		logDir = "/tmp"
	default:
		logDir = "/tmp"
	}
	
	return filepath.Join(logDir, "iptp_dns_queries.log")
}

// handleDNSRequest processes incoming DNS queries
func (dr *DNSRouter) handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	// Extract query info
	var domain string
	var qtype string
	
	if len(r.Question) > 0 {
		domain = r.Question[0].Name
		qtype = dns.TypeToString[r.Question[0].Qtype]
	}

	clientIP := w.RemoteAddr().String()

	// Forward to upstream DNS
	c := new(dns.Client)
	upstream := dr.upstreamDNS
	
	in, _, err := c.Exchange(r, upstream)
	if err != nil {
		log.Printf("DNS query failed for %s: %v", domain, err)
		msg.Rcode = dns.RcodeServerFailure
		w.WriteMsg(msg)
		dr.logQuery(clientIP, domain, qtype, "FAILED", upstream)
		return
	}

	// Extract response IP(s)
	var responseIPs []string
	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			responseIPs = append(responseIPs, a.A.String())
		}
		if aaaa, ok := ans.(*dns.AAAA); ok {
			responseIPs = append(responseIPs, aaaa.AAAA.String())
		}
	}

	response := "NO_ANSWER"
	if len(responseIPs) > 0 {
		response = responseIPs[0] // Log first IP
	}

	// Log the query
	dr.logQuery(clientIP, domain, qtype, response, upstream)

	// Send response back to client
	w.WriteMsg(in)
}

// logQuery logs a DNS query to file and memory
func (dr *DNSRouter) logQuery(clientIP, domain, qtype, response, upstream string) {
	query := DNSQuery{
		Timestamp:  time.Now().Format(time.RFC3339),
		ClientIP:   clientIP,
		Domain:     domain,
		QueryType:  qtype,
		Response:   response,
		Upstream:   upstream,
	}

	dr.logMutex.Lock()
	defer dr.logMutex.Unlock()

	// Add to memory
	dr.queries = append(dr.queries, query)

	// Keep only last 1000 queries in memory
	if len(dr.queries) > 1000 {
		dr.queries = dr.queries[len(dr.queries)-1000:]
	}

	// Write to file
	f, err := os.OpenFile(dr.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer f.Close()

	logLine, _ := json.Marshal(query)
	fmt.Fprintf(f, "%s\n", logLine)
}

// Start starts the DNS router service
func (dr *DNSRouter) Start() error {
	dr.statusMutex.Lock()
	if dr.running {
		dr.statusMutex.Unlock()
		return fmt.Errorf("DNS router is already running")
	}
	dr.running = true
	dr.statusMutex.Unlock()

	dns.HandleFunc(".", dr.handleDNSRequest)
	
	dr.server = &dns.Server{
		Addr: dr.listenAddr,
		Net:  "udp",
	}

	fmt.Printf("✓ DNS Router starting on %s\n", dr.listenAddr)
	fmt.Printf("  Upstream: %s\n", dr.upstreamDNS)
	fmt.Printf("  Log file: %s\n", dr.logFile)

	go func() {
		if err := dr.server.ListenAndServe(); err != nil {
			log.Printf("DNS server error: %v", err)
			dr.statusMutex.Lock()
			dr.running = false
			dr.statusMutex.Unlock()
		}
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Stop stops the DNS router service
func (dr *DNSRouter) Stop() error {
	dr.statusMutex.Lock()
	defer dr.statusMutex.Unlock()

	if !dr.running {
		return fmt.Errorf("DNS router is not running")
	}

	if dr.server != nil {
		err := dr.server.Shutdown()
		dr.running = false
		return err
	}

	return nil
}

// IsRunning checks if the DNS router is running
func (dr *DNSRouter) IsRunning() bool {
	dr.statusMutex.RLock()
	defer dr.statusMutex.RUnlock()
	return dr.running
}

// GetRecentQueries returns recent DNS queries
func (dr *DNSRouter) GetRecentQueries(count int) []DNSQuery {
	dr.logMutex.Lock()
	defer dr.logMutex.Unlock()

	if count > len(dr.queries) {
		count = len(dr.queries)
	}

	if count == 0 {
		return []DNSQuery{}
	}

	return dr.queries[len(dr.queries)-count:]
}

// GetStats returns DNS router statistics
func (dr *DNSRouter) GetStats() map[string]interface{} {
	dr.logMutex.Lock()
	defer dr.logMutex.Unlock()

	uniqueDomains := make(map[string]bool)
	for _, q := range dr.queries {
		uniqueDomains[q.Domain] = true
	}

	return map[string]interface{}{
		"total_queries":   len(dr.queries),
		"unique_domains":  len(uniqueDomains),
		"running":         dr.IsRunning(),
		"listen_address":  dr.listenAddr,
		"upstream_dns":    dr.upstreamDNS,
		"log_file":        dr.logFile,
	}
}

// InstallService installs the DNS router as a system service
func (dr *DNSRouter) InstallService() error {
	switch runtime.GOOS {
	case "linux":
		return dr.installLinuxService()
	case "darwin":
		return dr.installMacService()
	case "windows":
		return dr.installWindowsService()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// installLinuxService creates a systemd service
func (dr *DNSRouter) installLinuxService() error {
	serviceName := "iptp-dns-router"
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	serviceContent := fmt.Sprintf(`[Unit]
Description=IPTP DNS Router Service
After=network.target

[Service]
Type=simple
ExecStart=%s dns start --service
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
`, execPath)

	fmt.Println("To install as systemd service, create:")
	fmt.Printf("  %s\n", servicePath)
	fmt.Println("\nWith content:")
	fmt.Println(serviceContent)
	fmt.Println("\nThen run:")
	fmt.Println("  sudo systemctl daemon-reload")
	fmt.Println("  sudo systemctl enable iptp-dns-router")
	fmt.Println("  sudo systemctl start iptp-dns-router")

	return nil
}

// installMacService creates a launchd service
func (dr *DNSRouter) installMacService() error {
	serviceName := "com.intentixlab.iptp.dnsrouter"
	plistPath := fmt.Sprintf("/Library/LaunchDaemons/%s.plist", serviceName)
	
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>dns</string>
        <string>start</string>
        <string>--service</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
`, serviceName, execPath)

	fmt.Println("To install as launchd service, create:")
	fmt.Printf("  %s\n", plistPath)
	fmt.Println("\nWith content:")
	fmt.Println(plistContent)
	fmt.Println("\nThen run:")
	fmt.Println("  sudo launchctl load " + plistPath)

	return nil
}

// installWindowsService provides instructions for Windows service
func (dr *DNSRouter) installWindowsService() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	fmt.Println("To install as Windows service, use sc.exe:")
	fmt.Printf("  sc create IPTPDNSRouter binPath= \"%s dns start --service\"\n", execPath)
	fmt.Println("  sc start IPTPDNSRouter")

	return nil
}
