# Hyprtime Architecture

This document describes the technical architecture of Hyprtime, a modular screen time tracking system for Hyprland.

## Overview

Hyprtime is split into two independent processes that communicate via REST API:

1. **Daemon (hyprtimed)**: Background service that tracks window usage and owns the database
2. **GUI (hyprtime)**: Desktop application that displays statistics via API calls

## Design Goals

- **Separation of Concerns**: Daemon and GUI have distinct, well-defined responsibilities
- **Modularity**: Easy to add new frontends (CLI, web UI, etc.) without changing daemon
- **Maintainability**: Clear code organization by domain
- **Testability**: Components can be tested independently
- **Performance**: Unix socket for fast local IPC

## Component Details

### Daemon (hyprtimed)

**Location:** `cmd/hyprtimed/`, `internal/daemon/`

**Responsibilities:**
1. **Window Tracking** (`internal/daemon/tracker/`)
   - Subscribes to Hyprland IPC events
   - Tracks active window focus changes
   - Records window open/close events
   - Manages current session state

2. **Database Management** (`internal/daemon/database/`)
   - Owns the SQLite database exclusively
   - Handles all write operations
   - Provides read operations for API server
   - Schema: `apps` table (global stats) + `daily_stats` table (per-day stats)

3. **IPC Communication** (`internal/daemon/ipc/`)
   - Connects to Hyprland via Unix sockets
   - Parses window information (class, title, address)
   - Handles event stream from `.socket2.sock`

4. **REST API Server** (`internal/daemon/api/`)
   - HTTP server over Unix socket
   - Exposes tracking data to clients
   - JSON responses with formatted data
   - Endpoints: health check, today's stats, daily stats

**Key Files:**
- `cmd/hyprtimed/main.go` - Entry point, lifecycle management
- `internal/daemon/tracker/tracker.go` - Core tracking logic
- `internal/daemon/api/server.go` - HTTP server setup
- `internal/daemon/api/handlers.go` - API endpoint handlers
- `internal/daemon/database/write.go` - Database mutations
- `internal/daemon/database/query.go` - Database queries

### GUI (hyprtime)

**Location:** `main.go`, `internal/gui/`, `frontend/`

**Responsibilities:**
1. **API Client** (`internal/gui/client/`)
   - HTTP client for daemon API
   - Unix socket communication
   - Context-aware requests
   - Error handling and timeouts

2. **Wails Service** (`internal/gui/service/`)
   - Exposes Go methods to frontend
   - Wraps API client calls
   - Manages socket path resolution
   - Interface between Go backend and Svelte frontend

3. **Frontend** (`frontend/src/`)
   - Svelte components for UI
   - Date selection and navigation
   - Application usage visualization
   - Wails runtime integration

**Key Files:**
- `main.go` - GUI entry point, Wails setup
- `internal/gui/client/client.go` - HTTP API client
- `internal/gui/service/service.go` - Wails service layer
- `frontend/src/App.svelte` - Root component
- `frontend/src/views/DailyView.svelte` - Main view
- `frontend/src/components/AppList.svelte` - Usage list

### Shared Components

**Location:** `internal/shared/`, `internal/logger/`, `internal/utils/`

**Shared Models** (`internal/shared/models/`)
- `AppData` - Application usage data structure
- `DailyData` - Daily statistics container
- Used by both daemon (API responses) and GUI (API client)
- Auto-generated TypeScript bindings by Wails

**Logger** (`internal/logger/`)
- Centralized logging for daemon
- Log levels: Quiet, Normal, Verbose, Debug
- Structured output with prefixes

**Utils** (`internal/utils/`)
- Helper functions used across the project

## Communication Flow

### Data Write Flow (Tracking)

```
Hyprland Window Event
        ↓
  Hyprland IPC (.socket2.sock)
        ↓
  internal/daemon/ipc/ipc.go
        ↓
  internal/daemon/tracker/tracker.go
        ↓
  internal/daemon/database/write.go
        ↓
  SQLite Database
```

### Data Read Flow (GUI)

```
User Interaction (GUI)
        ↓
  frontend/src/views/DailyView.svelte
        ↓
  ScreenTimeService.GetDailyStats() (Wails binding)
        ↓
  internal/gui/service/service.go
        ↓
  internal/gui/client/client.go (HTTP request)
        ↓
  Unix Socket ($XDG_RUNTIME_DIR/hyprtime/daemon.sock)
        ↓
  internal/daemon/api/server.go
        ↓
  internal/daemon/api/handlers.go
        ↓
  internal/daemon/database/query.go
        ↓
  SQLite Database
        ↓
  JSON Response → API Client → GUI
```

## API Design

### Transport: Unix Socket

**Path:** `$XDG_RUNTIME_DIR/hyprtime/daemon.sock` (typically `/run/user/{uid}/hyprtime/daemon.sock`)

**Protocol:** HTTP over Unix Socket

**Why Unix Socket?**
- Fast: ~1μs overhead vs direct function call
- Secure: Filesystem permissions (0600)
- No network exposure: Only local access
- Standard: Works with existing HTTP clients

### Endpoints

