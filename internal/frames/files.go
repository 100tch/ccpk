package frames

import (
	"ccpk/internal/fsutil"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		return nil, fmt.Errorf("no frames matching pattern %q in %s", pattern, dir)
	}

	sort.SliceStable(files, func(i, j int) bool { return fsutil.NaturalLess(files[i], files[j]) })
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
