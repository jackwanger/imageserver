package nfntresize

import (
	_ "image/jpeg"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

var _ imageserver_native.Processor = &Processor{}

func TestProcessor(t *testing.T) {
	type TC struct {
		params             imageserver.Params
		expectedWidth      int
		expectedHeight     int
		expectedParamError string
	}
	for _, tc := range []TC{
		// no size
		{
			params:         imageserver.Params{},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params:         imageserver.Params{globalParam: imageserver.Params{}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":  0,
				"height": 0,
			}},
			expectedWidth:  1024,
			expectedHeight: 819,
		},
		// with size
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"height": 100,
			}},
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":  100,
				"height": 100,
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		// interpolation
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "nearest_neighbor",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "bilinear",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "bicubic",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "mitchell_netravali",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos2",
			}},
			expectedWidth: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "lanczos3",
			}},
			expectedWidth: 100,
		},
		// mode
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "resize",
			}},
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":  100,
				"height": 100,
				"mode":   "thumbnail",
			}},
			expectedWidth:  100,
			expectedHeight: 79, // 819 * 100 / 1024
		},
		// error
		{
			params:             imageserver.Params{globalParam: "invalid"},
			expectedParamError: globalParam,
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width": "invalid",
			}},
			expectedParamError: globalParam + ".width",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"height": "invalid",
			}},
			expectedParamError: globalParam + ".height",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width": -1,
			}},
			expectedParamError: globalParam + ".width",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"height": -1,
			}},
			expectedParamError: globalParam + ".height",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": false,
			}},
			expectedParamError: globalParam + ".interpolation",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width":         100,
				"interpolation": "invalid",
			}},
			expectedParamError: globalParam + ".interpolation",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
				"mode":  false,
			}},
			expectedParamError: globalParam + ".mode",
		},
		{
			params: imageserver.Params{globalParam: imageserver.Params{
				"width": 100,
				"mode":  "invalid",
			}},
			expectedParamError: globalParam + ".mode",
		},
	} {
		func() {
			defer func() {
				if t.Failed() {
					t.Logf("%#v", tc)
				}
			}()
			im, err := imageserver_native.Decode(imageserver_testdata.Medium)
			if err != nil {
				t.Fatal(err)
			}
			p := &Processor{}
			im, err = p.Process(im, tc.params)
			if err != nil {
				if err, ok := err.(*imageserver.ParamError); ok && err.Param == tc.expectedParamError {
					return
				}
				t.Fatal(err)
			}
			if tc.expectedWidth != 0 && im.Bounds().Dx() != tc.expectedWidth {
				t.Fatalf("unexpected width %d, wanted %d", im.Bounds().Dx(), tc.expectedWidth)
			}
			if tc.expectedHeight != 0 && im.Bounds().Dy() != tc.expectedHeight {
				t.Fatalf("unexpected height %d, wanted %d", im.Bounds().Dy(), tc.expectedHeight)
			}
		}()
	}
}
