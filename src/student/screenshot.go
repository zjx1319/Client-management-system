package main

import (
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

func screenshot1() {
	img, _ := screenshot.CaptureDisplay(0)
	file, _ := os.Create("1.png")
	defer file.Close()
	png.Encode(file, img)
}
