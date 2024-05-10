package mgdb

import (
	"fmt"
	"os"
	fsPath "path"
	"path/filepath"
	"strings"
)

type ExtendedFileInfo struct {
	Info     os.FileInfo
	Path     string
	FileName string
	Slug     string
}

func getFileNameFromFileInfo(info os.FileInfo) string {
	name := info.Name()
	return strings.TrimSuffix(fsPath.Base(name), fsPath.Ext(name))
}

func isExtAllowed(info os.FileInfo, filterExt []string) bool {
	ext := fsPath.Ext(info.Name())
	if len(filterExt) == 0 {
		return true
	}
	for _, v := range filterExt {
		if v == ext {
			return true
		}
	}
	return false
}

func GetDirContentsRecursive(rootPath string, filterExts []string) ([]ExtendedFileInfo, error) {
	files := make([]ExtendedFileInfo, 0)
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			fmt.Println("info nil")
			return nil
		}
		if info.IsDir() && info.Name() == ".mistergamesgui" {
			return filepath.SkipDir
		}
		//fmt.Println(path, info.Name())
		if !info.IsDir() && isExtAllowed(info, filterExts) {
			fName := getFileNameFromFileInfo(info)
			efInfo := ExtendedFileInfo{
				Info:     info,
				Path:     path,
				FileName: fName,
				Slug:     SlugifyString(fName),
			}
			files = append(files, efInfo)
			//fmt.Println(efInfo)
		}
		return nil
	})
	return files, err
}
