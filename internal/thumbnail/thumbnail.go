package thumbnail

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type Thumbnail struct {
	AbsolutePath string
	CanonicalUri string
	OriginalPath string
	UriHash      string
	IsThumbnail  bool
}

func FromPath(path string) (Thumbnail, error) {
	thumb := Thumbnail{
		OriginalPath: path,
	}

	if path, err := filepath.Abs(path); err == nil {
		thumb.AbsolutePath = path
	} else {
		return thumb, err
	}

	// Check if file is already a thumbnail and return itself.
	// As per the freedesktop thumbnail spec:, we "must load and use these files directly."
	if IsThumbnail(thumb.AbsolutePath) {
		thumb.CanonicalUri = "file://" + thumb.AbsolutePath
		thumb.UriHash = strings.Split(filepath.Base(thumb.AbsolutePath), ".")[0]
		thumb.IsThumbnail = true
		return thumb, nil
	}

	u := url.URL{Path: thumb.AbsolutePath}
	thumb.CanonicalUri = "file://" + u.EscapedPath()
	thumb.UriHash = fmt.Sprintf("%x", md5.Sum([]byte(thumb.CanonicalUri)))

	return thumb, nil
}
