package main

import (
	"fmt"
	"os"
)

// ExecuteCommand executes iptp in command mode (non-interactive)
// e.g., iptp goto /path
func ExecuteCommand(state *State, args []string) int {
	if len(args) == 0 {
		return 0
	}

	cmd := args[0]
	cmdArgs := args[1:]

	// Get current process name (from environment or default)
	currentProcess := os.Getenv("iptp_PROCESS")
	if currentProcess == "" {
		currentProcess = fmt.Sprintf("shell_%d", os.Getpid())
	}

	switch cmd {
	case "goto":
		return cmdGotoNonInteractive(state, currentProcess, cmdArgs)
	case "name":
		return cmdNameNonInteractive(state, currentProcess, cmdArgs)
	case "save":
		return cmdSaveNonInteractive(state, currentProcess)
	case "list":
		return cmdListNonInteractive(state)
	case "jump":
		return cmdJumpNonInteractive(state, cmdArgs)
	case "state":
		return cmdStateNonInteractive(state, currentProcess)
	case "help", "--help", "-h":
		cmdHelpNonInteractive()
		return 0
	case "version", "--version", "-v":
		fmt.Println("iptp version 1.0.0")
		fmt.Println("IPTP Shell Process Manager")
		return 0
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Run 'iptp help' for usage")
		return 1
	}
}

func cmdGotoNonInteractive(state *State, process string, args []string) int {
	if len(args) == 0 {
		fmt.Println("Usage: iptp goto PATH")
		return 1
	}

	path := args[0]
	oldDir, _ := os.Getwd()

	if err := os.Chdir(path); err != nil {
		fmt.Printf("✗ Cannot change directory: %v\n", err)
		return 1
	}

	newDir, _ := os.Getwd()
	state.UpdateDirectory(process, newDir, oldDir)
	state.Save()

	fmt.Printf("✓ Changed to: %s\n", newDir)
	return 0
}

func cmdNameNonInteractive(state *State, process string, args []string) int {
	if len(args) == 0 {
		if proc, ok := state.GetProcess(process); ok {
			fmt.Printf("Current process: %s\n", process)
			fmt.Printf("Intention: %s\n", proc.Intention)
		} else {
			fmt.Printf("Current process: %s\n", process)
		}
		return 0
	}

	intention := joinArgs(args)
	processName := ParseIntention(intention)

	currentDir, _ := os.Getwd()
	state.SetProcess(processName, intention, currentDir)
	state.Save()

	fmt.Printf("✓ Shell named: %s\n", processName)
	fmt.Printf("  Intention: %s\n", intention)
	return 0
}

func cmdSaveNonInteractive(state *State, process string) int {
	currentDir, _ := os.Getwd()
	state.UpdateDirectory(process, currentDir, "")
	state.Save()

	fmt.Printf("✓ Saved state for: %s @ %s\n", process, currentDir)
	return 0
}

func cmdListNonInteractive(state *State) int {
	processes := state.ListProcesses()
	if len(processes) == 0 {
		fmt.Println("No saved processes")
		return 0
	}

	fmt.Println("=== Available Processes ===")
	for _, name := range processes {
		if proc, ok := state.GetProcess(name); ok {
			fmt.Printf("  → %s: %s (PID: %d)\n", name, proc.CurrentDir, proc.PID)
		}
	}
	return 0
}

func cmdJumpNonInteractive(state *State, args []string) int {
	if len(args) == 0 {
		fmt.Println("Usage: iptp jump PROCESS")
		return 1
	}

	targetProcess := args[0]
	proc, ok := state.GetProcess(targetProcess)
	if !ok {
		fmt.Printf("✗ Process '%s' not found\n", targetProcess)
		return 1
	}

	fmt.Printf("Process: %s\n", targetProcess)
	fmt.Printf("Directory: %s\n", proc.CurrentDir)
	fmt.Printf("Intention: %s\n", proc.Intention)

	return 0
}

func cmdStateNonInteractive(state *State, process string) int {
	proc, ok := state.GetProcess(process)
	if !ok {
		fmt.Println("No state for current process")
		return 1
	}

	currentDir, _ := os.Getwd()

	fmt.Println("=== Current Process State (IPTP Format) ===")
	fmt.Printf("Process: %s\n", process)
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

	return 0
}

func cmdHelpNonInteractive() {
	fmt.Println("iptp - IPTP Shell Process Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  iptp              - Start interactive shell")
	fmt.Println("  iptp COMMAND      - Execute single command")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  name [INTENTION]    - Name current process")
	fmt.Println("  goto PATH           - Navigate to directory")
	fmt.Println("  save                - Save current state")
	fmt.Println("  list                - List all processes")
	fmt.Println("  jump PROCESS        - Show process info")
	fmt.Println("  state               - Show current state")
	fmt.Println("  help                - Show this help")
	fmt.Println("  version             - Show version")
	fmt.Println()
	fmt.Println("Interactive Mode:")
	fmt.Println("  Run 'iptp' with no arguments to enter interactive shell")
	fmt.Println("  Inside iptp shell, all commands work plus:")
	fmt.Println("    - getmethere: Interactive directory finder")
	fmt.Println("    - back: Navigate back in history")
	fmt.Println("    - Execute any script or command")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  iptp                                 # Start shell")
	fmt.Println("  iptp goto /var/www                   # Change directory")
	fmt.Println("  iptp name \"working on authentication\"  # Name process")
	fmt.Println("  iptp list                            # List processes")
	fmt.Println()
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}
