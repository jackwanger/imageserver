// Package tiff provides a TIFF Encoder.
package tiff

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
	"golang.org/x/image/tiff"
)

// Encoder encodes an Image to TIFF.
type Encoder struct {
}

var opts = &tiff.Options{
	Compression: tiff.Deflate,
	Predictor:   true,
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := tiff.Encode(buf, nim, opts)
	return buf.Bytes(), err
}

func init() {
	imageserver_native.RegisterEncoder("tiff", &Encoder{})
}
