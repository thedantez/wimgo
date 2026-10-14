# WimGo

minimalist, keyboard driven web browser written in Go using WebKit & CGO

![Status](https://img.shields.io/badge/status-alpha-orange)
![Go](https://img.shields.io/badge/go-1.22+-blue)
![Platform](https://img.shields.io/badge/platform-linux-lightgrey)

# features
- **modal ui (vim-like)** : modes (`normal`, `insert`, `command`, `visual`(in progress))
- **minimalist statusbar** : flat monochrome command entry & laod progress indicator
- **embedded dark mode** : integrated Dark Reader engine by JS injection
- **fast navigation** : url normalization & quick command execution

# keybindings
| key | action |
| :---: | :---: |
| `j`/`k` | scroll down/up |
| `o` | open command prompt w/ `:open`
| `:` | open clean command prompt |
| `i` | switch to Insert mode |
| `v` | switch to Visual mode (now not work) |
| `esc` | return to Normal mode / close command prompt |
| `r` | reload site (like soft reload: use unchanged data from cache) |
| `R` | hard reload site (full reload site) |

# commands
| command | action |
| :---: | :---: |
| `:q`/`:quit` | close window |
| `:o`/`:open` | open url |
| `:r`/`:reload` | reload site (like like soft reload: use unchanged data from cache)
| `:r!`/`:reload!` | hard reload site (full reload site) |

# dependencies
make sure u've GTK3 & WebKitGTK installed:

**arch-based**
```bash
sudo pacman -S base-devel go gtk3 webkit2gtk-4.1 # or u can use yay
```

**debian-based**
```bash
sudo apt isntall build-essential golang libgtk-3-dev libwebkit2gtk-4.1-dev
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

[x] command prompt (`:open`, `:quit`)

[ ] smart search handling

[ ] configuration file support (config.?)

[ ] gRPC bridge & core daemon separation

[ ] tab & session management


