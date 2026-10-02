# Animux 🐧 🐱 🐶

**Animux** (Animux-like) is a terminal-based virtual pet (Tamagotchi-style) that lives natively in your Linux/Unix environment. 
Built with Unix philosophy in mind, it operates as a background daemon and lets you interact with your pet using standard CLI commands or a real-time TUI observer.

## Features

- **Multiple Species:** Adopt a Tux-like Penguin, a Cat, or a Dog.
- **Real-time Lifecycle:** Your pet gets hungry, bored, and makes messes over time. Don't neglect them!
- **XDG Compliant:** Uses standard Linux directories (`$XDG_STATE_HOME`, `$XDG_DATA_HOME`, `$XDG_RUNTIME_DIR`).
- **Interactive Observer:** Use `animux show` to watch your pet continuously in a colorful UI (powered by Bubble Tea & Lipgloss).
- **Daemon Architecture:** A lightweight background process manages time decay and handles CLI requests via Unix Domain Sockets.

## Installation

### Homebrew (macOS / Linux)
```bash
brew install VinylStage/tap/animux-like
```

### Debian/Ubuntu (.deb)
Download the latest `.deb` package from the [Releases](https://github.com/VinylStage/animux-like/releases) page and run:
```bash
sudo dpkg -i animux_*.deb
```

### Manual Build
```bash
git clone https://github.com/VinylStage/animux-like.git
cd animux-like
go build -o animux
sudo mv animux /usr/local/bin/
```

## Quick Start

1. **Start the background daemon:**
   ```bash
   animux daemon &
   ```
   *(Logs are continuously written to `~/.local/state/animux/animux.log`)*

2. **Adopt your first pet:**
   ```bash
   animux adopt penguin Pingu
   ```
   *(Supported species: `penguin`, `cat`, `dog`)*

3. **Interact with your pet!**
   ```bash
   animux status  # Check how your pet is doing
   animux feed    # Give them some food
   animux play    # Play a game
   animux clean   # Clean up their messes
   ```

4. **Watch them in real-time:**
   ```bash
   animux show
   ```
   *In the `show` view, you can press `f` (feed), `p` (play), or `c` (clean) to interact directly! Press `q` to exit.*

## Storage & Configuration

Animux strictly follows the XDG Base Directory Specification:
- **Save file:** `~/.local/share/animux/state.json`
- **Daemon Logs:** `~/.local/state/animux/animux.log`
- **Unix Socket:** `/run/user/$(id -u)/animux.sock` (or `/tmp/animux.sock`)
