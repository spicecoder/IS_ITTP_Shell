package main

import (
	"encoding/json"
	"os"
	"time"
)

// Pulse represents a trivalent truth condition (IPTP)
type Pulse struct {
	Name     string `json:"name"`
	TV       string `json:"TV"` // "Y", "N", or "U"
	Response string `json:"response"`
}

// Process represents a named shell process with state
type Process struct {
	Intention  string   `json:"intention"`
	CurrentDir string   `json:"current_dir"`
	History    []string `json:"history"`
	PID        int      `json:"pid"`
	Timestamp  string   `json:"timestamp"`
	Pulses     []Pulse  `json:"pulses"`
}

// State represents the global gobash state
type State struct {
	Processes map[string]Process `json:"processes"`
	filepath  string
}

// NewState creates a new empty state
func NewState(filepath string) *State {
	return &State{
		Processes: make(map[string]Process),
		filepath:  filepath,
	}
}

// LoadState loads state from JSON file
func LoadState(filepath string) (*State, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	state.filepath = filepath
	return &state, nil
}

// Save writes state to JSON file
func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0644)
}

// SetProcess creates or updates a process
func (s *State) SetProcess(name, intention, currentDir string) {
	process := Process{
		Intention:  intention,
		CurrentDir: currentDir,
		History:    []string{},
		PID:        os.Getpid(),
		Timestamp:  time.Now().Format(time.RFC3339),
		Pulses: []Pulse{
			{Name: "process named", TV: "Y", Response: name},
			{Name: "directory saved", TV: "Y", Response: currentDir},
		},
	}

	// Preserve history if process already exists
	if existing, ok := s.Processes[name]; ok {
		process.History = existing.History
	}

	s.Processes[name] = process
}

// UpdateDirectory updates the current directory for a process
func (s *State) UpdateDirectory(processName, newDir, oldDir string) {
	process, ok := s.Processes[processName]
	if !ok {
		// Create new process if doesn't exist
		s.SetProcess(processName, "Working in "+processName, newDir)
		return
	}

	// Add old directory to history
	if oldDir != "" && oldDir != newDir {
		process.History = append(process.History, oldDir)
	}

	process.CurrentDir = newDir
	process.Timestamp = time.Now().Format(time.RFC3339)
	process.Pulses = []Pulse{
		{Name: "process named", TV: "Y", Response: processName},
		{Name: "directory saved", TV: "Y", Response: newDir},
	}

	s.Processes[processName] = process
}

// GetProcess retrieves a process by name
func (s *State) GetProcess(name string) (Process, bool) {
	process, ok := s.Processes[name]
	return process, ok
}

// ListProcesses returns all process names
func (s *State) ListProcesses() []string {
	names := make([]string, 0, len(s.Processes))
	for name := range s.Processes {
		names = append(names, name)
	}
	return names
}

// PopHistory removes and returns the last history entry for a process
func (s *State) PopHistory(processName string) (string, bool) {
	process, ok := s.Processes[processName]
	if !ok || len(process.History) == 0 {
		return "", false
	}

	lastIdx := len(process.History) - 1
	lastDir := process.History[lastIdx]
	process.History = process.History[:lastIdx]

	s.Processes[processName] = process
	return lastDir, true
}
