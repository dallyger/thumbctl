package thumbnail

import (
	"bytes"
	"fmt"
	"image/png"
	"slices"
)

func SizeToInt(v string) (int, error) {
	switch v {
	case "normal":
		return 128, nil
	case "large":
		return 256, nil
	case "x-large":
		return 512, nil
	case "xx-large":
		return 1024, nil
	default:
		return 0, fmt.Errorf("invalid size definition: %s", v)
	}
}

func IntToSize(v int) (string, error) {
	switch v {
	case 128:
		return "normal", nil
	case 256:
		return "large", nil
	case 512:
		return "x-large", nil
	case 1024:
		return "xx-large", nil
	default:
		return "", fmt.Errorf("invalid pixel size: %d", v)
	}
}

func PngToSize(b []byte) (string, error) {
	pngConf, err := png.DecodeConfig(bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("invalid png: %s", err)
	}

	if IsValidSize(pngConf.Width) {
		return IntToSize(pngConf.Width)
	} else {
		return IntToSize(pngConf.Height)
	}
}

func IsValidSize(v int) bool {
	return slices.Contains([]int{128, 256, 512, 1024}, v)
}

func IsValidSizeLabel(v string) bool {
	return slices.Contains([]string{"normal", "large", "x-large", "xx-large"}, v)
}
