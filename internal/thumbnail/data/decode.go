package data

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	pngembed "github.com/sabhiram/png-embed"
)

func UnmarshalTextualDataFromPath(path string, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return errors.New("UnmarshalTextualDataFromFile needs a non-nil pointer to a struct")
	}

	bs, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading file: %s", err)
	}

	if err := UnmarshalTextualData(bs, v); err != nil {
		return fmt.Errorf("failed reading metadata of file: %s", err)
	}

	return nil
}

func UnmarshalTextualData(b []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return errors.New("UnmarshalTextualData needs a non-nil pointer to a struct")
	}

	png, err := pngembed.Extract(b)
	if err != nil {
		return err
	}

	return decodeStruct(rv.Elem(), png)
}

func decodeStruct(rv reflect.Value, png map[string][]byte) error {
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		tag, ok := ft.Tag.Lookup("tEXt")
		if !ok || !ft.IsExported() {
			continue
		}

		if tag == "" {
			tag = ft.Name
		}

		if err := decodeValue(rv.Field(i), png, tag); err != nil {
			return fmt.Errorf("field %s: %s", ft.Name, err)
		}
	}

	return nil
}
func decodeValue(fv reflect.Value, png map[string][]byte, tag string) error {
	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := intDecoder(png[tag]); err == nil {
			fv.SetInt(int64(v))
		} else {
			return fmt.Errorf("unsupported value %s for type %s", string(png[tag]), fv.Kind())
		}

	case reflect.String:
		fv.SetString(string(png[tag]))

	case reflect.Struct:
		switch ptr := fv.Addr().Interface().(type) {

		case *time.Time:
			if v, err := timeDecoder(png[tag]); err != nil {
				return fmt.Errorf("unsupported time.Time value %s: %s", v.Format(time.RFC3339), err)
			} else {
				*ptr = v
			}

		default:
			return decodeStruct(fv, png)
		}

	default:
		return fmt.Errorf("unsupported kind %s", fv.Kind())
	}

	return nil
}

func intDecoder(value []byte) (int64, error) {
	if value == nil {
		return 0, nil
	}
	v, err := strconv.Atoi(string(value))
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func timeDecoder(value []byte) (time.Time, error) {
	if v, err := intDecoder(value); err != nil {
		return time.Time{}, fmt.Errorf("unsupported time.Time value %d: %s", v, err)
	} else if v == 0 {
		return time.Time{}, nil
	} else {
		return time.Unix(v, 0), nil
	}
}
