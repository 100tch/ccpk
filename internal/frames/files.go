package frames

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// returns files in directory whose names match the glob pattern in natural order
func ListFiles(dir, pattern string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if matched, _ := filepath.Match(pattern, entry.Name()); matched {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("[error] no frames matching pattern %q in %s", pattern, dir)
	}

	sort.SliceStable(files, func(i, j int) bool { return naturalLess(files[i], files[j]) })
	return files, nil
}

// returns the direct subdirectories of dir, sorted by name
func ListSubdirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(dir, entry.Name()))
		}
	}
	return dirs, nil
}

// filename without directory and extension
func stem(path string) string {
	name := filepath.Base(path)
	return strings.TrimSuffix(name, filepath.Ext(name))
}

var chunks = regexp.MustCompile(`\d+|\D+`)

// compares filenames so that numbers are ordered by value
// (frame-2 < frame-10) and text is compared ignoring case
func naturalLess(a, b string) bool {
	chunksA := chunks.FindAllString(stem(a), -1)
	chunksB := chunks.FindAllString(stem(b), -1)

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
