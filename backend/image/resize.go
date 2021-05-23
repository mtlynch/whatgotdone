package image

import (
	"image"
	"io"

	"github.com/disintegration/imaging"
)

type ResizedImage struct {
	Img   image.Image
	Width int
}

func Resize(img image.Image, resizeWidths []int) []ResizedImage {
	results := []ResizedImage{}
	for _, width := range resizeWidths {
		results = append(results, ResizedImage{
			Img:   imaging.Resize(img, width, 0, imaging.Lanczos),
			Width: width,
		})
	}
	return results
}

func ResizeFile(r io.Reader, resizeWidths []int) ([]ResizedImage, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return []ResizedImage{}, err
	}
	return Resize(img, resizeWidths), nil
}
