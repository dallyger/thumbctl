package thumbnail

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/thumbnail/data"
)

// Find thumbnail for file. Returns empty string without error on cache miss.
func Find(file metadata.Meta, size string) (string, error) {

	if IsThumbnail(file.AbsolutePath) {
		return file.AbsolutePath, nil
	}

	if _, err := SizeToInt(size); err != nil {
		return "", err
	}

	thumbPath := getPath(file, size)
	if IsCurrent(thumbPath) {
		return thumbPath, nil
	} else {
		return "", nil
	}

}

func GetEmbeddedMeta(path string) (meta metadata.Meta, err error) {

	err = data.UnmarshalTextualDataFromPath(path, &meta)
	if err != nil {
		return
	}

	u, err := url.Parse(meta.CanonicalUri)
	if err != nil {
		err = fmt.Errorf("invalid stored canonical URI: %s", err)
		return
	}
	if u.Scheme != "file" {
		err = fmt.Errorf("unsupported URI in thumbnail")
		return
	}

	meta.AbsolutePath = u.Path
	meta.FilePath = meta.AbsolutePath
	meta.UriHash = fmt.Sprintf("%x", md5.Sum([]byte(meta.CanonicalUri)))
	return
}

func getPath(meta metadata.Meta, size string) string {
	return filepath.Join(
		MustGetRootDir(),
		size,
		meta.UriHash+".png",
	)
}

// Checks if the thumbnail must be (re-)generated.
func IsCurrent(path string) bool {
	thumb, err := GetEmbeddedMeta(path)
	if err != nil {
		// Failed loading embedded data. Malformed file or missing.
		// Either way, not a usable thumbnail.
		return false
	}

	file, err := os.Stat(thumb.AbsolutePath)
	switch {
	case err == nil:

		if thumb.MTime.IsZero() {
			return false
		}

		if thumb.MTime.Truncate(time.Second) != file.ModTime().Truncate(time.Second) {
			return false
		}

		if thumb.FileSize != 0 && thumb.FileSize != file.Size() {
			return false
		}

		return true

	case errors.Is(err, fs.ErrNotExist):
		// does not exist anymore.
		return false
	default:
		// unknown: permission denied, I/O error, etc.
		// We cannot safely evaluate if the file is still current. Just to be
		// safe, let's delete it. Even it will be re-generated shortly after.
		return false
	}

}

// Checks if the file is located inside the thumbnails directory.
func IsThumbnail(path string) bool {
	rel, err := filepath.Rel(MustGetRootDir(), path)

	if err != nil {
		return false
	}

	return filepath.IsLocal(rel)
}
