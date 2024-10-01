package helpers

import (
	"bytes"
	"github.com/chai2010/webp"
	"image"
	_ "image/gif"  // Import support for GIF
	_ "image/jpeg" // Import support for JPEG
	_ "image/png"  // Import support for PNG
)

// decodeImage attempts to decode an image from the provided byte slice.
// It supports formats like JPEG, PNG, and GIF. Returns the decoded `image.Image` object or an error.
func decodeImage(fileBytes []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// ConvertImageToWebp converts an image from a supported format (JPEG, PNG, GIF) to WebP format.
// It returns the WebP image as a byte slice or an error if the conversion fails.
func ConvertImageToWebp(fileBytes []byte) ([]byte, error) {
	// Decode the original image from the provided bytes.
	img, err := decodeImage(fileBytes)
	if err != nil {
		return nil, err
	}

	// Create a buffer to store the WebP-encoded image.
	var buf bytes.Buffer

	// Set WebP encoding options, including quality and whether the encoding is lossless.
	options := &webp.Options{
		Lossless: false, // Set to lossy compression.
		Quality:  80,    // Set the image quality to 80%.
	}

	// Encode the decoded image into WebP format and write to the buffer.
	err = webp.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}

	// Return the WebP image bytes.
	return buf.Bytes(), nil
}
