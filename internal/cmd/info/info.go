package info

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"

	"go.vnbr.de/thumbctl/internal/metadata"
	"go.vnbr.de/thumbctl/internal/thumbnail/data"
)

type InfoCmd struct {
	Path   string `arg:"" help:"Path to the file to check."`
	Format string `enum:"human,json" default:"human" help:"Choose human or json."`
}

type File struct {
	Meta TextualData `json:"meta" tEXt:""`
	Path string      `json:"absolute_path" human:"Thumbnail Path"`
}

type TextualData struct {
	Title       string    `json:"title" tEXt:""`
	Author      string    `json:"author" tEXt:""`
	Comment     string    `json:"comment" tEXt:""`
	Description string    `json:"description" tEXt:""`
	Software    string    `json:"software" tEXt:""`
	MTime       time.Time `json:"modified_at" tEXt:"Thumb::MTime"`
	Size        int       `json:"size" tEXt:"Thumb::Size"`
	URI         string    `json:"canonical_uri" tEXt:"Thumb::URI"`
}

type UnsupportedFormatError struct {
	Format string
}

func (e UnsupportedFormatError) Error() string {
	return fmt.Sprintf("the requested format [%s] is not supported", e.Format)
}

func (cmd *InfoCmd) Run() error {
	meta, err := metadata.FromPath(cmd.Path)
	if err != nil {
		return err
	}

	if !meta.IsThumbnail {
		return fmt.Errorf("given path is not a thumbnail file")
	}

	file := File{
		Path: meta.AbsolutePath,
	}

	if err = data.UnmarshalTextualDataFromFile(meta.AbsolutePath, &file); err != nil {
		return fmt.Errorf("failed reading metadata of thumbnail at %s: %s", meta.AbsolutePath, err)
	}

	switch cmd.Format {
	case "human":
		return outputHuman(file)
	case "json":
		return outputJson(file)
	default:
		return UnsupportedFormatError{Format: cmd.Format}
	}
}

func outputHuman(file File) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\n\n", outputHumanStruct(reflect.ValueOf(file)))
	w.Flush()
	return nil
}

func outputHumanKey(t reflect.StructField) string {
	if n := t.Tag.Get("human"); n != "" {
		return n
	} else if n := t.Tag.Get("tEXt"); n != "" {
		return n
	} else {
		return t.Name
	}
}
func outputHumanStruct(v reflect.Value) string {
	t := v.Type()
	out := ""
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		out = out + outputHumanRow(fv, outputHumanKey(ft))
	}
	return strings.TrimSpace(out)
}

func outputHumanRow(v reflect.Value, prefix string) string {
	switch v.Kind() {
	case reflect.Int:
		return fmt.Sprintf("%s:\t%d\n", prefix, v.Int())
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			return fmt.Sprintf("%s:\t%s\n", prefix, v.Interface().(time.Time).Format(time.RFC1123))
		default:
			return outputHumanStruct(v) + "\n"
		}
	default:
		return fmt.Sprintf("%s:\t%s\n", prefix, v.String())
	}
}

func outputJson(file File) error {
	if out, err := json.Marshal(file); err != nil {
		return err
	} else {
		fmt.Println(string(out))
		return nil
	}
}
