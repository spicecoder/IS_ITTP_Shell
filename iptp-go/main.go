package main

import (
	"fmt"
	"os"
)

func main() {
	// Initialize state
	stateFile := getStateFilePath()
	state, err := LoadState(stateFile)
	if err != nil {
		// Create new state if doesn't exist
		state = NewState(stateFile)
	}

	// Check if running as a command (e.g., iptp goto /path)
	if len(os.Args) > 1 {
		// Command mode
		exitCode := ExecuteCommand(state, os.Args[1:])
		os.Exit(exitCode)
	}

	// Interactive REPL mode
	fmt.Println("🚀 iptp- IPTP Shell Process Manager")
	fmt.Println("   Type 'help' for commands, 'exit' to quit")
	fmt.Println()

	shell := NewShell(state)
	shell.Run()
}

// getStateFilePath returns the path to the state file
// Works on Windows, macOS, and Linux
func getStateFilePath() string {
	if os.Getenv("GOOS") == "windows" {
		return os.Getenv("TEMP") + "\\iptp_state.json"
	}
	return "/tmp/iptp_state.json"
}
