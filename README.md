# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer Python script GUI for your game library

Runs on MiSTer, using Groovy Mister core for headless graphics processing.

Utilities to allow parsing of EmulationStation gamelist.xml and image data in simple binary DAT files for fast consumption.

# No Stable Build Available

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
- Blit Loop @6hz
- Rectangle Bitmap region composite and buffer manipulation

## Blocking Issues
- No input listening is currently available in GroovyMiSTer or `/dev/input/` and must be simulated

## Roadmap
- Sprite Sheet blitting
- Font DAT generation and blitting
- gamelist.xml parsing to DAT file
- CLI utility for remote input simulation