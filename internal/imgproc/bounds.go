package imgproc

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"path/filepath"
)

func FullRect(frames []Frame) (image.Rectangle, error) {
	if len(frames) == 0 {
		return image.Rectangle{}, errors.New("no frames provided")
	}

	size := frames[0].Image.Bounds().Size()
	for _, frame := range frames[1:] {
		if got := frame.Image.Bounds().Size(); got != size {
			return image.Rectangle{}, fmt.Errorf("[error] frame %s is %dx%d, expected %dx%d (sizes must match without trimming)",
				frame.Path, got.X, got.Y, size.X, size.Y)
		}
	}
	return image.Rectangle{Max: size}, nil
}

func TrimRect(frames []Frame, threshold, padding int, debug io.Writer) (image.Rectangle, error) {
	var visible image.Rectangle
	for _, frame := range frames {
		box, ok := visibleBounds(frame.Image, threshold)
		if !ok {
			continue
		}
		if debug != nil {
			fmt.Fprintf(debug, "debug: %s -> bbox: (%d, %d, %d, %d)\n",
				filepath.Base(frame.Path), box.Min.X, box.Min.Y, box.Max.X, box.Max.Y)
		}
		visible = visible.Union(box)
	}
	if visible.Empty() {
		return image.Rectangle{}, errors.New("no non-transparent pixels found in any frame")
	}

	return image.Rect(
		max(visible.Min.X-padding, 0),
		max(visible.Min.Y-padding, 0),
		visible.Max.X+padding,
		visible.Max.Y+padding,
	), nil
}

func Crop(frames []Frame, rect image.Rectangle) []*image.RGBA {
	cells := make([]*image.RGBA, len(frames))
	for i, frame := range frames {
		cell := image.NewRGBA(image.Rectangle{Max: rect.Size()})
		draw.Draw(cell, cell.Bounds(), frame.Image, rect.Min, draw.Src)
		cells[i] = cell
	}
	return cells
}

func visibleBounds(img image.Image, threshold int) (box image.Rectangle, ok bool) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if int(alpha>>8) >= threshold {
				box = box.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	return box, !box.Empty()
}