See [api/openapi.yaml](api/openapi.yaml) for full API specification.

**Base URL:** `http://unix` (when using Unix socket transport)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/api/v1/stats/today` | Today's stats |
| GET | `/api/v1/stats/daily/{date}` | Stats for specific date |

### Response Format

All successful responses return JSON:

```json
{
  "date": "2026-03-21",
  "total_time": 3600,
  "total_time_formatted": "1h 0m 0s",
  "apps": [
    {
      "class": "kitty",
      "total_time": 1800,
      "total_time_formatted": "30m 0s",
      "open_count": 5,
      "last_seen": "2026-03-21",
      "percentage": 50.0
    }
  ]
}
```

## Database Schema

**Location:** `~/.local/share/hyprtime/hyprtime.db`

**Owner:** Daemon only (GUI has no direct access)

### Tables

**apps** - Global application statistics
```sql
CREATE TABLE apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class TEXT UNIQUE NOT NULL,
    total_time INTEGER DEFAULT 0,
    open_count INTEGER DEFAULT 0,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

**daily_stats** - Per-day statistics
```sql
CREATE TABLE daily_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    date DATE NOT NULL,
    total_time INTEGER DEFAULT 0,
    open_count INTEGER DEFAULT 0,
    FOREIGN KEY (app_id) REFERENCES apps(id),
    UNIQUE(app_id, date)
)
```

## Build System

The project uses [Task](https://taskfile.dev/) for build automation.

### Build Process

```
Build GUI:
  1. Install frontend deps (npm install)
  2. Generate Wails bindings (from Go models)
  3. Build frontend (Svelte → dist/)
  4. Embed frontend in Go binary
  5. Build Go binary with Wails

Build Daemon:
  1. Build Go binary from cmd/hyprtimed/
  (No frontend dependencies)
```

### Independent Builds

- Daemon can be built without frontend dependencies
- GUI requires frontend to be built first
- Both share Go models from `internal/shared/models/`

## Testing Strategy

### Unit Tests

- API handlers: Test with mock database
- API client: Test with mock HTTP server
- Database functions: Test with in-memory SQLite

### Integration Tests

- Daemon API: Start daemon, test endpoints
- GUI API client: Test against real daemon
- End-to-end: Full workflow from tracking to display

### Manual Testing

```bash
# 1. Start daemon
./bin/hyprtimed --verbose

# 2. Verify socket
ls -la $XDG_RUNTIME_DIR/hyprtime/daemon.sock

# 3. Test API
echo -e "GET /api/v1/stats/today HTTP/1.1\r\nHost: unix\r\n\r\n" | nc -U $XDG_RUNTIME_DIR/hyprtime/daemon.sock

# 4. Start GUI
./bin/hyprtime
```

## Extension Points

The modular architecture makes it easy to extend:

### Adding New API Endpoints

1. Add handler in `internal/daemon/api/handlers.go`
2. Register route in `internal/daemon/api/server.go`
3. Add client method in `internal/gui/client/client.go`
4. Update `api/openapi.yaml`

### Adding CLI Tool

Create a new client that uses `internal/gui/client/` to call the daemon API.

### Adding Web UI

Create a web frontend that calls the daemon API (might need to add HTTP server option).

### Multi-User Support

- Currently single-user (per-user database and socket)
- Could extend to support multiple users by isolating databases
- API already supports this via socket permissions

## Performance Considerations

**API Overhead:**
- Unix socket: ~1-2ms per request
- JSON encoding/decoding: ~1ms
- Total overhead: < 10ms (acceptable for UI)

**Database:**
- SQLite in WAL mode for concurrent reads
- Daemon owns writes (no conflicts)
- Periodic updates (1-minute intervals) reduce write frequency

**Memory:**
- Daemon: ~10 MB RSS
- GUI: ~50 MB RSS (includes Wails + Svelte)

## Security

- Unix socket with 0600 permissions (user-only access)
- No network exposure
- No authentication needed (local-only, filesystem-protected)
- Database at `~/.local/share/` (user-only)

## Future Enhancements

1. **CLI Tool** - Command-line interface for querying stats
2. **Web Dashboard** - Browser-based UI
3. **Metrics Export** - Prometheus exporter
4. **Cloud Sync** - Optional cloud backup/sync
5. **Aggregations** - Weekly/monthly reports
6. **Categories** - Group apps by category
7. **Notifications** - Time limit alerts
8. **Multi-Monitor** - Track per-monitor usage

## Migration from Direct DB Access

Prior to the refactoring, the GUI accessed the SQLite database directly. The new architecture:

**Before:**
```
GUI → SQLite Database ← Daemon
(both reading and writing)
```

**After:**
```
GUI → API (Unix socket) → Daemon → SQLite Database
(single owner)
```

**Benefits:**
- No read/write race conditions
- Single source of truth
- Clean API boundary
- Easy to add new clients

## References

- **Wails Documentation:** https://wails.io/docs/introduction
- **Hyprland IPC:** https://wiki.hyprland.org/IPC/
- **Task:** https://taskfile.dev/
- **Svelte:** https://svelte.dev/docs
