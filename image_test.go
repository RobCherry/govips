package govips

import (
	"image/color"
	"testing"
)

func Test_NRGBAVipsImage(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, &DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}})
	checkError(t, err)
	defer vi.Free()
	nrgba, err := NewNRGBAVipsImage(vi)
	checkError(t, err)
	defer nrgba.Free()
	if nrgba.ColorModel() != color.NRGBAModel {
		t.Fatal("Invalid color model")
	}
	if *nrgba.At(0, 0).(*color.NRGBA) != (color.NRGBA{249, 249, 249, 255}) {
		t.Fatalf("Invalid color: %v", nrgba.At(0, 0))
	}
}

func Test_CMYKVipsImage(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := test_DecodeJpegVips(t, "benchmark_images/1_cmyk.jpg", BENCHMARK_IMAGE_1_BOUNDS, &DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}})
	checkError(t, err)
	defer vi.Free()
	cmyk, err := NewCMYKVipsImage(vi)
	checkError(t, err)
	defer cmyk.Free()
	if cmyk.ColorModel() != color.CMYKModel {
		t.Fatal("Invalid color model")
	}
	if *cmyk.At(0, 0).(*color.CMYK) != (color.CMYK{0, 0, 0, 6}) {
		t.Fatalf("Invalid color: %v", cmyk.At(0, 0))
	}
}

func Test_GrayVipsImage(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := test_DecodeJpegVips(t, "benchmark_images/1_bw.jpg", BENCHMARK_IMAGE_1_BOUNDS, &DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}})
	checkError(t, err)
	defer vi.Free()
	gray, err := NewGrayVipsImage(vi)
	checkError(t, err)
	defer gray.Free()
	if gray.ColorModel() != color.GrayModel {
		t.Fatal("Invalid color model")
	}
	if *gray.At(0, 0).(*color.Gray) != (color.Gray{249}) {
		t.Fatalf("Invalid color: %v", gray.At(0, 0))
	}
}
