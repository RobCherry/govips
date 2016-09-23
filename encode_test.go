package govips

import (
	"bytes"
	_ "golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test_EncodeGifNative(t *testing.T) {
	i := decodeNative(t, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS)
	var w bytes.Buffer
	err := gif.Encode(&w, i, nil)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(w.Bytes()), "gif", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodeJpegNative(t *testing.T) {
	i := decodeNative(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS)
	var w bytes.Buffer
	err := jpeg.Encode(&w, i, nil)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(w.Bytes()), "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodePngNative(t *testing.T) {
	i := decodeNative(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS)
	var w bytes.Buffer
	err := png.Encode(&w, i)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(w.Bytes()), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodeWebpNative(t *testing.T) {
	t.Skip("Native does not support encoding to WEBP...")
}

func Test_EncodeGifVips(t *testing.T) {
	t.Skip("Vips does not support encoding to GIF...")
}

func Test_EncodeJpegFileVips(t *testing.T) {
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	file, err := ioutil.TempFile("", "")
	checkError(t, err)
	defer os.Remove(file.Name())
	defer file.Close()
	options := EncodeJpegOptions{Q: 92}
	err = EncodeJpegFile(vi, file, &options)
	checkError(t, err)
	file.Seek(0, io.SeekStart)
	checkEncoded(t, file, "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodeJpegBytesVips(t *testing.T) {
	vi := test_DecodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	options := EncodeJpegOptions{Q: 92}
	b, err := EncodeJpegBytes(vi, &options)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(b), "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodePngFileVips(t *testing.T) {
	vi := test_DecodePngVips(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	file, err := ioutil.TempFile("", "")
	checkError(t, err)
	defer os.Remove(file.Name())
	defer file.Close()
	options := EncodePngOptions{Compression: 6} // Use: EncodePngOptions{ Compression: 4, Filter: VIPS_PNG_FILTER_UP }
	err = EncodePngFile(vi, file, &options)
	checkError(t, err)
	file.Seek(0, io.SeekStart)
	checkEncoded(t, file, "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodePngBytesVips(t *testing.T) {
	vi := test_DecodePngVips(t, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	options := EncodePngOptions{Compression: 6} // Use: EncodePngOptions{ Compression: 4, Filter: VIPS_PNG_FILTER_UP }
	b, err := EncodePngBytes(vi, &options)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(b), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodeWebpFileVips(t *testing.T) {
	vi := test_DecodeWebpVips(t, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	file, err := ioutil.TempFile("", "")
	checkError(t, err)
	defer os.Remove(file.Name())
	defer file.Close()
	options := EncodeWebpOptions{Q: 92}
	err = EncodeWebpFile(vi, file, &options)
	checkError(t, err)
	file.Seek(0, io.SeekStart)
	checkEncoded(t, file, "webp", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Test_EncodeWebpBytesVips(t *testing.T) {
	vi := test_DecodeWebpVips(t, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	options := EncodeWebpOptions{Q: 92}
	b, err := EncodeWebpBytes(vi, &options)
	checkError(t, err)
	checkEncoded(t, bytes.NewReader(b), "webp", BENCHMARK_IMAGE_1_BOUNDS.Size())
}

func Benchmark_EncodeGifNative(b *testing.B) {
	m := decodeNative(b, "benchmark_images/1.gif", BENCHMARK_IMAGE_1_BOUNDS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var w bytes.Buffer
		err := gif.Encode(&w, m, nil)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(w.Bytes()), "gif", BENCHMARK_IMAGE_1_BOUNDS.Size())
	}
}

func Benchmark_EncodeJpegNative(b *testing.B) {
	m := decodeNative(b, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var w bytes.Buffer
		err := jpeg.Encode(&w, m, &jpeg.Options{Quality: 92})
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(w.Bytes()), "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
	}
}

func Benchmark_EncodePngNative(b *testing.B) {
	m := decodeNative(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var w bytes.Buffer
		err := png.Encode(&w, m)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(w.Bytes()), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	}
}

func Benchmark_EncodeWebpNative(b *testing.B) {
	b.Skip("Native does not support encoding to WEBP...")
}

func Benchmark_EncodeGifFileVips(b *testing.B) {
	b.Skip("Vips does not support encoding to GIF...")
}

func Benchmark_EncodeGifBytesVips(b *testing.B) {
	b.Skip("Vips does not support encoding to GIF...")
}

func Benchmark_EncodeJpegFileVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodeJpegVips(b, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		file, err := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		defer file.Close()
		checkError(b, err)
		options := EncodeJpegOptions{Q: 92}
		err = EncodeJpegFile(vi, file, &options)
		checkError(b, err)
		file.Seek(0, io.SeekStart)
		checkEncoded(b, file, "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodeJpegBytesVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodeJpegVips(b, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		options := EncodeJpegOptions{Q: 92}
		buf, err := EncodeJpegBytes(vi, &options)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(buf), "jpeg", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngFileVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		file, err := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		defer file.Close()
		checkError(b, err)
		options := EncodePngOptions{Compression: 6} // Use: EncodePngOptions{ Compression: 4, Filter: VIPS_PNG_FILTER_UP }
		err = EncodePngFile(vi, file, &options)
		checkError(b, err)
		file.Seek(0, io.SeekStart)
		checkEncoded(b, file, "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngBytesVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		options := EncodePngOptions{Compression: 6} // Use: EncodePngOptions{ Compression: 4, Filter: VIPS_PNG_FILTER_UP }
		buf, err := EncodePngBytes(vi, &options)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(buf), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngFileVipsFilterUp(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		file, err := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		defer file.Close()
		checkError(b, err)
		options := EncodePngOptions{Compression: 6, Filter: VIPS_PNG_FILTER_UP} // Use: EncodePngOptions{ Compression: 4, Filter: VIPS_PNG_FILTER_UP }
		err = EncodePngFile(vi, file, &options)
		checkError(b, err)
		file.Seek(0, io.SeekStart)
		checkEncoded(b, file, "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngBytesVipsFilterUp(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		options := EncodePngOptions{Compression: 6, Filter: VIPS_PNG_FILTER_UP}
		buf, err := EncodePngBytes(vi, &options)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(buf), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngFileVipsCompress4FilterUp(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		file, err := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		defer file.Close()
		checkError(b, err)
		options := EncodePngOptions{Compression: 4, Filter: VIPS_PNG_FILTER_UP}
		err = EncodePngFile(vi, file, &options)
		checkError(b, err)
		file.Seek(0, io.SeekStart)
		checkEncoded(b, file, "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodePngBytesVipsCompress4FilterUp(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodePngVips(b, "benchmark_images/1.png", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		options := EncodePngOptions{Compression: 4, Filter: VIPS_PNG_FILTER_UP}
		buf, err := EncodePngBytes(vi, &options)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(buf), "png", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodeWebpFileVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodeWebpVips(b, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		file, err := ioutil.TempFile("", "")
		defer os.Remove(file.Name())
		defer file.Close()
		checkError(b, err)
		options := EncodeWebpOptions{Q: 92}
		err = EncodeWebpFile(vi, file, &options)
		checkError(b, err)
		file.Seek(0, io.SeekStart)
		checkEncoded(b, file, "webp", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func Benchmark_EncodeWebpBytesVips(b *testing.B) {
	benchmark_EncodeVips(b, func() *VipsImage {
		return test_DecodeWebpVips(b, "benchmark_images/1.webp", BENCHMARK_IMAGE_1_BOUNDS, nil)
	}, func(vi *VipsImage) {
		options := EncodeWebpOptions{Q: 92}
		buf, err := EncodeWebpBytes(vi, &options)
		checkError(b, err)
		checkEncoded(b, bytes.NewReader(buf), "webp", BENCHMARK_IMAGE_1_BOUNDS.Size())
	})
}

func benchmark_EncodeVips(b *testing.B, decoder func() *VipsImage, runner func(*VipsImage)) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(b)
	if err != nil {
		b.Fatal(err)
	}
	vi := decoder()
	defer vi.Free()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner(vi)
	}
}

func checkEncoded(t testing.TB, r io.Reader, format string, dimensions image.Point) {
	c, f, err := image.DecodeConfig(r)
	checkError(t, err)
	if f != format {
		t.Fatalf("Incorrectly encoded %s to %s", format, f)
	}
	if dimensions != image.Pt(c.Width, c.Height) {
		t.Fatalf("Invalid dimensions for %s: %dx%d", format, c.Width, c.Height)
	}
}
