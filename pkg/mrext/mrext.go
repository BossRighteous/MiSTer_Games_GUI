package mrext

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/mgdb"
	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/utils"
	"github.com/wizzomafizzo/mrext/pkg/config"
	"github.com/wizzomafizzo/mrext/pkg/games"
	"github.com/wizzomafizzo/mrext/pkg/mister"
)

var cfg *config.UserConfig = &config.UserConfig{}

func GetSampleMgl() (string, error) {
	system := games.Systems["Nintendo64"]
	return mister.GenerateMgl(cfg, &system, fmt.Sprintf("/media/usb0/games/%s/game.n64", system.Folder), "")
}

/*
	write a utility to scan for gamelist.xml in rom directories mister.GamesFolders
	Bind to top level "Scan for ES gamelist.xml in games" to perform this
	Build a very basic mapping of SystemKeys to gamelistPaths string[]

	v1 keep it simple, no need for a grand indexing system beyond on-demand xml load.
*/

func GetSystemsByIDsString(idsStr string) []games.System {
	ids := strings.Split(idsStr, ",")
	systems := make([]games.System, 0)
	for _, id := range ids {
		fmt.Println(id)
		system, err := games.GetSystem(id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(*system)
		systems = append(systems, *system)
	}
	return systems
}

func GetGamesFolders() []string {
	// TODO: HEY IDIOT DON'T FORGET THIS
	isDev := true
	if isDev {
		return []string{
			"/mnt/c/Users/bossr/Code/MiSTer_Games_Data_Utils/cores",
			"/mnt/c/Users/bossr/Code/MiSTer_Games_GUI/games",
		}
	}
	return config.GamesFolders
}

func GetSystemsPaths(systems []games.System) []games.PathResult {
	var matches []games.PathResult
	gamesFolders := GetGamesFolders()
	fmt.Println(gamesFolders)
	for _, system := range systems {
		for _, gamesFolder := range gamesFolders {
			fmt.Println(gamesFolder)
			gf, err := games.FindFile(gamesFolder)
			if err != nil {
				continue
			}

			for _, folder := range system.Folder {
				systemFolder := filepath.Join(gf, folder)
				path, err := games.FindFile(systemFolder)
				if err != nil {
					continue
				}

				matches = append(matches, games.PathResult{System: system, Path: path})
			}
		}
	}
	return matches
}

func GetSystemsGamesPaths(systems []games.System) ([]string, error) {
	gameFiles := make([]string, 0)
	pathResults := GetSystemsPaths(systems)
	for _, pathResult := range pathResults {
		fmt.Println(pathResult)
		system := pathResult.System
		path := pathResult.Path
		pathGames, err := games.GetFiles(system.Id, path)
		if err != nil {
			return nil, err
		}
		gameFiles = append(gameFiles, pathGames...)
	}
	return gameFiles, nil
}

func GetRelativeGamePath(absPath string) (string, bool) {
	// Attempt to trim GetGamesFolders paths from abs for storage
	for _, rootPath := range GetGamesFolders() {
		relPath, wasCut := utils.CutPrefix(absPath, rootPath)
		if wasCut {
			return relPath, true
		}
	}
	return absPath, false
}

func GetFirstGamePathFromRelative(relPath string) (string, bool) {
	// usb pathing can/will change, use relative pathing for on-demand loading
	for _, gFolder := range GetGamesFolders() {
		absPath := filepath.Join(gFolder, relPath)
		if games.FileExists(absPath) {
			return absPath, true
		}
	}
	return relPath, false
}

func ScanMGDBGames(client mgdb.MGDBClient) (chan string, chan bool) {
	outBuffer := make(chan string, 1024)
	completedOk := make(chan bool)

	go func() {
		bPrint := func(msg string) {
			fmt.Println(msg)
			outBuffer <- msg
		}

		info, err := client.GetMGDBInfo()
		if err != nil {
			bPrint(err.Error())
			completedOk <- false
			return
		}

		// TODO: reset IsIndexed on Games table to 0

		systems := GetSystemsByIDsString(info.SupportedSystemIds)
		roms, err := GetSystemsGamesPaths(systems)
		if err != nil {
			bPrint(err.Error())
			completedOk <- false
			return
		}

		bPrint(fmt.Sprintf("Found %v ROMs", len(roms)))
		indexedCount := 0
		for _, romAbsPath := range roms {
			fmt.Println(romAbsPath)
			romRelPath, ok := GetRelativeGamePath(romAbsPath)
			if !ok {
				bPrint("Relative Pathing Error on Scanned ROM")
				bPrint(romRelPath)
				continue
			}
			bPrint(fmt.Sprintf("Found %v", romRelPath))

			// Match to DB
			// Update Game row IsIndexed
			// Add IndexedRom record
			indexedCount++
			bPrint(fmt.Sprintf("Indexed %v", romRelPath))

		}
		bPrint(fmt.Sprintf("Indexed %v ROMs", len(roms)))

		completedOk <- true
	}()

	return outBuffer, completedOk
}
