package thumbnailer

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

type Entries map[string]Entry

func (c Entries) Sorted() iter.Seq2[string, Entry] {
	return func(yield func(string, Entry) bool) {
		for _, k := range slices.Sorted(maps.Keys(c)) {
			if !yield(k, c[k]) {
				return
			}
		}
	}
}

type Entry struct {
	Exec     string   `ini:"Exec"`
	TryExec  string   `ini:"TryExec"`
	MimeType []string `ini:"MimeType"`
}

func FromFile(f []byte) (Entry, error) {
	e := Entry{}

	c, err := ini.LoadSources(ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}, f)
	if err != nil {
		return e, err
	}

	s := c.Section("Thumbnailer Entry")
	if v, ok := getValue(s, "Exec"); ok {
		e.Exec = v
	}

	if v, ok := getValue(s, "TryExec"); ok {
		e.TryExec = v
	}

	if v, ok := getValue(s, "MimeType"); ok {
		e.MimeType = strings.Split(strings.Trim(v, ";"), ";")
	}

	return e, nil
}

func FromPath(f string) (Entry, error) {
	if h, err := os.Open(f); err != nil {
		return Entry{}, err
	} else {
		return FromReader(h)
	}
}

func FromReader(f io.Reader) (Entry, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return Entry{}, err
	}
	return FromFile(b)
}

func GetEntries() (Entries, []error) {
	entries := make(map[string]Entry)
	var errors []error
	for _, f := range GetEntryFiles() {
		if e, err := FromPath(f); err != nil {
			errors = append(errors, fmt.Errorf("%s: %s", f, err))
		} else {
			entries[f] = e
		}
	}
	return entries, errors
}
func getValue(s *ini.Section, key string) (value string, ok bool) {
	if val, err := s.GetKey(key); err == nil {
		return val.String(), true
	} else {
		return "", false
	}
}
