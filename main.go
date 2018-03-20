package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/aquilax/go-perlin"
)

var (
	// Image properties
	width, height int    = 5120, 2880
	filename      string = "islands.png"

	// Colornames
	water    color.RGBA = color.RGBA{85, 172, 238, 255}
	sand     color.RGBA = color.RGBA{194, 177, 128, 255}
	grass    color.RGBA = color.RGBA{99, 229, 33, 255}
	forest   color.RGBA = color.RGBA{49, 180, 58, 255}
	mountain color.RGBA = color.RGBA{230, 230, 230, 255}

	// Generation settings
	zoom float64 = 4000
	seed int     = 300

	// Abundance of each material
	numSnow  uint8 = 85
	numForr  uint8 = 50
	numGrass uint8 = 20
	numSand  uint8 = 15
)

func saveImage(fname string, img *image.RGBA) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	png.Encode(f, img)
}

func input() string {
	return ask("")
}

func ask(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func config() {
	sizeStr := ask("Enter Window Size (WIDTHxHEIGHT) ex. 300x450.\nOr type '5k', '4k', '2k'(2560x1440), or '1k'(1920x1080)\n")
	switch sizeStr {
	case "5k":
		width = 5120
		height = 2880
		break
	case "4k":
		width = 3480
		height = 2160
		break
	case "2k":
		width = 2560
		height = 1440
		break
	case "1k":
		width = 1920
		height = 1080
		break
	default:
		split := strings.Split(sizeStr, "x")
		var err error
		width, err = strconv.Atoi(split[0])
		height, err = strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}
		break
	}
	zoomStr := ask("Enter zoom level (1-5000): ")
	zoomInt, err := strconv.Atoi(zoomStr)
	if err != nil {
		panic(err)
	}
	zoom = float64(zoomInt)
}

func main() {
	config()
	fmt.Printf("Size: (%v, %v)\n", width, height)
	p := perlin.NewPerlin(2, 2, 3, 300)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			noiseVal := p.Noise2D(float64(x)/zoom, float64(y)/zoom)

			col := uint8((noiseVal + 1) * 127.5)
			var myCol color.RGBA
			if col > (255 - numSnow) {
				myCol = mountain
			} else if col > (255 - numSnow - numForr) {
				myCol = forest
			} else if col > (255 - numSnow - numForr - numGrass) {
				myCol = grass
			} else if col > (255 - numSnow - numForr - numGrass - numSand) {
				myCol = sand
			} else {
				myCol = water
			}
			// img.Set(x, y, color.RGBA{col, col, 0, 255})
			img.Set(x, y, myCol)
		}
	}
	saveImage(filename, img)
}
