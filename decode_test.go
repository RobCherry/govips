package govips

import (
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"testing"
)

var BENCHMARK_IMAGE_1_BOUNDS = image.Rect(0, 0, 4608, 3456)

func Test_DecodeGifNative(t *testing.T) {
	decodeNative(t, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS)
}

func Test_DecodeJpegNative(t *testing.T) {
	decodeNative(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS)
}

func Test_DecodePngNative(t *testing.T) {
	decodeNative(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS)
}

func Test_DecodeWebpNative(t *testing.T) {
	decodeNative(t, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS)
}

func Test_DecodeGifVips(t *testing.T) {
	t.Skip("GIF support using giflib/giflib5 is buggy right now...")
	options := DecodeGifOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeGifVips(t, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeJpegVips(t *testing.T) {
	options := DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeJpegVipsWithShrink(t *testing.T) {
	options := DecodeJpegOptions{Shrink: 2, DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeJpegVips(t, "benchmark_images/1.jpg", image.Rect(0, 0, BENCHMARK_IMAGE_1_BOUNDS.Dx()/2, BENCHMARK_IMAGE_1_BOUNDS.Dy()/2), &options).Free()
}

func Test_DecodePngVips(t *testing.T) {
	options := DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}
	decodePngVips(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeWebpVips(t *testing.T) {
	options := DecodeWebpOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeWebpVips(t, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeGifMagick(t *testing.T) {
	options := DecodeMagickOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeMagickVips(t, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeJpegMagick(t *testing.T) {
	options := DecodeMagickOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeMagickVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodePngMagick(t *testing.T) {
	options := DecodeMagickOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeMagickVips(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Test_DecodeWebpMagick(t *testing.T) {
	t.Skip("WEBP support using libmagick is buggy right now...")
	options := DecodeMagickOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}}
	decodeMagickVips(t, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, &options).Free()
}

func Benchmark_DecodeGifNative(b *testing.B) {
	benchmark_DecodeNative(b, "benchmark_images/1.gif")
}

func Benchmark_DecodeJpegNative(b *testing.B) {
	benchmark_DecodeNative(b, "benchmark_images/1.jpg")
}

func Benchmark_DecodePngNative(b *testing.B) {
	benchmark_DecodeNative(b, "benchmark_images/1.png")
}

func Benchmark_DecodeWebpNative(b *testing.B) {
	benchmark_DecodeNative(b, "benchmark_images/1.webp")
}

func Benchmark_DecodeGifVips(b *testing.B) {
	b.Skip("GIF support using giflib/giflib5 is buggy right now...")
	benchmark_DecodeVips(b, "benchmark_images/1.gif", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeGifReader(imageReader, nil)
	})
}

func Benchmark_DecodeJpegVips(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.jpg", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeJpegReader(imageReader, nil)
	})
}

func Benchmark_DecodePngVips(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.png", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodePngReader(imageReader, nil)
	})
}

func Benchmark_DecodeWebpVips(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.webp", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeWebpReader(imageReader, nil)
	})
}

func Benchmark_DecodeGifMagick(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.gif", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeMagickReader(imageReader, nil)
	})
}

func Benchmark_DecodeJpegMagick(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.jpg", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeMagickReader(imageReader, nil)
	})
}

func Benchmark_DecodePngMagick(b *testing.B) {
	benchmark_DecodeVips(b, "benchmark_images/1.png", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeMagickReader(imageReader, nil)
	})
}

func Benchmark_DecodeWebpMagick(b *testing.B) {
	b.Skip("WEBP support using libmagick is buggy right now...")
	benchmark_DecodeVips(b, "benchmark_images/1.webp", func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeMagickReader(imageReader, nil)
	})
}

func Benchmark_DecodeConfigGifNative(b *testing.B) {
	benchmark_DecodeConfigNative(b, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Benchmark_DecodeConfigJpegNative(b *testing.B) {
	benchmark_DecodeConfigNative(b, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Benchmark_DecodeConfigPngNative(b *testing.B) {
	benchmark_DecodeConfigNative(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Benchmark_DecodeConfigWebpNative(b *testing.B) {
	benchmark_DecodeConfigNative(b, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func benchmark_DecodeVips(b *testing.B, file string, runner func(io.Reader) (*VipsImage, error)) {
	err := Initialize()
	defer ThreadShutdown()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decodeVips(b, file, BENCHMARK_IMAGE_1_BOUNDS, runner).Free()
	}
}

func decodeVips(t testing.TB, file string, bounds image.Rectangle, runner func(io.Reader) (*VipsImage, error)) *VipsImage {
	defer checkErrorBuffer(t)
	imageReader, err := os.Open(file)
	checkError(t, err)
	vi, err := runner(imageReader)
	imageReader.Close()
	if err != nil {
		t.Fatalf("%s: %s", err, ErrorBuffer())
	}
	if bounds != vi.Bounds() {
		t.Fatalf("Invalid bounds for %s: %v", file, vi.Bounds())
	}
	return vi
}

func decodeGifVips(t testing.TB, file string, bounds image.Rectangle, options *DecodeGifOptions) *VipsImage {
	return decodeVips(t, file, bounds, func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeGifReader(imageReader, options)
	})
}

func decodeJpegVips(t testing.TB, file string, bounds image.Rectangle, options *DecodeJpegOptions) *VipsImage {
	return decodeVips(t, file, bounds, func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeJpegReader(imageReader, options)
	})
}

func decodePngVips(t testing.TB, file string, bounds image.Rectangle, options *DecodeOptions) *VipsImage {
	return decodeVips(t, file, bounds, func(imageReader io.Reader) (*VipsImage, error) {
		return DecodePngReader(imageReader, options)
	})
}

func decodeWebpVips(t testing.TB, file string, bounds image.Rectangle, options *DecodeWebpOptions) *VipsImage {
	return decodeVips(t, file, bounds, func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeWebpReader(imageReader, options)
	})
}

func decodeMagickVips(t testing.TB, file string, bounds image.Rectangle, options *DecodeMagickOptions) *VipsImage {
	return decodeVips(t, file, bounds, func(imageReader io.Reader) (*VipsImage, error) {
		return DecodeMagickReader(imageReader, options)
	})
}

func benchmark_DecodeNative(b *testing.B, file string) {
	for i := 0; i < b.N; i++ {
		decodeNative(b, file, BENCHMARK_IMAGE_1_BOUNDS)
	}
}

func decodeNative(t testing.TB, file string, bounds image.Rectangle) image.Image {
	imageReader, err := os.Open(file)
	checkError(t, err)
	m, _, err := image.Decode(imageReader)
	checkError(t, err)
	if bounds != m.Bounds() {
		t.Fatalf("Invalid bounds for %s: %v", file, m.Bounds())
	}
	imageReader.Close()
	return m
}

func benchmark_DecodeConfigNative(b *testing.B, file string, dimensions image.Point) {
	for i := 0; i < b.N; i++ {
		imageReader, err := os.Open(file)
		if err != nil {
			b.Fatal(err)
		}
		c, _, err := image.DecodeConfig(imageReader)
		if err != nil {
			b.Fatal(err)
		}
		if dimensions != image.Pt(c.Width, c.Height) {
			b.Fatalf("Invalid dimensions for %s: %dx%d", file, c.Width, c.Height)
		}
		imageReader.Close()
	}
}

func checkError(t testing.TB, err error) {
	if err != nil {
		t.Error(err)
	}
	checkErrorBuffer(t)
}

func checkErrorBuffer(t testing.TB) {
	err := ErrorBuffer()
	if err != nil {
		t.Fatal(err)
	}
}
