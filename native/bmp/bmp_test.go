package bmp

import (
	"testing"

	imageserver_native "github.com/pierrre/imageserver/native"
	imageserver_native_test "github.com/pierrre/imageserver/native/_test"
)

var _ imageserver_native.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	imageserver_native_test.TestEncoder(t, &Encoder{}, "bmp")
}
