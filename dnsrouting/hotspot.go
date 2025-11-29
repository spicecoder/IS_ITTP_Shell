package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// HotspotManager manages WiFi hotspot functionality
type HotspotManager struct {
	platform string
}

// NewHotspotManager creates a new hotspot manager
func NewHotspotManager() *HotspotManager {
	return &HotspotManager{
		platform: runtime.GOOS,
	}
}

// IsConnectedToWiFi checks if the machine is currently connected to WiFi
func (hm *HotspotManager) IsConnectedToWiFi() (bool, error) {
	switch hm.platform {
	case "darwin":
		return hm.isConnectedMac()
	case "linux":
		return hm.isConnectedLinux()
	case "windows":
		return hm.isConnectedWindows()
	default:
		return false, fmt.Errorf("unsupported platform: %s", hm.platform)
	}
}

// EnableHotspot enables WiFi hotspot on the current machine
func (hm *HotspotManager) EnableHotspot(ssid, password string) error {
	switch hm.platform {
	case "darwin":
		return hm.enableHotspotMac(ssid, password)
	case "linux":
		return hm.enableHotspotLinux(ssid, password)
	case "windows":
		return hm.enableHotspotWindows(ssid, password)
	default:
		return fmt.Errorf("unsupported platform: %s", hm.platform)
	}
}

// DisableHotspot disables WiFi hotspot
func (hm *HotspotManager) DisableHotspot() error {
	switch hm.platform {
	case "darwin":
		return hm.disableHotspotMac()
	case "linux":
		return hm.disableHotspotLinux()
	case "windows":
		return hm.disableHotspotWindows()
	default:
		return fmt.Errorf("unsupported platform: %s", hm.platform)
	}
}

// GetHotspotStatus checks if hotspot is currently enabled
func (hm *HotspotManager) GetHotspotStatus() (bool, error) {
	switch hm.platform {
	case "darwin":
		return hm.getHotspotStatusMac()
	case "linux":
		return hm.getHotspotStatusLinux()
	case "windows":
		return hm.getHotspotStatusWindows()
	default:
		return false, fmt.Errorf("unsupported platform: %s", hm.platform)
	}
}

// GetIPAddress gets the machine's IP address for the hotspot interface
func (hm *HotspotManager) GetIPAddress() (string, error) {
	switch hm.platform {
	case "darwin":
		return hm.getIPAddressMac()
	case "linux":
		return hm.getIPAddressLinux()
	case "windows":
		return hm.getIPAddressWindows()
	default:
		return "", fmt.Errorf("unsupported platform: %s", hm.platform)
	}
}

// ===== macOS Implementation =====

func (hm *HotspotManager) isConnectedMac() (bool, error) {
	// Simple and reliable: Check if en0 (standard WiFi interface) has an IP address
	// This works even when networksetup/airport commands fail
	cmd := exec.Command("sh", "-c", "ifconfig en0 | grep 'inet ' | grep -v 'inet6' | grep -v '127.0.0.1'")
	output, err := cmd.Output()
	
	// If we got output, en0 has an IP address = connected
	if err == nil && len(strings.TrimSpace(string(output))) > 0 {
		return true, nil
	}
	
	// Fallback: check other potential WiFi interfaces (en1, en2, en3)
	for i := 1; i <= 3; i++ {
		iface := fmt.Sprintf("en%d", i)
		cmd = exec.Command("sh", "-c", fmt.Sprintf("ifconfig %s | grep 'inet ' | grep -v 'inet6' | grep -v '127.0.0.1'", iface))
		output, err = cmd.Output()
		
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			return true, nil
		}
	}
	
	// No WiFi interface has an IP
	return false, nil
}

func (hm *HotspotManager) enableHotspotMac(ssid, password string) error {
	// On macOS, we need to use System Settings
	// Internet Sharing can be enabled via sharing preferences
	
	fmt.Println("📱 Enabling WiFi Hotspot on macOS...")
	fmt.Println()
	fmt.Println("To enable Internet Sharing (WiFi Hotspot) on macOS:")
	fmt.Println("1. Open System Settings")
	fmt.Println("2. Go to General → Sharing")
	fmt.Println("3. Enable 'Internet Sharing'")
	fmt.Println("4. Share from: Ethernet or your internet connection")
	fmt.Println("5. To computers using: Wi-Fi")
	fmt.Println("6. Click 'Wi-Fi Options' to set:")
	fmt.Printf("   Network Name: %s\n", ssid)
	fmt.Printf("   Password: %s\n", password)
	fmt.Println()
	fmt.Println("Or use this quick AppleScript method:")
	fmt.Println()
	
	// Try using AppleScript to enable Internet Sharing
	script := `tell application "System Events"
		tell network preferences
			set internetSharing to service "Internet Sharing" of sharing preferences
			set enabled of internetSharing to true
		end tell
	end tell`
	
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		fmt.Println("Note: Automatic enabling requires accessibility permissions.")
		fmt.Println("You may need to enable it manually in System Settings.")
		return nil // Don't fail, just inform
	}
	
	fmt.Println("✓ Internet Sharing enabled!")
	return nil
}

