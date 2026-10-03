package mimetype

import (
	"fmt"
	"os"

	"github.com/go-freedesktop/mime"
)

func Aliases(m string) []string {
	return mime.Default().Aliases(m)
}

func FromPath(p string) (m string, e error) {
	f, err := os.Open(p)

	if err != nil {
		e = fmt.Errorf("failed opening file: %s", err)
		return
	}

	head := make([]byte, 256)
	n, err := f.Read(head)

	if err != nil {
		e = fmt.Errorf("failed reading file: %s", err)
		return
	}

	return FromBytes(p, head[:n])
}
func FromBytes(p string, b []byte) (m string, e error) {
	m = mime.TypeByNameAndContent(p, b)
	return
}
