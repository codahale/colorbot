// Command colorbot analyzes a given image and prints out a list of the N most
// dominant colors in the image.
//
//     $ colorbot test-assets/hodges-research.png
//     #000000
//     #af905a
//     #f6d185
//     #f9e2b3
//     #ffffff
//
//     $ curl https://dribbble.com/dribbble-logo.png | colorbot -
//     #000000
//     #621131
//     #c32361
//     #d83976
//     #ea4c89
//
// Colorbot supports GIF, JPEG, and PNG images.
package main

import (
	"flag"
	"fmt"
	_ "image/gif"  // support GIF images
	_ "image/jpeg" // support JPEG images
	_ "image/png"  // support PNG images
	"os"

	"github.com/codahale/colorbot"
)

func main() {
	var (
		n         = flag.Int("n", 5, "number of colors to return")
		maxBytes  = flag.Int64("maxbytes", 10*1024*1024, "max image size in bytes")
		maxPixels = flag.Int64("maxpixels", 10*1024*1024, "max image size in pixels")
	)
	flag.Parse()

	var in *os.File
	if s := flag.Args()[0]; s == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(s)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		in = f
	}

	img, err := colorbot.DecodeImage(in, *maxBytes, *maxPixels)
	if err != nil {
		panic(err)
	}

	colors := colorbot.DominantColors(img, *n)
	for _, color := range colors {
		r, g, b, _ := color.RGBA()
		fmt.Printf("#%02x%02x%02x\n", r/256, g/256, b/256)
	}
}
