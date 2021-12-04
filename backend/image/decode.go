package image

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type DecodeLimits struct {
	MinWidthPixels  int
	MinHeightPixels int
	MaxWidthPixels  int
	MaxHeightPixels int
}

func Decode(r io.Reader, dl DecodeLimits) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	if imageWidth(img) < dl.MinWidthPixels {
		return nil, fmt.Errorf("image width is too small: %dpx (min: %dpx)", imageWidth(img), dl.MinWidthPixels)
	}
	if imageWidth(img) > dl.MaxWidthPixels {
		return nil, fmt.Errorf("image width exceeds maximum: %dpx (max: %dpx)", imageWidth(img), dl.MaxWidthPixels)
	}
	if imageHeight(img) < dl.MinHeightPixels {
		return nil, fmt.Errorf("image height is too small: %dpx (min: %dpx)", imageHeight(img), dl.MinHeightPixels)
	}
	if imageHeight(img) > dl.MaxHeightPixels {
		return nil, fmt.Errorf("image height exceeds maximum: %dpx (max: %dpx)", imageHeight(img), dl.MaxHeightPixels)
	}

	return img, nil
}
