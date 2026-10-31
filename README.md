# WimGo

minimalist, keyboard driven web browser written in Go using WebKit & CGO

![Status](https://img.shields.io/badge/status-alpha-orange)
![Go](https://img.shields.io/badge/go-1.22+-blue)
![Platform](https://img.shields.io/badge/platform-linux-lightgrey)

# features
- **modal ui** : modes (`normal`, `insert`, `command`, `visual`(in progress))
- **minimalist statusbar** : flat monochrome command entry & laod progress indicator
- **native dark mode** : via GTK & WebKit settings
- **fast work** : quick command execution & native work in linux via GTK
- **dynamic configuration** : line-by-line evaluated config file (`~/.config/wimgo/config`)
- **vim-like navigation** : shows hints to emulate mouse clicks
- **customizable keybindings** : rebind/disable shorcuts on the fly via `:bind` & `:unbind`

# commands
| command | action |
| :---: | :---: |
| `:q`/`:quit` | close window |
| `:o`/`:open` | open url |
| `:r`/`:reload` | reload site (like like soft reload: use unchanged data from cache)
| `:r!`/`:reload!` | hard reload site (full reload site) |
| `:b`/`:back` | go back history |
| `:f`/`:forward` | go forward history |
| `:set <key> <val>` | change configuration setting |
| `:config-save` | save current configuration to `~/.config/wimgo/config` |
| `:config-load` | load configuration from `~/.config/wimgo/config` |
| `:stop` | stop loading page |
| `:bind <key> <cmd>` | bind `<key>` as action `<cmd>` |
| `:unbind <key>` | unbind `<key>` |
| `:bind-default`/`:reset-default` | reset to default factory mapping |

# configuration
wimgo uses a plain text command-driven configuration file located at `~/.config/wimgo/config` (creating in first start)
```wimgo
# example wimgo configuration file
set homePage https://google.com
set searchEngine https://google.com/search?q=
set scrollStep 100
bind j scroll_down
bind k scroll_up
```
| variable | value | example (& default values) |
| :---: | :---: | :---: |
| `homePage` | string | `https://google.com` |
| `searchEngine` | string | `https://google.com/search?q=` |
| `scrollStep` | integer | 100 |
| `chordTimeout` | integer | 1200 |
| `webgl` | bool | true |
| `bind <key> <value>`| string, string | `h history_back`

All variables (w/out bindings) can be used via:
- lowerCamelCase (like `homePage`)
- snake_case     (like `home_page`)
- lowerflatcase  (like `homepage`)

 default actions
| key    | cmd           | action |
| :---:  | :---:         | :---: |
| `j` | `scroll_down` | scroll down |
| `k` | `scroll_up`   | scroll up |
| `gg` | `scroll_top`  | go to page's top |
| `G` | `scroll_bottom` | go to page's bottom |
| `H` | `history_back` | go back/forward history |
| `L` | `history_forward` | go back/forward history |
| `o` | - | open command prompt w/ `:open` |
| `:` | - | open clean command prompt |
| `i` | - | switch to Insert mode |
| `v` | - | switch to Visual mode (now not work) |
| `esc` | - | return to Normal mode / close command prompt |
| `r` | `soft_reload` | reload site (like soft reload: use unchanged data from cache) |
| `R` | `hard_reload` | hard reload site (full reload site) |
| `-` | `zoom_out` | zoom out |
| `+` | `zoom_in` | zoom in |
| `=` | `reset_zoom` | reset zoom |
| `f` | `hints` | show hints |

# fast installing (only in Arch-based)
```bash
git clone https://github.com/thedantez/wimgo.git
cd wimgo/app
makepkg
```

# regular installing (building from source)
## dependencies
make sure u've GTK3 & WebKitGTK installed:

**arch-based**
```bash
sudo pacman -S base-devel go gtk3 webkit2gtk-4.1 # or u can use yay

# optional: GStreamer
sudo pacman -S gst-plugins-good gst-plugins-bad gst-plugins-ugly gst-libav
```

**debian-based**
```bash
sudo apt install build-essential golang libgtk-3-dev libwebkit2gtk-4.1-dev

# optional: GStreamer
sudo apt install gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav
```

## building & running
```bash
git clone https://github.com/thedantez/wimgo.git
cd wimgo

# for build binary file for ur system
go build .
./wimgo

# for fast run w/out building
go run .
```

# flags
| flag | mean |
| :---: | :---: |
| `--no-cookies` | run w/out cookies |

# roadmap
[x] WebKitGTK window w/ CGO bridge

[x] FSM keyboard navigation

[x] command prompt (`:open`, `:quit`, `:set`, `:bind`)

[x] smart search handling

[x] configuration file support

[ ] gRPC bridge & core daemon separation

[ ] tab & session management
