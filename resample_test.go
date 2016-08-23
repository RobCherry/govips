package govips

import (
	"image"
	_ "image/jpeg"
	"testing"
)

func Test_Affine(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	tests := map[float64]image.Point{
		0.25: BENCHMARK_IMAGE_1_BOUNDS.Size().Div(4),
		0.5:  BENCHMARK_IMAGE_1_BOUNDS.Size().Div(2),
		1:    BENCHMARK_IMAGE_1_BOUNDS.Size(),
		2:    BENCHMARK_IMAGE_1_BOUNDS.Size().Mul(2),
	}
	for scale, expected := range tests {
		vi2, err := Affine(vi, scale, 0, 0, scale, nil)
		checkError(t, err)
		defer vi2.Free()
		if expected != vi2.Bounds().Size() {
			t.Fatalf("Invalid bounds for scale with factor of %0.2f: %v", scale, vi2.Bounds())
		}
	}
}

func Test_Similarity(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	tests := map[float64]image.Point{
		0.25: BENCHMARK_IMAGE_1_BOUNDS.Size().Div(4),
		0.5:  BENCHMARK_IMAGE_1_BOUNDS.Size().Div(2),
		1:    BENCHMARK_IMAGE_1_BOUNDS.Size(),
		2:    BENCHMARK_IMAGE_1_BOUNDS.Size().Mul(2),
	}
	for scale, expected := range tests {
		vi2, err := Similarity(vi, &SimilarityOptions{Scale: scale})
		checkError(t, err)
		defer vi2.Free()
		if expected != vi2.Bounds().Size() {
			t.Fatalf("Invalid bounds for scale with factor of %0.2f: %v", scale, vi2.Bounds())
		}
	}
}

func Test_SimilarityRotate(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	tests := map[float64]image.Rectangle{
		0:   BENCHMARK_IMAGE_1_BOUNDS,
		90:  image.Rect(0, 0, BENCHMARK_IMAGE_1_BOUNDS.Dy(), BENCHMARK_IMAGE_1_BOUNDS.Dx()),
		180: BENCHMARK_IMAGE_1_BOUNDS,
		270: image.Rect(0, 0, BENCHMARK_IMAGE_1_BOUNDS.Dy(), BENCHMARK_IMAGE_1_BOUNDS.Dx()),
		360: BENCHMARK_IMAGE_1_BOUNDS,
	}
	for angle, expected := range tests {
		vi2, err := Similarity(vi, &SimilarityOptions{Angle: angle})
		checkError(t, err)
		defer vi2.Free()
		if expected != vi2.Bounds() {
			t.Fatalf("Invalid bounds for %0.0f degree rotation: %v", angle, vi2.Bounds())
		}
	}
}

func Test_ReduceVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Reduce(vi, 2, 2, VIPS_KERNEL_LANCZOS3)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx()/2, vi.Bounds().Dy()/2) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ReduceHVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := ReduceH(vi, 2, VIPS_KERNEL_LANCZOS3)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx()/2, vi.Bounds().Dy()) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ReduceVVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := ReduceV(vi, 2, VIPS_KERNEL_LANCZOS3)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx(), vi.Bounds().Dy()/2) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ShrinkVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := Shrink(vi, 2, 2)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx()/2, vi.Bounds().Dy()/2) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ShrinkHVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := ShrinkH(vi, 2)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx()/2, vi.Bounds().Dy()) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}

func Test_ShrinkVVips(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, nil)
	defer vi.Free()
	vi2, err := ShrinkV(vi, 2)
	checkError(t, err)
	defer vi2.Free()
	if image.Rect(0, 0, vi.Bounds().Dx(), vi.Bounds().Dy()/2) != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
}
