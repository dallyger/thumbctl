package metadata

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Meta struct {
	AbsolutePath string
	CanonicalUri string `tEXt:"Thumb::URI"`
	FilePath     string
	FileSize     int64     `tEXt:"Thumb::Size"`
	MTime        time.Time `tEXt:"Thumb::MTime"`
	Software     string    `tEXt:"Software"`
	UriHash      string
}

func FromPath(path string) (Meta, error) {
	var err error
	meta := Meta{
		FilePath: path,
	}

	if path, err := filepath.Abs(path); err == nil {
		meta.AbsolutePath = path
	} else {
		return meta, err
	}

	var info os.FileInfo
	if info, err = os.Stat(meta.AbsolutePath); err == nil {
		meta.MTime = info.ModTime()
		meta.FileSize = info.Size()
	}

	u := url.URL{Path: meta.AbsolutePath}
	meta.CanonicalUri = "file://" + u.EscapedPath()
	meta.UriHash = fmt.Sprintf("%x", md5.Sum([]byte(meta.CanonicalUri)))

	return meta, nil
}
