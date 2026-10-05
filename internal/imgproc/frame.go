package imgproc

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Frame struct {
	Path  string
	Image image.Image
}

func Load(paths []string) ([]Frame, error) {
	loaded := make([]Frame, 0, len(paths))
	for _, path := range paths {
		img, err := decodeFile(path)
		if err != nil {
			return nil, err
		}
		loaded = append(loaded, Frame{Path: path, Image: img})
	}
	return loaded, nil
}

func SavePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func decodeFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return img, nil
}
