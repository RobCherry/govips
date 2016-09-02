package govips

import (
	"image"
	"image/color"
	"testing"
)

func Test_Embed(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
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
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
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
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Crop(vi, 100, 100, 400, 300)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, 400, 300) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_Flatten(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	red := color.NRGBA{255, 0, 0, 255}
	vi := test_DecodePngVips(t, "benchmark_images/2x1_transparent_red.png", image.Rect(0, 0, 2, 1), nil)
	defer vi.Free()

	runTest := func(expected color.Color, options *FlattenOptions) {
		o, err := Flatten(vi, options)
		checkError(t, err)
		nrgba, err := NewNRGBAVipsImage(o)
		checkError(t, err)
		if expected != *nrgba.At(0, 0).(*color.NRGBA) {
			t.Fatalf("Invalid color: %v", nrgba.At(0, 0))
		}
		if red != *nrgba.At(1, 0).(*color.NRGBA) {
			t.Fatalf("Invalid color: %v", nrgba.At(1, 0))
		}
		nrgba.Free()
	}

	runTest(color.NRGBAModel.Convert(color.Black), nil)
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: []float64{0}})
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: []float64{0, 0, 0}})
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: VIPS_BACKGROUND_BLACK})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: []float64{255}})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: []float64{255, 255, 255}})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: VIPS_BACKGROUND_WHITE})
	runTest(color.NRGBA{255, 0, 0, 255}, &FlattenOptions{Background: []float64{255, 0, 0}})
	runTest(color.NRGBA{0, 255, 0, 255}, &FlattenOptions{Background: []float64{0, 255, 0}})
	runTest(color.NRGBA{0, 0, 255, 255}, &FlattenOptions{Background: []float64{0, 0, 255}})
}
