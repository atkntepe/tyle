# tyle

Layout manager for [Ghostty](https://ghostty.org) terminal. Pick a split layout from a TUI and apply it instantly.

<!-- screenshot or gif here -->

## Install

Download a binary from [Releases](https://github.com/atkntepe/tyle/releases).

Or install with Go:

```bash
brew install go
go install github.com/atkntepe/tyle@latest
```

## Setup

Add this to your Ghostty config (`~/.config/ghostty/config` or `~/Library/Application Support/com.mitchellh.ghostty/config`):

```
keybind = cmd+shift+l=text:tyle\x0d
```

This binds **Cmd+Shift+L** to launch the picker.

You also need to grant **Accessibility** permission to your terminal in System Settings > Privacy & Security > Accessibility.

## Usage

```bash
tyle                  # open the layout picker
tyle apply <id>       # apply a layout directly
tyle apply <id> --dry-run  # preview steps without executing
tyle list             # list available layouts
tyle list --all       # include hidden layouts
```

### Custom layouts

```bash
tyle add              # create a custom layout interactively
tyle hide <id>        # hide a layout from the picker
tyle show <id>        # unhide a layout
tyle reset            # close all splits
```

`tyle add` walks you through creating a layout by specifying the number of columns and rows per column.

## Build from source

```bash
git clone https://github.com/atkntepe/tyle.git
cd tyle
go build -o tyle .
./tyle
```

## Requirements

- macOS (uses AppleScript to send keystrokes)
- Ghostty terminal

## License

MIT
