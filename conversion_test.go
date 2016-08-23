package govips

import (
	"image"
	_ "image/jpeg"
	"testing"
)

func Test_Embed(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	width := 5000
	height := 5000
	x := (width - BENCHMARK_IMAGE_1_BOUNDS.Dx()) / 2
	y := (height - BENCHMARK_IMAGE_1_BOUNDS.Dy()) / 2
	vi2, err := Embed(vi, x, y, width, height, nil)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, width, height) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ExtractArea(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := ExtractArea(vi, 100, 100, 400, 300)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, 400, 300) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_Crop(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Crop(vi, 100, 100, 400, 300)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, 400, 300) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_Flatten(t *testing.T) {
	t.Skip()
}
