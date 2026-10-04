package frames

import (
	"fmt"
	"strings"
)

type Direction string

const (
	north     Direction = "north"
	northEast Direction = "north_east"
	east      Direction = "east"
	southEast Direction = "south_east"
	south     Direction = "south"
	southWest Direction = "south_west"
	west      Direction = "west"
	northWest Direction = "north_west"
)

var (
	Dir4 = []Direction{north, east, south, west}
	Dir8 = []Direction{north, northEast, east, southEast, south, southWest, west, northWest}
)

var fileSuffixes = map[string]Direction{
	"N":  north,
	"NE": northEast,
	"E":  east,
	"SE": southEast,
	"S":  south,
	"SW": southWest,
	"W":  west,
	"NW": northWest,
}

// pick, for each direction in want, the file carrying that
// direction's suffix and returns the files in the order of want
//
// every file in paths must have a direction suffix and all of them must
// share one base name.
func OrderByDirection(paths []string, want []Direction) ([]string, error) {
	byDirection := make(map[Direction]string)
	var commonBase string

	for i, path := range paths {
		base, direction, ok := splitDirection(stem(path))
		if !ok {
			return nil, fmt.Errorf("[error] file %s has no recognized direction suffix (expected one of -N -NE -E -SE -S -SW -W -NW)", path)
		}
		if i == 0 {
			commonBase = base
		} else if base != commonBase {
			return nil, fmt.Errorf("[error] file %s base name %q does not match other files' base %q — directional frames must share one base name", path, base, commonBase)
		}
		if _, exists := byDirection[direction]; exists {
			return nil, fmt.Errorf("[error] direction %s has more than one matching file (last: %s)", direction, path)
		}
		byDirection[direction] = path
	}

	ordered := make([]string, 0, len(want))
	var missing []string
	for _, direction := range want {
		if path, ok := byDirection[direction]; ok {
			ordered = append(ordered, path)
		} else {
			missing = append(missing, string(direction))
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing frames for direction(s): %s", strings.Join(missing, ", "))
	}
	return ordered, nil
}

// splits filename with prefix into the base name and direction
// ok is false if there is no known direction suffix
func splitDirection(name string) (base string, direction Direction, ok bool) {
	dash := strings.LastIndex(name, "-")
	if dash == -1 {
		return name, "", false
	}
	direction, ok = fileSuffixes[strings.ToUpper(name[dash+1:])]
	if !ok {
		return name, "", false
	}
	return name[:dash], direction, true
}
