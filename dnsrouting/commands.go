package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteCommand runs a command with arguments (used for non-interactive mode)
func ExecuteCommand(state *State, args []string) int {
	if len(args) == 0 {
		return 0
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "goto":
		return execGoto(state, cmdArgs)
	case "version":
		fmt.Println("iptp version 1.0.0 (IPTP Shell)")
		return 0
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Available commands: goto, version")
		return 1
	}
}

// execGoto handles goto command in non-interactive mode
func execGoto(state *State, args []string) int {
	if len(args) == 0 {
		fmt.Println("Usage: iptp goto PATH")
		return 1
	}

	path := args[0]
	if err := os.Chdir(path); err != nil {
		fmt.Printf("Error: Cannot change to %s: %v\n", path, err)
		return 1
	}

	newDir, _ := os.Getwd()
	processName := fmt.Sprintf("shell_%d", os.Getpid())
	state.UpdateDirectory(processName, newDir, "")
	state.Save()

	fmt.Printf("Changed to: %s\n", newDir)
	return 0
}

// ExecuteScript runs an external command or script
func ExecuteScript(parts []string) {
	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	cmdArgs := parts[1:]

	// Determine shell based on OS
	var shellCmd string
	var shellArgs []string

	switch runtime.GOOS {
	case "windows":
		shellCmd = "cmd"
		shellArgs = []string{"/C", cmdName}
	default:
		shellCmd = "sh"
		shellArgs = []string{"-c", cmdName}
	}

	shellArgs = append(shellArgs, cmdArgs...)

	cmd := exec.Command(shellCmd, shellArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Set environment variables for the script
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "IPTP_SHELL=true")
	cmd.Env = append(cmd.Env, fmt.Sprintf("IPTP_PID=%d", os.Getpid()))

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Printf("Error running command: %v\n", err)
	}
}
