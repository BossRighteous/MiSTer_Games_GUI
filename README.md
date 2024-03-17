### Notice
This is a development WIP! All planned features are 'proven' in isolation but full support is in progress. \
Currently Supported Features listed below

# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer Python script GUI for your game library.

Runs on MiSTer, using Groovy Mister core for headless graphics processing.

Will include utilities for data processing from scrapper API(s). All image data will need transcoded to uncompressed BGR8 bytes

# WIP install directions
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
- Installable and runnable via sh Script
- Utilize existing Python3 build available in MiSTer Main
- Require no external dependencies via PIP or C Bindings
- Balance usable interface concerns and responsiveness with single threaded low CPU/Memory demand
- Offload otherwise expensive computation to off-MiSTer data parsing utilities and DAT files

## MiSTer Requirements
The GUI requires the (currently dev) GroovyMiSTer core for meaningful operation

## Currently Supported Features
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

## Roadmap
- Implement directory navigation
- Mock data seeder for meta / image
- Undecided meta format, XML/JSON/DAT
- Arcade15 theme state machine and 'screens' containing meta and image data
- ScreenScraper.fr data util for library parsing