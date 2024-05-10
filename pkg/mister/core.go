package mister

import (
	"os"
	"path/filepath"
)

type CoreMGL struct {
	RBF   string
	File  string
	Delay string
	Type  string
	Index string
	Path  string
}

type Core struct {
	Name    string
	DirName string
	Slug    string
	MGL     CoreMGL
}

// This should get loaded from DB/JSON for extensibility
func GetCoresFromJSON() []Core {
	return []Core{
		{Name: "Nintendo 64", DirName: "N64", Slug: "n64"},
	}
}

func GetCorePathByPriority(core Core, fsRoot string) (string, bool) {
	coreDirs := []string{
		filepath.Join(fsRoot, "media", "usb0", core.DirName),
		filepath.Join(fsRoot, "media", "usb1", core.DirName),
		filepath.Join(fsRoot, "media", "usb2", core.DirName),
		filepath.Join(fsRoot, "media", "usb3", core.DirName),
		filepath.Join(fsRoot, "media", "usb4", core.DirName),
		filepath.Join(fsRoot, "media", "usb5", core.DirName),
		filepath.Join(fsRoot, "media", "usb0", "games", core.DirName),
		filepath.Join(fsRoot, "media", "usb1", "games", core.DirName),
		filepath.Join(fsRoot, "media", "usb2", "games", core.DirName),
		filepath.Join(fsRoot, "media", "usb3", "games", core.DirName),
		filepath.Join(fsRoot, "media", "usb4", "games", core.DirName),
		filepath.Join(fsRoot, "media", "usb5", "games", core.DirName),
		filepath.Join(fsRoot, "media", "network", core.DirName),
		filepath.Join(fsRoot, "media", "network", "games", core.DirName),
		filepath.Join(fsRoot, "media", "fat", "cifs", core.DirName),
		filepath.Join(fsRoot, "media", "fat", "cifs", "games", core.DirName),
		filepath.Join(fsRoot, "media", "fat", core.DirName),
		filepath.Join(fsRoot, "media", "fat", "games", core.DirName),
	}

	for _, path := range coreDirs {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path, true
		}
	}
	return "", false
}
