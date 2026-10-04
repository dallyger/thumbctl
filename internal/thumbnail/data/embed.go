package data

import (
	"fmt"
	"reflect"
	"time"

	pngembed "github.com/sabhiram/png-embed"
)

func EmbedTextualData(b []byte, v any) (out []byte, err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return out, fmt.Errorf("EmbedTextualData needs a struct")
	}

	out, err = embedStruct(b, rv)
	return
}

func embedStruct(b []byte, rv reflect.Value) ([]byte, error) {
	rt := rv.Type()
	var err error
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		tag, ok := ft.Tag.Lookup("tEXt")
		if tag == "" {
			tag = ft.Name
		}
		if !ok || !ft.IsExported() {
			continue
		}

		if b, err = embedValue(b, tag, rv.Field(i)); err != nil {
			return b, fmt.Errorf("field %s: %s", ft.Name, err)
		}
	}

	return b, nil
}

func embedValue(b []byte, tag string, rv reflect.Value) ([]byte, error) {
	var err error
	switch rv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b, err = pngembed.Embed(b, tag, rv.Int())

	case reflect.String:
		b, err = pngembed.Embed(b, tag, rv.String())

	case reflect.Struct:
		switch v := rv.Interface().(type) {

		case time.Time:
			b, err = pngembed.Embed(b, tag, v.Unix())
			if err != nil {
				return b, err
			}

		default:
			b, err = embedStruct(b, rv)
		}

	default:
		return b, fmt.Errorf("unsupported kind %s", rv.Kind())
	}

	return b, err

}
