package govips

import (
	_ "image/jpeg"
	"testing"
)

func Test_Blur(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Blur(vi, 5, nil)
	checkError(t, err)
	defer vi2.Free()
	if BENCHMARK_IMAGE_1_BOUNDS != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_Sharpen(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Sharpen(vi, nil)
	checkError(t, err)
	defer vi2.Free()
	if BENCHMARK_IMAGE_1_BOUNDS != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}
