# Hyprtime

A screen time tracking tool for Hyprland window manager, providing detailed insights into application usage.

## Components

Hyprtime consists of two main components:

### 1. Daemon (`hyprtimed`)
A background service that:
- Monitors active Hyprland windows in real-time
- Tracks time spent in each application
- Records window open/close events
- Stores data in SQLite database
- Exposes REST API over Unix socket

### 2. GUI (`hyprtime`)
A desktop application that:
- Displays daily screen time statistics
- Shows per-application usage breakdown
- Provides date navigation for historical data
- Communicates with daemon via API

## Architecture

The project follows a modular architecture with clear separation of concerns:

```
┌─────────────────┐  Unix Socket    ┌──────────────────┐
│   hyprtimed     │   HTTP API      │    hyprtime      │
│   (daemon)      │◄────────────────┤    (Wails GUI)   │
│                 │                 │                  │
│  - Tracker      │                 │  - API Client    │
│  - Database     │                 │  - Frontend      │
│  - API Server   │                 │    (Svelte)      │
└─────────────────┘                 └──────────────────┘
         │
         ▼
    SQLite DB
```

**Key Design Principles:**
- **Single Source of Truth**: Daemon owns the database
- **API-First**: GUI communicates only via REST API
- **Unix Sockets**: Fast, secure local IPC
- **Modular**: Clear separation between daemon and GUI logic

## Directory Structure

```
hyprtime/
├── cmd/
│   └── hyprtimed/              # Daemon entry point
├── internal/
│   ├── daemon/                 # Daemon-specific code
│   │   ├── api/                # REST API server
│   │   ├── database/           # Database access layer
│   │   ├── tracker/            # Window tracking logic
│   │   └── ipc/                # Hyprland IPC communication
│   ├── gui/                    # GUI-specific code
│   │   ├── client/             # API client
│   │   └── service/            # Wails service layer
│   ├── shared/                 # Shared code
│   │   └── models              # Data models
│   ├── logger/                 # Logging utilities
│   └── utils/                  # Helper functions
├── frontend/                   # Svelte frontend
├── api/                        # API documentation
│   └── openapi.yaml
└── main.go                     # GUI entry point
```

## Build

### Prerequisites
- Go 1.21+
- Node.js 18+
- Task (taskfile)
- Wails v3 (for GUI development)

### Build Commands

```bash
# Build both daemon and GUI
task build

# Build daemon only
task build:daemon

# Build GUI only
task build:gui

# Development mode (with hot reload)
task dev
```

### Binaries

After building, you'll find:
- `bin/hyprtimed` - Daemon (7.5 MB)
- `bin/hyprtime` - GUI application (19 MB)

## Installation

1. Build the project:
   ```bash
   task build
   ```

2. Copy binaries to your PATH:
   ```bash
   sudo cp bin/hyprtimed /usr/local/bin/
   sudo cp bin/hyprtime /usr/local/bin/
   ```

3. Create database directory:
   ```bash
   mkdir -p ~/.local/share/hyprtime
   ```

## Usage

### Running the Daemon

The daemon must be running for the GUI to function:

```bash
# Run daemon with default logging
hyprtimed

# Verbose logging
hyprtimed --verbose

# Debug logging
hyprtimed --debug

# Quiet mode (errors only)
hyprtimed --quiet
```

The daemon:
- Stores data at: `~/.local/share/hyprtime/hyprtime.db`
- Creates API socket at: `$XDG_RUNTIME_DIR/hyprtime/daemon.sock`
- Saves tracking data every 1 minute

### Running the GUI

Once the daemon is running:

```bash
hyprtime
```

The GUI will:
- Connect to the daemon's Unix socket
- Display today's screen time by default
- Allow navigation to previous dates
- Show per-application usage statistics

### Autostart (Optional)

To start the daemon automatically, create a systemd user service:

```bash
mkdir -p ~/.config/systemd/user
```

Create `~/.config/systemd/user/hyprtimed.service`:

```ini
[Unit]
Description=Hyprtime Daemon
After=hyprland-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/hyprtimed
Restart=on-failure

[Install]
WantedBy=hyprland-session.target
```

Enable and start:

```bash
systemctl --user enable --now hyprtimed.service
```

## API

The daemon exposes a REST API over Unix socket. See [api/openapi.yaml](api/openapi.yaml) for full specification.

### Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/stats/today` - Today's statistics
- `GET /api/v1/stats/daily/{date}` - Stats for specific date (YYYY-MM-DD)

### Testing the API

You can test the API using netcat or curl with socat:

```bash
# Health check (using netcat)
echo -e "GET /api/v1/stats/today HTTP/1.1\r\nHost: localhost\r\n\r\n" | nc -U $XDG_RUNTIME_DIR/hyprtime/daemon.sock

# Or use curl with socat
alias unixcurl='socat - UNIX-CONNECT:$XDG_RUNTIME_DIR/hyprtime/daemon.sock'
echo -e "GET /api/v1/health HTTP/1.1\r\n\r\n" | unixcurl
```

## Development

### Project Setup

```bash
# Install dependencies
go mod download
cd frontend && npm install

# Generate Wails bindings
task common:generate:bindings
```

### Running in Development Mode

```bash
# Start dev server with hot reload
task dev
```

### Code Organization

The codebase is organized by domain:

- **`internal/daemon/`** - All daemon logic (tracking, database, API)
- **`internal/gui/`** - All GUI logic (API client, Wails service)
- **`internal/shared/`** - Code shared between daemon and GUI (models)
- **`internal/logger/`** and **`internal/utils/`** - Common utilities

This separation ensures:
- Clear ownership of code
- Easy testing of individual components
- GUI can be replaced without touching daemon
- Future CLI or web UI can reuse the same API

## Troubleshooting

### Daemon not starting

Check if Hyprland is running and the IPC socket exists:
```bash
ls -la /tmp/hypr/$HYPRLAND_INSTANCE_SIGNATURE/
```

### GUI shows "failed to connect to daemon"

1. Verify daemon is running:
   ```bash
   ps aux | grep hyprtimed
   ```

2. Check if socket exists:
   ```bash
   ls -la $XDG_RUNTIME_DIR/hyprtime/daemon.sock
   ```

3. Check daemon logs for errors

### Database locked errors

If you see database locked errors, ensure only one hyprtimed instance is running:
```bash
pkill hyprtimed
hyprtimed
```

## Contributing

Contributions are welcome! The modular architecture makes it easy to:

- Add new API endpoints (modify `internal/daemon/api/`)
- Add new frontend features (modify `frontend/src/`)
- Improve tracking logic (modify `internal/daemon/tracker/`)
- Add database migrations (modify `internal/daemon/database/`)

## License

[Add your license here]

## Acknowledgments

- Built with [Wails v3](https://wails.io/)
- Frontend framework: [Svelte](https://svelte.dev/)
- Window manager: [Hyprland](https://hyprland.org/)
