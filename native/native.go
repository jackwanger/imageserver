// Package native provides a Server that uses native Go Image.
package native

import (
	"bytes"
	"fmt"
	"image"

	"github.com/pierrre/imageserver"
)

// Server is an Image Server that uses the native Go Image.
//
// Steps:
//
// - Get the Image from the underlying Server.
//
// - Decode (from raw data to Go image).
//
// - Process (Go image).
//
// - Encode (from Go image to raw data).
type Server struct {
	imageserver.Server
	Processor Processor
}

// Get implements Server.
func (server *Server) Get(params imageserver.Params) (*imageserver.Image, error) {
	im, err := server.Server.Get(params)
	if err != nil {
		return nil, err
	}
	nim, err := Decode(im)
	if err != nil {
		return nil, err
	}
	nim, err = server.Processor.Process(nim, params)
	if err != nil {
		return nil, err
	}
	enc, format, err := getEncoderFormat(im, params)
	if err != nil {
		return nil, err
	}
	im, err = encode(enc, format, nim, params)
	if err != nil {
		return nil, err
	}
	return im, nil
}

// Decode decodes a raw Image to a native Image.
func Decode(im *imageserver.Image) (image.Image, error) {
	nim, format, err := image.Decode(bytes.NewReader(im.Data))
	if err != nil {
		return nil, &imageserver.ImageError{Message: err.Error()}
	}
	if format != im.Format {
		return nil, &imageserver.ImageError{Message: fmt.Sprintf("format \"%s\" does not match decoded format \"%s\"", im.Format, format)}
	}
	return nim, nil
}

func getEncoderFormat(im *imageserver.Image, params imageserver.Params) (Encoder, string, error) {
	format, fromParams, err := getFormat(im, params)
	if err != nil {
		return nil, "", err
	}
	enc, err := getEncoder(format)
	if err != nil {
		if fromParams {
			err = &imageserver.ParamError{Param: "format", Message: err.Error()}
		} else {
			err = &imageserver.ImageError{Message: err.Error()}
		}
		return nil, "", err
	}
	return enc, format, nil
}

func getFormat(im *imageserver.Image, params imageserver.Params) (string, bool, error) {
	if !params.Has("format") {
		return im.Format, false, nil
	}
	format, err := params.GetString("format")
	if err != nil {
		return "", false, err
	}
	return format, true, nil
}

func encode(enc Encoder, format string, nim image.Image, params imageserver.Params) (*imageserver.Image, error) {
	data, err := enc.Encode(nim, params)
	if err != nil {
		return nil, err
	}
	im := &imageserver.Image{
		Format: format,
		Data:   data,
	}
	return im, nil
}

// Processor represents a native Go Image processor.
type Processor interface {
	Process(image.Image, imageserver.Params) (image.Image, error)
}

// ProcessorFunc is a Processor func.
type ProcessorFunc func(image.Image, imageserver.Params) (image.Image, error)

// Process implements Processor.
func (f ProcessorFunc) Process(nim image.Image, params imageserver.Params) (image.Image, error) {
	return f(nim, params)
}

// Encoder represents an Image encoder.
//
// It encodes a native Image to []byte.
//
// An Encoder must encode to only one specific format.
type Encoder interface {
	Encode(image.Image, imageserver.Params) ([]byte, error)
}

// EncoderFunc is a Encoder func.
type EncoderFunc func(image.Image, imageserver.Params) ([]byte, error)

// Encode implements Encoder.
func (f EncoderFunc) Encode(nim image.Image, params imageserver.Params) ([]byte, error) {
	return f(nim, params)
}

var encoders = make(map[string]Encoder)

// RegisterEncoder registers an Encoder for a format.
func RegisterEncoder(format string, encoder Encoder) {
	encoders[format] = encoder
}

func getEncoder(format string) (Encoder, error) {
	encoder, ok := encoders[format]
	if !ok {
		return nil, fmt.Errorf("no registered encoder for format \"%s\"", format)
	}
	return encoder, nil
}
