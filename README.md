### Notice
Currently migrating to Go for the rest of dev

After some initial work in python, I found the runtime too slow to be usable.

`src` is python, will be deleted when Go is working better
`pkg/cmd` are go paths

# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer script GUI for your game library.

Will run on MiSTer ARM chip, using Groovy Mister core for headless graphics processing.

Will include utilities for data processing from scrapper API(s)? 

# Go install direction
- build release TBD

# Old Python install directions
- ssh to MiSTer
- `cd /media/fat/Scripts/`
- `mkdir MiSTer_Games_GUI`
- (On other computer)
- clone repo
- cd `{REPO}/src`
- `rsync -r ./ root@{MISTER_HOST}:/media/fat/Scripts/MiSTer_Games_GUI`
- (Back on MiSTer ssh shell)
- `cd /media/fat/Scripts/MiSTer_Games_GUI/`
- `python3 mister_game_gui.py`

Pull the repo and rsync `src` as needed.

## Goals
- Provide a simple analog friendly graphical interface for browsing your library
- Installable via DB json ini
- Go executable build for Cyclone V chip
- Require no external dependencies
- Balance usable interface concerns and responsiveness with single threaded low CPU/Memory demand

## MiSTer Requirements
The GUI requires the (currently dev) GroovyMiSTer core for meaningful operation.

I may be a bit behind on versions too

## Go Supported Features
- UDP connection to GroovyMiSTer via localhost loopback
- Modeline parsing (hardcoded)
- GroovyMister API basic implementation
- Image/TTF-font embeds and rendering w/ transparency
- FPS display

## Go Roadmap
- Controller input handling support (wizzo mrext)
- Directory navigation
- MGL temp writes for Core/Game loading
- Display meta-data for selected game
- Display image(s) for selected game
- LZ4 compression for GroovyMister API
- Interlace support
- Configurable ini settings
- Alternate or adaptive GUI for 480p vs 240p
- Scaper to SQLite routine (maybe wizzo mrext)

## Old Python Supported Features
- Executable on MiSTer via SSH
- UDP connection to GroovyMiSTer via localhost loopback
- 256x240p@60hz default resolution rendering
- Blit Loop up to 60hz
- Rectangle Bitmap region composite and buffer manipulation
- Image loading from data_utils DAT creation
- Bitmap Font loading from data_utils DAT creation
- FontSheet with trimmable text box rendering
- SpriteSheet cell based blitting minimally working
- "Themes" abstract with hooks for rendering
- Idle frame queues for cross-frame work or deferred blit sequencing
- Joystick input support
- 24fps stable with unoptimized Joystick and fps output
