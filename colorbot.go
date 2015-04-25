// Package colorbot provides image analysis routines to determine the dominant
// colors in images.
package colorbot

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/soniakeys/quant/median"
)

// DominantColors returns the N most dominant colors in the given image.
//
// This uses the median cut quantization algorithm.
func DominantColors(img image.Image, n int) color.Palette {
	return median.Quantizer(n).Palette(img).ColorPalette()
}

var (
	// ErrImageTooLarge is returned when the image is too large to be processed.
	ErrImageTooLarge = errors.New("image is too large")
)

// DecodeImage decodes the given reader as either a GIF, JPEG, or PNG image.
//
// If the image is larger than maxBytes or maxPixels, it returns
// ErrImageTooLarge.
func DecodeImage(r io.Reader, maxBytes, maxPixels int64) (image.Image, error) {
	// limit images to maxBytes
	lr := &io.LimitedReader{
		R: r,
		N: maxBytes,
	}

	// read the header
	header := make([]byte, maxHeaderSize)
	if _, err := io.ReadFull(lr, header); err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	// parse just the image size
	hr := bytes.NewReader(header)
	config, _, err := image.DecodeConfig(hr)
	if err != nil {
		return nil, err
	}

	// check to see the image isn't too big
	pixels := int64(config.Height) * int64(config.Width)
	if pixels > maxPixels {
		return nil, ErrImageTooLarge
	}

	// recombine the image
	_, _ = hr.Seek(0, 0)
	ir := io.MultiReader(hr, lr)

	// decode the image
	img, _, err := image.Decode(ir)
	if err == io.ErrUnexpectedEOF {
		return nil, ErrImageTooLarge
	}
	return img, err
}

const (
	maxHeaderSize = 4096 // bytes
)
