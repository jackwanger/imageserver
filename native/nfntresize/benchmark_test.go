package nfntresize

import (
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_native "github.com/pierrre/imageserver/native"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkResizeSmall(b *testing.B) {
	benchmarkResize(b, testdata.Small)
}

func BenchmarkResizeMedium(b *testing.B) {
	benchmarkResize(b, testdata.Medium)
}

func BenchmarkResizeLarge(b *testing.B) {
	benchmarkResize(b, testdata.Large)
}

func BenchmarkResizeHuge(b *testing.B) {
	benchmarkResize(b, testdata.Huge)
}

func benchmarkResize(b *testing.B, im *imageserver.Image) {
	nim, err := imageserver_native.Decode(im)
	if err != nil {
		b.Fatal(err)
	}
	proc := &Processor{}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width": 100,
		},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := proc.Process(nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}
