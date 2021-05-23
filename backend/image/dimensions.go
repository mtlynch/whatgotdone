package image

import (
	"image"
)

func imageWidth(img image.Image) int {
	return img.Bounds().Max.X
}

func imageHeight(img image.Image) int {
	return img.Bounds().Max.Y
}
