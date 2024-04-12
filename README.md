### Notice
This is currently under active development. This is Babby's First Go project so feel free to recommend changes. Channel structures are still a bit wonky I think.

# MiSTer_Games_GUI
Low-resolution analog friendly MiSTer script GUI for your game library.

Will run on MiSTer ARM chip, using Groovy Mister core for headless graphics processing.

Will include utilities for data processing from scrapper API(s)? 

# Go install direction
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

I am definitely be a bit behind on versions too!

## Go Supported Features
- UDP connection to GroovyMiSTer via localhost loopback
- Modeline parsing (hardcoded)
- GroovyMister API basic implementation
- Image/TTF-font embeds and rendering w/ transparency
- FPS display
- Basic goroutine/channel support for lazy-loads against blit cycle
- Loading meta JSON from disk into overlay Images

## Go Roadmap
- Update for latest Groovy core packet support
- Controller input handling support (wizzo mrext or groovymister udp?)
- Directory navigation
- MGL temp writes for Core/Game loading
- Display meta-data for selected game
- Display image(s) for selected game
- LZ4 compression for GroovyMister API
- Interlace support
- Configurable ini settings
- Alternate or adaptive GUI for 480p vs 240p
- Scaper to SQLite routine (maybe wizzo mrext)