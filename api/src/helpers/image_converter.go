package helpers

import (
	"bytes"
	"github.com/chai2010/webp"
	"image"
)

func decodeImage(fileBytes []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ConvertImageToWebp(fileBytes []byte) ([]byte, error) {
	img, err := decodeImage(fileBytes)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	options := &webp.Options{
		Lossless: false,
		Quality:  80,
	}

	err = webp.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
