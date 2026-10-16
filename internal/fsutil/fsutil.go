package fsutil

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func Exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func HasPNGExt(p string) bool {
	return strings.EqualFold(path.Ext(p), ".png")
}

func SavePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		return fmt.Errorf("encode %s: %w", path, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}

// filename without directory and extension
func Stem(path string) string {
	name := filepath.Base(path)
	return strings.TrimSuffix(name, filepath.Ext(name))
}

var chunks = regexp.MustCompile(`\d+|\D+`)

// compares filenames so that numbers are ordered by value
// (frame-2 < frame-10) and text is compared ignoring case
func NaturalLess(a, b string) bool {
	chunksA := chunks.FindAllString(Stem(a), -1)
	chunksB := chunks.FindAllString(Stem(b), -1)

	for i := range min(len(chunksA), len(chunksB)) {
		numA, errA := strconv.Atoi(chunksA[i])
		numB, errB := strconv.Atoi(chunksB[i])
		if errA == nil && errB == nil {
			if numA != numB {
				return numA < numB
			}
			continue
		}

		textA, textB := strings.ToLower(chunksA[i]), strings.ToLower(chunksB[i])
		if textA != textB {
			return textA < textB
		}
	}
	return len(chunksA) < len(chunksB)
}
