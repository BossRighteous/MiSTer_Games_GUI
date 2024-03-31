package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func main() {
	start := time.Now()
	imagePath := os.Args[1]
	fmt.Println(imagePath)

	src, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	src = imaging.Fit(src, 256, 240, imaging.Lanczos)
	bounds := src.Bounds()
	var pix []uint8
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			pix = append(pix, uint8(b), uint8(g), uint8(r))
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(pix)
	fmt.Println(elapsed)
}
