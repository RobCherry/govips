package govips

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
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
	b := createTestPNG(t, image.Pt(2, 1), []color.Color{color.NRGBA{0, 0, 0, 0}, red})
	vi, err := DecodePngBytes(b, nil)
	checkError(t, err)
	defer vi.Free()

	runTest := func(expected color.Color, options *FlattenOptions) {
		o, err := Flatten(vi, options)
		checkError(t, err)

		b2, err := EncodePngBytes(o, nil)
		checkError(t, err)
		m, _, err := image.Decode(bytes.NewBuffer(*b2))
		checkError(t, err)
		if expected != color.NRGBAModel.Convert(m.At(0, 0)) {
			t.Fatalf("Invalid color: %v", m.At(0, 0))
		}
		if red != color.NRGBAModel.Convert(m.At(1, 0)) {
			t.Fatalf("Invalid color: %v", m.At(1, 0))
		}
	}
	runTest(color.NRGBAModel.Convert(color.Black), nil)
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: &[]float64{0}})
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: &[]float64{0, 0, 0}})
	runTest(color.NRGBAModel.Convert(color.Black), &FlattenOptions{Background: VIPS_BACKGROUND_BLACK})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: &[]float64{255}})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: &[]float64{255, 255, 255}})
	runTest(color.NRGBAModel.Convert(color.White), &FlattenOptions{Background: VIPS_BACKGROUND_WHITE})
	runTest(color.NRGBA{255, 0, 0, 255}, &FlattenOptions{Background: &[]float64{255, 0, 0}})
	runTest(color.NRGBA{0, 255, 0, 255}, &FlattenOptions{Background: &[]float64{0, 255, 0}})
	runTest(color.NRGBA{0, 0, 255, 255}, &FlattenOptions{Background: &[]float64{0, 0, 255}})
}

func createTestPNG(t testing.TB, size image.Point, pixels []color.Color) []byte {
	if len(pixels) != size.X*size.Y {
		t.Fatalf("Invalid pixel count: %d", len(pixels))
	}
	m := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			m.Set(x, y, pixels[x+(y*size.Y)])
		}
	}
	var b bytes.Buffer
	checkError(t, png.Encode(&b, m))
	return b.Bytes()
}
