package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/primefaces/tyle/internal/layout"
)

type Config struct {
	Settings      Settings       `toml:"settings"`
	CustomLayouts []CustomLayout `toml:"custom_layouts"`
}

type Settings struct {
	DelayBetweenSplitsMs int    `toml:"delay_between_splits_ms"`
	AutoEqualize         bool   `toml:"auto_equalize"`
	PickerColumns        int    `toml:"picker_columns"`
	GhosttyConfigPath    string `toml:"ghostty_config_path"`
}

type CustomLayout struct {
	ID          string             `toml:"id"`
	Name        string             `toml:"name"`
	Description string             `toml:"description"`
	Preview     []string           `toml:"preview"`
	PaneCount   int                `toml:"pane_count"`
	Steps       []CustomLayoutStep `toml:"steps"`
}

type CustomLayoutStep struct {
	Action    string `toml:"action"`
	Direction string `toml:"direction"`
	DelayMs   int    `toml:"delay_ms"`
}

func DefaultConfig() Config {
	return Config{
		Settings: Settings{
			DelayBetweenSplitsMs: 200,
			AutoEqualize:         true,
			PickerColumns:        3,
		},
	}
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "tyle", "config.toml")
}

func Load() Config {
	cfg := DefaultConfig()

	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return DefaultConfig()
	}

	return cfg
}

func (c Config) ToLayouts() []layout.Layout {
	var layouts []layout.Layout
	for _, cl := range c.CustomLayouts {
		var steps []layout.LayoutStep
		for _, s := range cl.Steps {
			steps = append(steps, layout.LayoutStep{
				Action:    layout.StepAction(s.Action),
				Direction: layout.Direction(s.Direction),
				DelayMs:   s.DelayMs,
			})
		}
		layouts = append(layouts, layout.Layout{
			ID:          cl.ID,
			Name:        cl.Name,
			Description: cl.Description,
			Preview:     cl.Preview,
			PaneCount:   cl.PaneCount,
			Steps:       steps,
		})
	}
	return layouts
}
