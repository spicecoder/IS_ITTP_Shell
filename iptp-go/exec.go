package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// ExecuteScript runs an external command or script in gobash context
func ExecuteScript(parts []string) {
	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	args := parts[1:]

	// Create command
	cmd := exec.Command(cmdName, args...)

	// Set up environment - pass GOBASH context
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOBASH_SHELL=1")
	if cwd, err := os.Getwd(); err == nil {
		cmd.Env = append(cmd.Env, "GOBASH_CWD="+cwd)
	}

	// Connect to stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Command ran but returned non-zero exit code
			if runtime.GOOS != "windows" {
				// On Unix, we can get the exit code
				fmt.Printf("Command exited with code %d\n", exitErr.ExitCode())
			}
		} else {
			// Command couldn't be executed
			fmt.Printf("✗ Error executing command: %v\n", err)
		}
	}
}

// IsScriptOrCommand checks if a command is likely a script or external command
func IsScriptOrCommand(cmd string) bool {
	// Check if it's a file that exists
	if _, err := os.Stat(cmd); err == nil {
		return true
	}

	// Check if it has path separators (./script.sh, ../bin/tool)
	if containsPathSep(cmd) {
		return true
	}

	// Check if it's in PATH
	if _, err := exec.LookPath(cmd); err == nil {
		return true
	}

	return false
}

func containsPathSep(path string) bool {
	if runtime.GOOS == "windows" {
		return os.PathSeparator == '\\' && (path[0] == '.' || path[0] == '\\')
	}
	return path[0] == '.' || path[0] == '/'
}
