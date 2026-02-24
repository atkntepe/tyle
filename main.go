package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(hideCmd())
	rootCmd.AddCommand(showCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func allLayouts(cfg config.Config) []layout.Layout {
	layouts := layout.Presets()
	layouts = append(layouts, cfg.ToLayouts()...)
	return layouts
}

func visibleLayouts(cfg config.Config) []layout.Layout {
	all := allLayouts(cfg)
	var visible []layout.Layout
	for _, l := range all {
		if !cfg.IsHidden(l.ID) {
			visible = append(visible, l)
		}
	}
	return visible
}

func runTUI(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	layouts := visibleLayouts(cfg)
	if len(layouts) == 0 {
		return fmt.Errorf("no visible layouts — run 'tyle show' to unhide layouts")
	}

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
			layouts := allLayouts(cfg)

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
	var showAll bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available layouts",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			layouts := allLayouts(cfg)

			for _, l := range layouts {
				hidden := ""
				if cfg.IsHidden(l.ID) {
					if !showAll {
						continue
					}
					hidden = " (hidden)"
				}
				fmt.Printf("  %-20s %s (%d panes)%s\n", l.ID, l.Name, l.PaneCount, hidden)
			}
		},
	}

	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show hidden layouts too")
	return cmd
}

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Create a custom layout interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Layout name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name == "" {
				return fmt.Errorf("name cannot be empty")
			}

			fmt.Print("Columns: ")
			colStr, _ := reader.ReadString('\n')
			colStr = strings.TrimSpace(colStr)
			numCols, err := strconv.Atoi(colStr)
			if err != nil || numCols < 1 || numCols > 6 {
				return fmt.Errorf("columns must be a number between 1 and 6")
			}

			rowsPerCol := make([]int, numCols)
			for i := 0; i < numCols; i++ {
				fmt.Printf("Column %d rows: ", i+1)
				rowStr, _ := reader.ReadString('\n')
				rowStr = strings.TrimSpace(rowStr)
				rows, err := strconv.Atoi(rowStr)
				if err != nil || rows < 1 || rows > 6 {
					return fmt.Errorf("rows must be a number between 1 and 6")
				}
				rowsPerCol[i] = rows
			}

			l := layout.GenerateLayout(name, rowsPerCol)

			fmt.Println()
			for _, line := range l.Preview {
				fmt.Printf("  %s\n", line)
			}
			fmt.Printf("  %d panes\n\n", l.PaneCount)

			cfg := config.Load()
			cfg.AddLayout(config.FromLayout(l))
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Saved \"%s\" to %s\n", l.ID, config.ConfigPath())
			return nil
		},
	}
}

func hideCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hide [layout-id]",
		Short: "Hide a layout from the picker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()

			found := false
			for _, l := range allLayouts(cfg) {
				if l.ID == args[0] {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("layout '%s' not found — run 'tyle list --all' to see all layouts", args[0])
			}

			cfg.HideLayout(args[0])
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Hidden \"%s\" from the picker\n", args[0])
			return nil
		},
	}
}

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [layout-id]",
		Short: "Unhide a layout in the picker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load()

			cfg.ShowLayout(args[0])
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Showing \"%s\" in the picker\n", args[0])
			return nil
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
