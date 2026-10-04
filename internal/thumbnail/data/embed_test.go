package data

import (
	"bytes"
	"image"
	"image/png"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pngembed "github.com/sabhiram/png-embed"
)

func TestEmbedTextualData(t *testing.T) {
	var refTime = time.Date(2026, 1, 2, 15, 4, 5, 0, time.Local)

	type testStruct struct {
		Software   string    `tEXt:""`
		ThumbMTime time.Time `tEXt:"Thumb::MTime"`
		ThumbSize  int       `tEXt:"Thumb::Size"`
		ThumbURI   string    `tEXt:"Thumb::URI"`
	}

	type testStructNestingInner struct {
		MTime time.Time `tEXt:"Thumb::MTime"`
		Size  int       `tEXt:"Thumb::Size"`
		URI   string    `tEXt:"Thumb::URI"`
	}

	type testStructNesting struct {
		Software  string                 `tEXt:""`
		Thumbnail testStructNestingInner `tEXt:""`
	}

	tests := []struct {
		name   string
		data   any
		expect map[string][]byte
	}{
		{
			"generic",
			testStruct{
				Software:   "Foo bar",
				ThumbURI:   "/tmp/foo.txt",
				ThumbSize:  480,
				ThumbMTime: refTime,
			},
			map[string][]byte{
				"Software":     []byte("Foo bar"),
				"Thumb::MTime": []byte(strconv.Itoa(int(refTime.Unix()))),
				"Thumb::Size":  []byte(strconv.Itoa(480)),
				"Thumb::URI":   []byte("/tmp/foo.txt"),
			},
		},
		{
			"nested-struct",
			testStructNesting{
				Software: "Foo bar",
				Thumbnail: testStructNestingInner{
					URI:   "/tmp/foo.txt",
					Size:  480,
					MTime: refTime,
				},
			},
			map[string][]byte{
				"Software":     []byte("Foo bar"),
				"Thumb::MTime": []byte(strconv.Itoa(int(refTime.Unix()))),
				"Thumb::Size":  []byte(strconv.Itoa(480)),
				"Thumb::URI":   []byte("/tmp/foo.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			png.Encode(&buf, image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}}))
			file := buf.Bytes()

			out, err := EmbedTextualData(file, tt.data)
			if err != nil {
				t.Errorf("error while embedding data: %s", err)
			}

			actual, err := pngembed.Extract(out)
			if err != nil {
				t.Errorf("error extracting data again: %s", err)
			}

			if diff := cmp.Diff(tt.expect, actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
