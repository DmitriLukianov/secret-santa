package image

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
)

const (
	maxDimension = 1200
	jpegQuality  = 85
)

// Optimize resizes the image so neither side exceeds maxDimension, then
// re-encodes it as JPEG at jpegQuality. Returns the bytes and "image/jpeg".
func Optimize(data []byte) ([]byte, string, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	b := src.Bounds()
	if b.Dx() > maxDimension || b.Dy() > maxDimension {
		src = imaging.Fit(src, maxDimension, maxDimension, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, src, imaging.JPEG, imaging.JPEGQuality(jpegQuality)); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "image/jpeg", nil
}
