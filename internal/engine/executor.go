package engine

import (
	"fmt"
	"time"

	"github.com/primefaces/tyle/internal/layout"
)

func ExecuteLayout(l layout.Layout, bindings map[string]KeyCombo, delayMs int) error {
	if !IsGhosttyRunning() {
		return fmt.Errorf("ghostty is not running")
	}

	if err := EnsureGhosttyFocused(); err != nil {
		return fmt.Errorf("failed to focus ghostty: %w", err)
	}

	time.Sleep(100 * time.Millisecond)

	delay := time.Duration(delayMs) * time.Millisecond

	for _, step := range l.Steps {
		switch step.Action {
		case layout.ActionSplit:
			action := fmt.Sprintf("new_split:%s", step.Direction)
			combo, ok := bindings[action]
			if !ok {
				return fmt.Errorf("no keybinding found for %s — add it to your Ghostty config", action)
			}
			if err := SendKeystroke(combo); err != nil {
				return fmt.Errorf("failed to execute %s: %w", action, err)
			}

		case layout.ActionFocus:
			action := fmt.Sprintf("goto_split:%s", step.Direction)
			combo, ok := bindings[action]
			if !ok {
				return fmt.Errorf("no keybinding found for %s — add it to your Ghostty config", action)
			}
			if err := SendKeystroke(combo); err != nil {
				return fmt.Errorf("failed to execute %s: %w", action, err)
			}

		case layout.ActionEqualize:
			combo, ok := bindings["equalize_splits"]
			if !ok {
				continue
			}
			if err := SendKeystroke(combo); err != nil {
				return fmt.Errorf("failed to equalize: %w", err)
			}

		case layout.ActionDelay:
			time.Sleep(time.Duration(step.DelayMs) * time.Millisecond)
			continue
		}

		time.Sleep(delay)
	}

	return nil
}
