# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer Python script GUI for your game library

Runs on MiSTer, using Groovy Mister core for headless graphics processing.

Utilities to allow parsing of EmulationStation gamelist.xml and image data in simple binary DAT files for fast consumption.

# Semi-stable install directions, needs love
- ssh to MiSTer
- `cd /media/fat/Scripts/`
- `mkdir MiSTer_Games_GUI`
- (On other computer)
- clone repo and rsync src/ to `/media/fat/Scripts/
- (Back on MiSTer ssh shell)
- `cd /media/fat/Scripts/MiSTer_Games_GUI/`
- `python3 mister_game_gui.py`

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
- Blit Loop @60hz
- Rectangle Bitmap region composite and buffer manipulation
- Image loading from data_utils DAT creation
- Bitmap Font loading from data_utils DAT creation
- FontSheet with trimmable text box rendering
- SpriteSheet cell based blitting minimally working
- Awkward implementation of CLI args, will rework
- "Themes" abstract with hooks for rendering
- Idle frame queues for cross-frame work or deferred blit sequencing

## Blocking Issues
- No input listening is currently available in GroovyMiSTer or `/dev/input/` and must be simulated

## Roadmap
- Reworking args and app level boot/settings
- CLI utility for remote input simulation
- Arcade15 theme support for state machine and 'screens'
- Will expiriment with direct XML parsing and PNG native loading
- Othewise utils for gamelist.xml parsing and images to DAT files