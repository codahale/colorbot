package colorbot_test

import (
	"fmt"
	"image"
	_ "image/gif"  // support GIF images
	_ "image/jpeg" // support JPEG images
	_ "image/png"  // support PNG images
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/codahale/colorbot"
)

func TestDominantColors(t *testing.T) {
	actual := map[string][]string{}

	filepath.Walk("./test-assets", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}

		p := colorbot.DominantColors(img, 3)
		s := []string{}
		for _, c := range p {
			r, g, b, _ := c.RGBA()
			s = append(s, fmt.Sprintf("#%02x%02x%02x", r/256, g/256, b/256))
		}
		actual[path] = s

		return nil
	})

	expected := map[string][]string{
		"test-assets/black-and-yellow.gif": {"#000000", "#8e8e00", "#ffff00"},
		"test-assets/black-and-yellow.jpg": {"#000000", "#010008", "#c0c004"},
		"test-assets/etsy.png":             {"#000000", "#68310d", "#d5641c"},
		"test-assets/hodges-research.png":  {"#000000", "#e0bd78", "#ffffff"},
		"test-assets/icon-1.png":           {"#39454f", "#4a525a", "#b6b8ba"},
		"test-assets/icon-2.png":           {"#bb1b59", "#d32d6c", "#ed538d"},
		"test-assets/icon-3.png":           {"#283c77", "#43588f", "#a9b4ce"},
		"test-assets/icon-4.png":           {"#335752", "#6f9c36", "#fbf1d9"},
		"test-assets/icon-5.png":           {"#abb4b2", "#dfddd6", "#f4f4ef"},
		"test-assets/icon-6.png":           {"#7da03d", "#8ab042", "#c2d899"},
		"test-assets/icon-7.png":           {"#4282ea", "#538dec", "#8cb2f2"},
		"test-assets/interlaced.png":       {"#000000", "#060606", "#464547"},
		"test-assets/low-contrast.png":     {"#ffed50", "#fffbd0", "#ffffff"},
	}

	for file, v := range actual {
		if want := expected[file]; !reflect.DeepEqual(v, want) {
			t.Errorf("Palette for %s was %v but expected %v", file, v, want)
		}
	}
}

func TestDecodeGIF(t *testing.T) {
	f, err := os.Open("./test-assets/black-and-yellow.gif")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = colorbot.DecodeImage(f, 1024*1024, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeJPEG(t *testing.T) {
	f, err := os.Open("./test-assets/black-and-yellow.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = colorbot.DecodeImage(f, 1024*1024, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodePNG(t *testing.T) {
	f, err := os.Open("./test-assets/hodges-research.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = colorbot.DecodeImage(f, 1024*1024, 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecodeTooManyBytes(t *testing.T) {
	f, err := os.Open("./test-assets/hodges-research.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = colorbot.DecodeImage(f, 1024, 1024*1024)
	if err != colorbot.ErrImageTooLarge {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestDecodeTooManyPixels(t *testing.T) {
	f, err := os.Open("./test-assets/hodges-research.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = colorbot.DecodeImage(f, 1024*1024, 10)
	if err != colorbot.ErrImageTooLarge {
		t.Errorf("Unexpected error: %v", err)
	}
}

func BenchmarkDominantColors(b *testing.B) {
	f, err := os.Open("./test-assets/icon-7.png")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		colorbot.DominantColors(img, 4)
	}
}
