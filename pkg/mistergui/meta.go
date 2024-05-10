package mistergui

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
)

type Meta struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	ReleaseDate string `json:"releasedate"`
	Developer   string `json:"developer"`
	Publisher   string `json:"publisher"`
	Genre       string `json:"genre"`
	Players     string `json:"players"`
	Source      string `json:"source"`
	Id          string `json:"id"`
}

func LoadMetaImages(path string, rect image.Rectangle) []image.NRGBA {
	meta := loadMetaFile(path)
	return parseMetaImages(meta, rect)
}

func loadMetaFile(path string) Meta {
	var meta Meta
	dat, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("path not readable")
	}

	err = json.Unmarshal(dat, &meta)
	if err != nil {
		fmt.Println("Meta JSON unparsable")
	}
	if meta.Name == "" {
		meta.Name = filepath.Base(path)
	}

	return meta
}

func parseMetaImages(meta Meta, rect image.Rectangle) []image.NRGBA {
	images := make([]image.NRGBA, 0)
	// These are gueses for now
	charsPerLine := int(math.Floor(float64(rect.Dx()) / 7.0))
	linesPerImage := int(math.Floor(float64(rect.Dy()) / 20.0))

	textLines := make([]string, 1)
	textLines[0] = fmt.Sprintf("Name: %s", meta.Name)

	if meta.Description != "" {
		textLines = append(textLines, "Description:")
		desc := meta.Description
		offset := 0
		end := len(meta.Description) - 1
		for offset < end {
			lookahead := charsPerLine
			if lookahead+offset >= end {
				lookahead = end - offset
			}
			fmt.Println("LoadMetaImages description looping", offset, lookahead, end)
			subslice := desc[offset : offset+lookahead]
			if lookahead == charsPerLine {
				for lookahead > 0 {
					if string(subslice[lookahead-1]) == " " {
						break
					}
					lookahead--
				}
			}
			fmt.Println("LoadMetaImages description looping", offset, lookahead, end)
			textLines = append(textLines, desc[offset:offset+lookahead])
			offset += lookahead
		}
	}
	fmt.Println("LoadMetaImages description looped")

	if meta.ReleaseDate != "" {
		textLines = append(textLines, fmt.Sprintf("Release Date: %s", meta.ReleaseDate))
	}
	if meta.Developer != "" {
		textLines = append(textLines, fmt.Sprintf("Developer: %s", meta.Developer))
	}
	if meta.Publisher != "" {
		textLines = append(textLines, fmt.Sprintf("Publisher: %s", meta.Publisher))
	}
	if meta.Genre != "" {
		textLines = append(textLines, fmt.Sprintf("Genre: %s", meta.Genre))
	}
	if meta.Genre != "" {
		textLines = append(textLines, fmt.Sprintf("Players: %s", meta.Genre))
	}
	if meta.Players != "" {
		textLines = append(textLines, fmt.Sprintf("Players: %s", meta.Players))
	}

	lineOffset := 0
	lineEnd := len(textLines) - 1
	for lineOffset < lineEnd {
		if lineOffset+linesPerImage >= lineEnd {
			linesPerImage = lineEnd - lineOffset
		}
		images = append(images, *DrawText(textLines[lineOffset:lineOffset+linesPerImage], rect, image.Transparent))
		lineOffset += linesPerImage
	}
	return images
}
