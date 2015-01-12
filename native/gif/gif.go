// Package gif provides a GIF Encoder.
package gif

import (
	"bytes"
	"image"
	"image/gif"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
)

// Encoder encodes an Image to GIF.
type Encoder struct {
}

// Encode implements Encoder.
func (e *Encoder) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := gif.Encode(buf, nim, &gif.Options{})
	return buf.Bytes(), err
}

func init() {
	imageserver_native.RegisterEncoder("gif", &Encoder{})
}
