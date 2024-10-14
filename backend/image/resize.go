package image

import (
	"image"

	"github.com/disintegration/imaging"
)

type ResizedImage struct {
	Img   image.Image
	Width int
}

func Resize(img image.Image, width int) ResizedImage {
	return ResizedImage{
		Img:   imaging.Resize(img, width, 0, imaging.Lanczos),
		Width: width,
	}
}
