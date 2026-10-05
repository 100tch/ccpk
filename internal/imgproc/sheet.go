package imgproc

import (
	"image"
	"image/draw"
	"math"
)

type Sheet struct {
	Image   *image.RGBA
	Columns int
	Rows    int
}

func BuildSheet(cells []*image.RGBA, columns int) Sheet {
	if columns <= 0 {
		columns = int(math.Ceil(math.Sqrt(float64(len(cells)))))
	}
	rows := (len(cells) + columns - 1) / columns
	cellSize := cells[0].Bounds().Size()

	img := image.NewRGBA(image.Rect(0, 0, cellSize.X*columns, cellSize.Y*rows))
	for i, cell := range cells {
		origin := image.Pt((i%columns)*cellSize.X, (i/columns)*cellSize.Y)
		target := image.Rectangle{Min: origin, Max: origin.Add(cellSize)}
		draw.Draw(img, target, cell, image.Point{}, draw.Src)
	}
	return Sheet{Image: img, Columns: columns, Rows: rows}
}
