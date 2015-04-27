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
//     $ colorbot https://www.google.com/images/srpr/logo11w.png
//     #000000
//     #009553
//     #0c60a7
//     #166bed
//     #c56937
//
// Colorbot supports GIF, JPEG, and PNG images.
package main

import (
	"flag"
	"fmt"
	_ "image/gif"  // support GIF images
	_ "image/jpeg" // support JPEG images
	_ "image/png"  // support PNG images
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/codahale/colorbot"
)

func main() {
	var (
		n         = flag.Int("n", 5, "number of colors to return")
		maxBytes  = flag.Int64("maxbytes", 10*1024*1024, "max image size in bytes")
		maxPixels = flag.Int64("maxpixels", 10*1024*1024, "max image size in pixels")
	)
	flag.Parse()

	var in io.Reader
	if s := flag.Args()[0]; s == "-" {
		in = os.Stdin
	} else if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		resp, err := http.Get(s)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		in = resp.Body
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
