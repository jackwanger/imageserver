package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkImageMarshalBinarySmall(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Small)
}

func BenchmarkImageMarshalBinaryMedium(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Medium)
}

func BenchmarkImageMarshalBinaryLarge(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Large)
}

func BenchmarkImageMarshalBinaryHuge(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Huge)
}

func BenchmarkImageMarshalBinaryAnimated(b *testing.B) {
	benchmarkImageMarshalBinary(b, testdata.Animated)
}

func benchmarkImageMarshalBinary(b *testing.B, im *Image) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := im.MarshalBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}

func BenchmarkImageUnmarshalBinarySmall(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Small)
}

func BenchmarkImageUnmarshalBinaryMedium(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Medium)
}

func BenchmarkImageUnmarshalBinaryLarge(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Large)
}

func BenchmarkImageUnmarshalBinaryHuge(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Huge)
}

func BenchmarkImageUnmarshalBinaryAnimated(b *testing.B) {
	benchmarkImageUnmarshalBinary(b, testdata.Animated)
}

func benchmarkImageUnmarshalBinary(b *testing.B, im *Image) {
	data, _ := im.MarshalBinary()
	imNew := new(Image)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := imNew.UnmarshalBinary(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}

func BenchmarkImageUnmarshalBinaryNoCopySmall(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Small)
}

func BenchmarkImageUnmarshalBinaryNoCopyMedium(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Medium)
}

func BenchmarkImageUnmarshalBinaryNoCopyLarge(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Large)
}

func BenchmarkImageUnmarshalBinaryNoCopyHuge(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Huge)
}

func BenchmarkImageUnmarshalBinaryNoCopyAnimated(b *testing.B) {
	benchmarkImageUnmarshalBinaryNoCopy(b, testdata.Animated)
}

func benchmarkImageUnmarshalBinaryNoCopy(b *testing.B, im *Image) {
	data, _ := im.MarshalBinary()
	imNew := new(Image)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := imNew.UnmarshalBinaryNoCopy(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.SetBytes(int64(len(im.Data)))
}
