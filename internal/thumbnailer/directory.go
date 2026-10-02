package thumbnailer

import (
	"path"
	"path/filepath"

	"go4.org/xdgdir"
)

// Get directories in which thumbnailers are installed to.
func GetLocations() []string {
	var dirs []string
	for _, dir := range xdgdir.Data.SearchPaths() {
		dirs = append(dirs, filepath.Join(dir, "thumbnailers"))
	}

	return dirs
}

func GetEntryFiles() []string {
	var files []string
	for _, dir := range GetLocations() {
		matches, _ := filepath.Glob(path.Join(dir, "*.thumbnailer"))
		for _, match := range matches {
			files = append(files, match)
		}
	}
	return files
}
