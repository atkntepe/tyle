package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/atkntepe/tyle/internal/layout"
)

type Config struct {
	Settings      Settings       `toml:"settings"`
	CustomLayouts []CustomLayout `toml:"custom_layouts"`
}

type Settings struct {
	DelayBetweenSplitsMs int      `toml:"delay_between_splits_ms"`
	AutoEqualize         bool     `toml:"auto_equalize"`
	PickerColumns        int      `toml:"picker_columns"`
	GhosttyConfigPath    string   `toml:"ghostty_config_path,omitempty"`
	HiddenLayouts        []string `toml:"hidden_layouts,omitempty"`
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
	Direction string `toml:"direction,omitempty"`
	DelayMs   int    `toml:"delay_ms,omitempty"`
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

func Save(cfg Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}

func FromLayout(l layout.Layout) CustomLayout {
	var steps []CustomLayoutStep
	for _, s := range l.Steps {
		steps = append(steps, CustomLayoutStep{
			Action:    string(s.Action),
			Direction: string(s.Direction),
			DelayMs:   s.DelayMs,
		})
	}
	return CustomLayout{
		ID:          l.ID,
		Name:        l.Name,
		Description: l.Description,
		Preview:     l.Preview,
		PaneCount:   l.PaneCount,
		Steps:       steps,
	}
}

func (c *Config) AddLayout(cl CustomLayout) {
	for i, existing := range c.CustomLayouts {
		if existing.ID == cl.ID {
			c.CustomLayouts[i] = cl
			return
		}
	}
	c.CustomLayouts = append(c.CustomLayouts, cl)
}

func (c *Config) HideLayout(id string) {
	for _, h := range c.Settings.HiddenLayouts {
		if h == id {
			return
		}
	}
	c.Settings.HiddenLayouts = append(c.Settings.HiddenLayouts, id)
}

func (c *Config) ShowLayout(id string) {
	var filtered []string
	for _, h := range c.Settings.HiddenLayouts {
		if h != id {
			filtered = append(filtered, h)
		}
	}
	c.Settings.HiddenLayouts = filtered
}

func (c Config) IsHidden(id string) bool {
	for _, h := range c.Settings.HiddenLayouts {
		if h == id {
			return true
		}
	}
	return false
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
