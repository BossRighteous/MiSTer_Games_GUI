### Notice
This is currently under active development. This is Babby's First Go project so feel free to recommend changes. Channel structures are still a bit wonky I think.

# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer script GUI for your game library.

Aiming for something like EmulationStation with filesystem browsing paired to Meta data and images.

Will run on MiSTer ARM chip, using Groovy Mister core for headless graphics processing.

Will include utilities for data processing from scrapper API(s)? 

# WIP local dev direction
- clone repo
- cd {repo_path}
- go run cmd/mistergamesgui/main.go

# WIP Install direction (unstable at the moment)
- mister ssh
- cd /media/fat/Scripts
- wget https://github.com/BossRighteous/MiSTer_Games_GUI/blob/main/_bin/linux_arm/mistergamesgui
- chmod 755
- (Open the Groovy Core on the MiSTer, need to pin down my old version or update :/)
- ./mistergamesgui

## Goals
- Provide a simple analog friendly graphical interface for browsing your library
- Installable via DB json ini
- Go executable build for Cyclone V chip
- Require no external dependencies
- Balance usable interface concerns and responsiveness with single threaded low CPU/Memory demand

## MiSTer Requirements
The GUI requires the (currently dev) GroovyMiSTer core for meaningful operation.
Developing against build Groovy_20240327.rbf

## Go Supported Features
- UDP connection to GroovyMiSTer via localhost loopback
- INI settings for modeline etc
- GroovyMister API basic implementation
- Image/TTF-font embeds and rendering w/ transparency
- FPS display
- Basic goroutine/channel support for lazy-loads against blit cycle
- Loading meta JSON from disk into overlay Images
- Input support for 2 9 Button digital controllers

## Go Roadmap
- Directory navigation
- MGL temp writes for Core/Game loading
- Display meta-data for selected game
- Display image(s) for selected game
- Interlace support
- Alternate or adaptive GUI for 480p vs 240p
- Scaper to Meta routine
- Meta dumps offloaded to static hosting to avoid scraping needs