// Package bmp provides a BMP Encoder.
package bmp

import (
	"bytes"
	"image"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
	"golang.org/x/image/bmp"
)

// Encoder encodes an Image to BMP.
type Encoder struct {
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := bmp.Encode(buf, nim)
	return buf.Bytes(), err
}

func init() {
	imageserver_native.RegisterEncoder("bmp", &Encoder{})
}
