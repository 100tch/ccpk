package fsutil

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/fs"
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

func SavePNG(path string, img image.Image, overwrite bool) error {
	file, err := Create(path, overwrite)
	if err != nil {
		return err
	}

	if err := png.Encode(file, img); err != nil {
		file.Close()
		os.Remove(path)
		return fmt.Errorf("encode %s: %w", path, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	return nil
}

func Create(path string, overwrite bool) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(path, flags, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return nil, fmt.Errorf("%s already exists (use -force to overwrite)", path)
	}
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", path, err)
	}

	return file, nil
}

func SamePath(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	return errA == nil && errB == nil && absA == absB
}

func Inside(parent, child string) bool {
	absParent, errP := filepath.Abs(parent)
	absChild, errC := filepath.Abs(child)
	if errP != nil || errC != nil {
		return false
	}
	rel, err := filepath.Rel(absParent, absChild)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
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
