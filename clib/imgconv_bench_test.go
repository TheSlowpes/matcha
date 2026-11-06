package clib

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

// BenchmarkDecodeToPNG converts a newsletter-sized JPEG to PNG.
func BenchmarkDecodeToPNG(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 600))
	for y := range 600 {
		for x := range 1200 {
			img.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(x ^ y), 255})
		}
	}
	var src bytes.Buffer
	if err := jpeg.Encode(&src, img, &jpeg.Options{Quality: 85}); err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if _, ok := DecodeToPNG(src.Bytes()); !ok {
			b.Fatal("decode failed")
		}
	}
}
