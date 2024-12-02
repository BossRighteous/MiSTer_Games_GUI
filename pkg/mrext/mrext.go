package mrext

import (
	"fmt"

	"github.com/wizzomafizzo/mrext/pkg/config"
	"github.com/wizzomafizzo/mrext/pkg/games"
	"github.com/wizzomafizzo/mrext/pkg/mister"
)

func GetSampleMgl() (string, error) {
	config := &config.UserConfig{}
	system := games.Systems["Nintendo64"]
	return mister.GenerateMgl(config, &system, fmt.Sprintf("/media/usb0/games/%s/game.n64", system.Folder), "")
}

/*
	write a utility to scan for gamelist.xml in rom directories mister.GamesFolders
	Bind to top level "Scan for ES gamelist.xml in games" to perform this
	Build a very basic mapping of SystemKeys to gamelistPaths string[]

	v1 keep it simple, no need for a grand indexing system beyond on-demand xml load.
*/