func (hm *HotspotManager) disableHotspotMac() error {
	script := `tell application "System Events"
		tell network preferences
			set internetSharing to service "Internet Sharing" of sharing preferences
			set enabled of internetSharing to false
		end tell
	end tell`
	
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable Internet Sharing: %v", err)
	}
	
	return nil
}

func (hm *HotspotManager) getHotspotStatusMac() (bool, error) {
	// Check if Internet Sharing is enabled
	cmd := exec.Command("defaults", "read", "/Library/Preferences/SystemConfiguration/com.apple.nat", "NAT", "-dict", "Enabled")
	output, err := cmd.Output()
	if err != nil {
		return false, nil // Not enabled or error reading
	}
	
	return strings.Contains(string(output), "1"), nil
}

func (hm *HotspotManager) getIPAddressMac() (string, error) {
	// Get IP address of the bridge interface (usually bridge100)
	cmd := exec.Command("ifconfig", "bridge100")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}
	
	return "", fmt.Errorf("no IP address found")
}

// ===== Linux Implementation =====

func (hm *HotspotManager) isConnectedLinux() (bool, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "DEVICE,STATE", "device")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "connected") && strings.Contains(line, "wlan") {
			return true, nil
		}
	}
	
	return false, nil
}

func (hm *HotspotManager) enableHotspotLinux(ssid, password string) error {
	fmt.Println("📱 Enabling WiFi Hotspot on Linux...")
	
	// Use NetworkManager to create hotspot
	cmd := exec.Command("nmcli", "device", "wifi", "hotspot",
		"ssid", ssid,
		"password", password)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable hotspot: %v\nOutput: %s", err, output)
	}
	
	fmt.Printf("✓ WiFi hotspot enabled!\n")
	fmt.Printf("  SSID: %s\n", ssid)
	fmt.Printf("  Password: %s\n", password)
	
	return nil
}

func (hm *HotspotManager) disableHotspotLinux() error {
	// Stop the hotspot connection
	cmd := exec.Command("nmcli", "connection", "down", "Hotspot")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable hotspot: %v", err)
	}
	
	return nil
}

func (hm *HotspotManager) getHotspotStatusLinux() (bool, error) {
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,TYPE,DEVICE", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	return strings.Contains(string(output), "Hotspot"), nil
}

func (hm *HotspotManager) getIPAddressLinux() (string, error) {
	// Get IP of wlan0 or wlp interface
	cmd := exec.Command("ip", "addr", "show")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	inWlan := false
	
	for _, line := range lines {
		if strings.Contains(line, "wlan") || strings.Contains(line, "wlp") {
			inWlan = true
		} else if inWlan && strings.Contains(line, "inet ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Remove /24 or /16 suffix
				ip := strings.Split(parts[1], "/")[0]
				return ip, nil
			}
		}
	}
	
	return "", fmt.Errorf("no IP address found")
}

// ===== Windows Implementation =====

func (hm *HotspotManager) isConnectedWindows() (bool, error) {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	return strings.Contains(string(output), "State                  : connected"), nil
}

func (hm *HotspotManager) enableHotspotWindows(ssid, password string) error {
	fmt.Println("📱 Enabling WiFi Hotspot on Windows...")
	
	// Set up hosted network
	cmd := exec.Command("netsh", "wlan", "set", "hostednetwork",
		"mode=allow",
		fmt.Sprintf("ssid=%s", ssid),
		fmt.Sprintf("key=%s", password))
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to configure hotspot: %v\nOutput: %s", err, output)
	}
	
	// Start hosted network
	cmd = exec.Command("netsh", "wlan", "start", "hostednetwork")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start hotspot: %v\nOutput: %s", err, output)
	}
	
	fmt.Printf("✓ WiFi hotspot enabled!\n")
	fmt.Printf("  SSID: %s\n", ssid)
	fmt.Printf("  Password: %s\n", password)
	
	return nil
}

func (hm *HotspotManager) disableHotspotWindows() error {
	cmd := exec.Command("netsh", "wlan", "stop", "hostednetwork")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop hotspot: %v", err)
	}
	
	return nil
}

func (hm *HotspotManager) getHotspotStatusWindows() (bool, error) {
	cmd := exec.Command("netsh", "wlan", "show", "hostednetwork")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	return strings.Contains(string(output), "Status                 : Started"), nil
}

func (hm *HotspotManager) getIPAddressWindows() (string, error) {
	cmd := exec.Command("ipconfig")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IPv4 Address") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	
	return "", fmt.Errorf("no IP address found")
}