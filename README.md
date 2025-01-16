# MiSTer Games GUI (MiSTer FPGA)

Low-resolution analog friendly rich-media GUI for your MiSTer FPGA game library.

"They" said it was impossible? At least Sorg doesn't have to approve a core :|

[Imgur Preview GIF](https://imgur.com/a/mister-gui-wip-2025-01-05-1xgubYu)

## Public Beta

After many hours of work, rework; failure and success I've reached an operable POC!

Instead of try to perfect everything I see wrong, now is the time to get it out and collect feedback and bug reports.

Discord channel for discussion and bug reports TBD

If you are worried about SDCard writes, please avoid downloading all collections. They are subject to change during Beta.

**Groovy Core Note**

I can help try to debug Groovy Core issues, but the author [psakhis](https://github.com/psakhis) unfortunately passed away late last year. It is an amazing core but please be understanding there is no active developer that has taken on a primary fork. There are also operational unknowns. IT'S USUALLY A NETWORK ISSUE, which shouldn't be a problem on loopback for this case.



## What it does:

- **Fast, responsive GUI for navigating your games and loading them directly from the GUI**
- Provides [full media libraries](https://github.com/BossRighteous/MiSTer_Games_Data_Utils) for common MGL Loadable Systems - [mrext](https://github.com/wizzomafizzo/mrext) & [Zaparoo](https://github.com/ZaparooProject/zaparoo-core)
- No need to scrape your library on the internet!
- Games as primary browsing unit, ROMs are children
- Low-resolution, 15khz analog display friendly
- Minimal input button requirments (NES / JAMMA interfaces)
- Executes as a Script from OSD
- Uses the wonderful [GroovyMister Core](https://github.com/psakhis/Groovy_MiSTer) as an interface for Video and Input handling

## What it does not do:

- **Replace the built in OSD Browser experience!!**
- **Seriously, this is not meant for primary on-boot loading!**
- Load the cores/script quickly. There is a timing based handshake.
- Guarantee handshake with the core and script 100% of the time. It's a brittle timing process!
- Cover the entire MiSTer core library (Limited to MGL and RetroarchDB overlap)
- Animation or video
- Bridge and index across systems. System DBs and indexes are standalone

## Known Feature Roadmap

- ROM Index reports for debugging - NDJSON
- Improvement of UI layout, overscan, navigation
- Image scaling/centering. Description scrolling
- Ability to hold directional input for repeat fire.
- Modifier key for D-pad -> List start/end/next-alpha
- Ability to navigate via PS2/keyboard
- Arcade MGDB curation in other repo
- Fuzzy /usb[0-5]/ path resolver for collections_path
- Input idle timeout with slideshow/random routine or self-exit
- Integrated collection manager (download, update, remove)
- Integrated script self-update
- Improved detection and handling of Groovy core handshake

## What it may be capable of doing later

- Being run via a separate daemon with a hotkey-combo listener to return to Groovy/GUI on demand
- Trim DB on demand to fit local index, remove unused DB entries
- Search by title
- Filtering by Genre, Developer, Publisher
- Supporting Custom multi-system Collections (All franchise games, favorites, etc)
- Matching by hash or ROM data - currently requires RetroarchDB NoIntro/Tosec/Redump naming for match


# Beta Installation Directions

Replace {} with your inputs

**Assumes GroovyMiSTer core is installed and tested with bouncing logo**

The script features some painfully long delays that will be tweaked, but can be changed via INI

To be executed via `ssh root@misterIP`
```
ssh root@{misterIP}
# Password: 1
cd /media/fat/Scripts
mkdir mistergamesgui
cd mistergamesgui
wget "https://github.com/BossRighteous/MiSTer_Games_GUI/releases/download/beta/load-mistergamesgui.sh"
wget "https://github.com/BossRighteous/MiSTer_Games_GUI/releases/download/beta/mistergamesgui"
wget "https://github.com/BossRighteous/MiSTer_Games_GUI/releases/download/beta/mistergamesgui.ini"
mkdir collections
cd collections

# Download Collections
wget "{Links from https://github.com/BossRighteous/MiSTer_Games_Database_MGDB/releases/tag/latest}"
# Example https://github.com/BossRighteous/MiSTer_Games_Database_MGDB/releases/download/latest/SNES.Console.1990-11-21.mgdb

# Open the Groovy core via OSD
# Should see bouncing logo, press primary controller button to activate
# Run the script via shell
cd /media/fat/Scripts/mistergamesgui && ./mistergamesgui

# OR reboot and from OSD run Scripts > mistergamesgui > load-mistergamesgui.sh
```

**GroovyMiSTer Install directions**

To be executed via `ssh root@misterIP`
```
ssh root@{misterIP}
# Password: 1
cd /media/fat/
wget "https://github.com/psakhis/Groovy_MiSTer/releases/download/0.7/MiSTer_groovy"
cd _Utility
wget "https://github.com/psakhis/Groovy_MiSTer/releases/download/0.7/Groovy_20240922.rbf"
cd ..

# This must be done for any active .ini!
nano {MiSTer}.ini

# ADD LINES TO BOTTOM OF INI
[Groovy]
main=MiSTer_groovy
```

## Thank you
- [psakhis](https://github.com/psakhis) for your amazing Core and collaborating on the GMC Command runner. I wish you could see this
- [wizzomafizzo](https://wizzo.dev/) for the mrext project which made programming the loaders and building for target a breeze
- My wife for putting up with my hundreds of hours huddled away on the PC, love you
- [#nogpu GroovyArcade discord](https://discord.com/channels/649595547308785664/1030412595884204082) for all the help over the last year
- [lllllllllllllllllll](https://forums.somethingawful.com/showthread.php?threadid=4058840#post542644664) For being the only fuckin' person to respond to my devlog thread on the SA Forums (lol. lmao even.)