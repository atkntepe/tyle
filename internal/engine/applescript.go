package engine

import (
	"fmt"
	"os/exec"
	"strings"
)

type KeyCombo struct {
	Key       string
	Modifiers []string // "command", "shift", "control", "option"
}

func SendKeystroke(combo KeyCombo) error {
	mods := make([]string, len(combo.Modifiers))
	for i, m := range combo.Modifiers {
		mods[i] = m + " down"
	}
	modStr := strings.Join(mods, ", ")

	script := fmt.Sprintf(
		`tell application "System Events" to tell process "Ghostty" to keystroke "%s" using {%s}`,
		combo.Key, modStr,
	)

	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func EnsureGhosttyFocused() error {
	cmd := exec.Command("osascript", "-e",
		`tell application "Ghostty" to activate`)
	return cmd.Run()
}

func CheckAccessibilityPermission() bool {
	cmd := exec.Command("osascript", "-e",
		`tell application "System Events" to get name of first process`)
	return cmd.Run() == nil
}

func IsGhosttyRunning() bool {
	cmd := exec.Command("osascript", "-e",
		`tell application "System Events" to (name of processes) contains "Ghostty"`)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}
