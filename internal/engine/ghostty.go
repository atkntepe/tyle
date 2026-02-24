package engine

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func DefaultKeybindings() map[string]KeyCombo {
	return map[string]KeyCombo{
		"new_split:right":     {Key: "d", Modifiers: []string{"command"}},
		"new_split:down":      {Key: "d", Modifiers: []string{"command", "shift"}},
		"goto_split:previous": {Key: "[", Modifiers: []string{"command"}},
		"goto_split:next":     {Key: "]", Modifiers: []string{"command"}},
		"equalize_splits":     {Key: "=", Modifiers: []string{"command", "shift"}},
		"close_surface":       {Key: "w", Modifiers: []string{"command", "shift"}},
	}
}

func GhosttyConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support",
		"com.mitchellh.ghostty", "config")
}

func ParseGhosttyKeybindings(configPath string) (map[string]KeyCombo, error) {
	bindings := DefaultKeybindings()

	file, err := os.Open(configPath)
	if err != nil {
		return bindings, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.HasPrefix(line, "keybind") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		value := strings.TrimSpace(parts[1])
		bindParts := strings.SplitN(value, "=", 2)
		if len(bindParts) != 2 {
			continue
		}

		trigger := strings.TrimSpace(bindParts[0])
		action := strings.TrimSpace(bindParts[1])

		combo := parseTrigger(trigger)
		if combo != nil {
			bindings[action] = *combo
		}
	}

	return bindings, scanner.Err()
}

func parseTrigger(trigger string) *KeyCombo {
	parts := strings.Split(trigger, "+")
	if len(parts) == 0 {
		return nil
	}

	combo := &KeyCombo{
		Key: parts[len(parts)-1],
	}

	modMap := map[string]string{
		"cmd":     "command",
		"command": "command",
		"shift":   "shift",
		"ctrl":    "control",
		"control": "control",
		"alt":     "option",
		"opt":     "option",
		"option":  "option",
	}

	for _, part := range parts[:len(parts)-1] {
		if mod, ok := modMap[strings.ToLower(part)]; ok {
			combo.Modifiers = append(combo.Modifiers, mod)
		}
	}

	return combo
}
