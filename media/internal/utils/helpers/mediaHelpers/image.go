package mediahelpers

import (
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

func BlurImage(img image.Image, sigma float64) ([]byte, image.Image, error) {
	blurredImg := imaging.Blur(img, sigma)
	var buf bytes.Buffer
	err := imaging.Encode(&buf, blurredImg, imaging.JPEG, imaging.JPEGQuality(90))
	if err != nil {
		return nil, nil, fmt.Errorf("error encoding blurred image: %v", err)
	}

	return buf.Bytes(), blurredImg, nil
}
