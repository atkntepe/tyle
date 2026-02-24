package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/primefaces/tyle/internal/config"
	"github.com/primefaces/tyle/internal/engine"
	"github.com/primefaces/tyle/internal/layout"
	"github.com/primefaces/tyle/internal/tui"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:          "tyle",
		Short:        "Layout manager for Ghostty terminal",
		Version:      version,
		RunE:         runTUI,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(applyCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(resetCmd())
	rootCmd.AddCommand(initCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	layouts := layout.Presets()
	layouts = append(layouts, cfg.ToLayouts()...)

	model := tui.NewModel(layouts)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	m := finalModel.(tui.Model)

	if m.Cancelled() || m.Selected() == nil {
		return nil
	}

	ghosttyConfig := cfg.Settings.GhosttyConfigPath
	if ghosttyConfig == "" {
		ghosttyConfig = engine.GhosttyConfigPath()
	}
	bindings, err := engine.ParseGhosttyKeybindings(ghosttyConfig)
	if err != nil {
		return fmt.Errorf("failed to read Ghostty config: %w", err)
	}

	fmt.Printf("Applying layout: %s...\n", m.Selected().Name)
	time.Sleep(200 * time.Millisecond)

	if err := engine.ExecuteLayout(*m.Selected(), bindings, cfg.Settings.DelayBetweenSplitsMs); err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

func applyCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "apply [layout-id]",
		Short: "Apply a layout directly without the picker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()
			layouts := layout.Presets()
			layouts = append(layouts, cfg.ToLayouts()...)

			var target *layout.Layout
			for i, l := range layouts {
				if l.ID == args[0] {
					target = &layouts[i]
					break
				}
			}
			if target == nil {
				return fmt.Errorf("layout '%s' not found — run 'tyle list' to see available layouts", args[0])
			}

			if dryRun {
				fmt.Printf("Layout: %s (%d panes)\n\n", target.Name, target.PaneCount)
				fmt.Println("Steps:")
				for i, step := range target.Steps {
					switch step.Action {
					case layout.ActionSplit:
						fmt.Printf("  %d. Split %s\n", i+1, step.Direction)
					case layout.ActionFocus:
						fmt.Printf("  %d. Focus %s\n", i+1, step.Direction)
					case layout.ActionEqualize:
						fmt.Printf("  %d. Equalize splits\n", i+1)
					case layout.ActionDelay:
						fmt.Printf("  %d. Delay %dms\n", i+1, step.DelayMs)
					}
				}
				return nil
			}

			bindings, _ := engine.ParseGhosttyKeybindings(engine.GhosttyConfigPath())
			return engine.ExecuteLayout(*target, bindings, cfg.Settings.DelayBetweenSplitsMs)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print the keystroke sequence without executing")
	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available layouts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			layouts := layout.Presets()
			layouts = append(layouts, cfg.ToLayouts()...)

			for _, l := range layouts {
				fmt.Printf("  %-20s %s (%d panes)\n", l.ID, l.Name, l.PaneCount)
			}
		},
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Print the Ghostty keybind to add to your config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Add this line to your Ghostty config:")
			fmt.Printf("  %s\n\n", engine.GhosttyConfigPath())
			fmt.Println(`  keybind = cmd+shift+l=text:tyle\x0d`)
			fmt.Println()
			fmt.Println("This binds Cmd+Shift+L to launch the layout picker.")
			fmt.Println()
			fmt.Println("Optional: directional focus bindings (for custom layouts):")
			fmt.Println()
			fmt.Println("  keybind = cmd+alt+left=goto_split:left")
			fmt.Println("  keybind = cmd+alt+right=goto_split:right")
			fmt.Println("  keybind = cmd+alt+up=goto_split:top")
			fmt.Println("  keybind = cmd+alt+down=goto_split:bottom")
		},
	}
}

func resetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Close all splits in the current tab",
		RunE: func(cmd *cobra.Command, args []string) error {
			bindings, _ := engine.ParseGhosttyKeybindings(engine.GhosttyConfigPath())
			combo, ok := bindings["close_surface"]
			if !ok {
				return fmt.Errorf("no keybinding found for close_surface")
			}

			fmt.Println("Closing splits...")
			for i := 0; i < 10; i++ {
				if err := engine.SendKeystroke(combo); err != nil {
					break
				}
				time.Sleep(150 * time.Millisecond)
			}
			return nil
		},
	}
}
