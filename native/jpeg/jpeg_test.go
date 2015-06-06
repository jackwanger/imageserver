package jpeg

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
	imageserver_native_test "github.com/pierrre/imageserver/native/_test"
)

var _ imageserver_native.Encoder = &Encoder{}

func TestEncoder(t *testing.T) {
	testEncoder(t, &Encoder{})
}

func TestEncoderDefaultQuality(t *testing.T) {
	encoder := &Encoder{
		DefaultQuality: 90,
	}
	testEncoder(t, encoder)
}

func TestEncoderQuality(t *testing.T) {
	params := imageserver.Params{
		"quality": 90,
	}
	testEncoderParams(t, &Encoder{}, params)
}

func TestEncoderErrorQuality(t *testing.T) {
	im := imageserver_native_test.NewImage()
	encoder := &Encoder{}
	for _, quality := range []interface{}{"foo", -1, 101} {
		_, err := encoder.Encode(im, imageserver.Params{"quality": quality})
		if err == nil {
			t.Fatal("no error")
		}
		errParam, ok := err.(*imageserver.ParamError)
		if !ok {
			t.Fatalf("wrong error type: %T", err)
		}
		if errParam.Param != "quality" {
			t.Fatalf("wrong param: %s", errParam.Param)
		}
	}
}

func testEncoder(t *testing.T, encoder *Encoder) {
	imageserver_native_test.TestEncoder(t, encoder, "jpeg")
}

func testEncoderParams(t *testing.T, encoder *Encoder, params imageserver.Params) {
	imageserver_native_test.TestEncoderParams(t, encoder, params, "jpeg")
}
