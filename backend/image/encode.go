package image

import (
	goimage "image"
	"image/jpeg"
	"io"
)

func Encode(img goimage.Image, w io.Writer) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 80})
}
