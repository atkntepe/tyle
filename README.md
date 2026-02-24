# tyle

Layout manager for [Ghostty](https://ghostty.org) terminal. Pick a split layout from a TUI and apply it instantly.

<!-- screenshot or gif here -->

## Install

```bash
go install github.com/atkntepe/tyle@latest
```

Or download a binary from [Releases](https://github.com/atkntepe/tyle/releases).

## Setup

Add this to your Ghostty config (`~/.config/ghostty/config` or `~/Library/Application Support/com.mitchellh.ghostty/config`):

```
keybind = cmd+shift+l=text:tyle\x0d
```

This binds **Cmd+Shift+L** to launch the picker.

You also need to grant **Accessibility** permission to your terminal in System Settings > Privacy & Security > Accessibility.

## Usage

```bash
tyle              # open the layout picker
tyle list         # list available layouts
tyle apply <id>   # apply a layout directly
tyle add          # create a custom layout interactively
tyle hide <id>    # hide a layout from the picker
tyle show <id>    # unhide a layout
tyle reset        # close all splits
```

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
