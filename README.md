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

# keybindings
| key | action |
| :---: | :---: |
| `j`/`k` | scroll down/up |
| `gg` | go to page's top |
| `gG`/`G` | go to page's bottom |
| `H`/`L` | go back/forward history |
| `o` | open command prompt w/ `:open` |
| `:` | open clean command prompt |
| `i` | switch to Insert mode |
| `v` | switch to Visual mode (now not work) |
| `esc` | return to Normal mode / close command prompt |
| `r` | reload site (like soft reload: use unchanged data from cache) |
| `R` | hard reload site (full reload site) |
| `-`/`zo` | zoom out |
| `+`/`zi` | zoom in |
| `=`/`zz` | reset zoom |
| `f` | show hints |

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

# configuration
wimgo uses a plain text command-driven configuration file located at `~/.config/wimgo/config` (creating in first start)
```wimgo
# wimgo configuration file
set homePage https://google.com
set searchEngine https://google.com/search?q=
set scrollStep 100
set chordTimeout 1200
```

# dependencies
make sure u've GTK3 & WebKitGTK installed:

**arch-based**
```bash
sudo pacman -S base-devel go gtk3 webkit2gtk-4.1 # or u can use yay
```

**debian-based**
```bash
sudo apt install build-essential golang libgtk-3-dev libwebkit2gtk-4.1-dev
```

# building & running
```bash
git clone https://github.com/thedantez/wimgo.git
cd wimgo

# for build binary file for ur system
go build .
./wimgo

# for fast run
go run .
```

# roadmap
[x] WebKitGTK window w/ CGO bridge

[x] FSM keyboard navigation

[x] command prompt (`:open`, `:quit`, `:set`)

[x] smart search handling

[x] configuration file support

[ ] gRPC bridge & core daemon separation

[ ] tab & session management
