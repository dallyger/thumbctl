package data

import (
	"bytes"
	"image"
	"image/png"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pngembed "github.com/sabhiram/png-embed"
)

var refTime = time.Date(2026, 1, 2, 15, 4, 5, 0, time.Local)

func TestUnmarshalTextualData(t *testing.T) {

	type testStruct struct {
		Comment    string    `tEXt:""`
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
		Comment   string                 `tEXt:""`
		Thumbnail testStructNestingInner `tEXt:""`
	}

	tests := []struct {
		name   string
		data   map[string]any
		actual any
		expect any
	}{
		{
			"generic",
			map[string]any{
				"Comment":      "Foo bar",
				"Thumb::MTime": refTime.Unix(),
				"Thumb::Size":  480,
				"Thumb::URI":   "/tmp/foo.txt",
			},
			&testStruct{},
			&testStruct{
				Comment:    "Foo bar",
				ThumbURI:   "/tmp/foo.txt",
				ThumbSize:  480,
				ThumbMTime: refTime,
			},
		},
		{
			"nested-struct",
			map[string]any{
				"Comment":      "Foo bar",
				"Thumb::MTime": refTime.Unix(),
				"Thumb::Size":  480,
				"Thumb::URI":   "/tmp/foo.txt",
			},
			&testStructNesting{},
			&testStructNesting{
				Comment: "Foo bar",
				Thumbnail: testStructNestingInner{
					URI:   "/tmp/foo.txt",
					Size:  480,
					MTime: refTime,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := UnmarshalTextualData(dummyPng(tt.data), tt.actual); err != nil {
				t.Fail()
				t.Errorf("error: %s", err)
			}

			if diff := cmp.Diff(tt.expect, tt.actual); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})
	}
}

func dummyPng(data map[string]any) []byte {
	var buf bytes.Buffer
	png.Encode(&buf, image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}}))

	file := buf.Bytes()
	for k, v := range data {
		file, _ = pngembed.Embed(file, k, v)
	}

	return file
}
