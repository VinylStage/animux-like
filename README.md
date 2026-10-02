# Animux 🐧 🐱 🐶

**Animux** (Animux-like) is a terminal-based virtual pet (Tamagotchi-style) that lives natively in your Linux/Unix environment. 
Built with Unix philosophy in mind, it operates as a background daemon and lets you interact with your pet using standard CLI commands or a real-time TUI observer.

## Installation

You can install Animux quickly using standard package managers!

### Ubuntu / Debian (`apt`)
```bash
# Add the repository and install
echo "deb [trusted=yes] https://apt.fury.io/vinylstage/ /" | sudo tee /etc/apt/sources.list.d/animux.list
sudo apt update
sudo apt install animux
```

### Snap Store (`snap`)
```bash
sudo snap install animux
```

### macOS / Linux (`brew`)
```bash
brew tap VinylStage/animux-like https://github.com/VinylStage/animux-like
brew install animux
```

## Quick Start

1. **Start the background daemon:**
   ```bash
   animux daemon &
   ```

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
