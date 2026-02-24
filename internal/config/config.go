package config

import "github.com/primefaces/tyle/internal/layout"

type Config struct {
	Settings      Settings
	CustomLayouts []CustomLayout
}

type Settings struct {
	DelayBetweenSplitsMs int
	AutoEqualize         bool
	PickerColumns        int
	GhosttyConfigPath    string
}

type CustomLayout struct {
	ID          string
	Name        string
	Description string
	Preview     []string
	PaneCount   int
	Steps       []CustomLayoutStep
}

type CustomLayoutStep struct {
	Action    string
	Direction string
	DelayMs   int
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

func Load() Config {
	return DefaultConfig()
}

func (c Config) ToLayouts() []layout.Layout {
	return nil
}
