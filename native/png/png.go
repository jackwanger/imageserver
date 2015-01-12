// Package png provides a PNG Encoder.
package png

import (
	"bytes"
	"image"
	"image/png"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
)

// Encoder encodes an Image to PNG.
type Encoder struct {
}

var encoder = &png.Encoder{
	CompressionLevel: png.BestCompression,
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := encoder.Encode(buf, nim)
	return buf.Bytes(), err
}

func init() {
	imageserver_native.RegisterEncoder("png", &Encoder{})
}
