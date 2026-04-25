# Hyprtime

A minimal screen time tracker for hyprland.

## Features

- Track application and window usage in real-time
- View daily screen time statistics
- Persistent data storage with SQLite
- Clean and minimal UI

## Installation

```bash
yay -S hyprtime-bin
```

## Usage

1. Add the daemon to your Hyprland config (`~/.config/hypr/hyprland.conf`):

```
exec-once = hyprtimed
```

2. Restart Hyprland

3. Launch the application from your application menu or run:

```bash
hyprtime
```

The daemon will automatically track your screen time in the background, and you can open the application to view your daily usage statistics.

## License

See LICENSE file for details.
