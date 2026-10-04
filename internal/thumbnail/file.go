package thumbnail

import (
	"path/filepath"
)

// Checks if the file is located inside the thumbnails directory.
func IsThumbnail(path string) bool {
	rel, err := filepath.Rel(MustGetRootDir(), path)

	if err != nil {
		return false
	}

	return filepath.IsLocal(rel)
}
