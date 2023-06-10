# hyprwatch - hyprland event watcher daemon

## Config

`hyprwatch` expects config file to be present at `$HOME/.config/hyprwatch/config.yaml`

```yaml
# top level keys should be event names from https://wiki.hyprland.org/IPC/
monitoradded:
    - data: "DP-1" # all data is matched as string
      callback: "echo 'hello world'" # shell command to execute if data matches
```

