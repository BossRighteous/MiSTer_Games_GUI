package mgdb

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BossRighteous/MiSTer_Games_GUI/pkg/utils"
)

func SlugifyString(input string) string {
	r := regexp.MustCompile(`(\(.*\))|(\[.*\])|(\.\w*$)|[^a-z0-9A-Z]`)
	rep := r.ReplaceAllStringFunc(input, func(m string) string {
		return ""
	})
	return strings.ToLower(rep)
}

func MakeIndexedRomFromPath(path string, gameID int) IndexedRom {
	rom := IndexedRom{
		Path:               path,
		FileExt:            filepath.Ext(path),
		GameID:             gameID,
		SupportedSystemIds: "",
	}
	file, _ := utils.CutSuffix(filepath.Base(path), rom.FileExt)
	rom.FileName = file
	return rom
}
