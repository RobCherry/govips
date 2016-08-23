package govips

import (
	_ "image/jpeg"
	"testing"
)

func Test_Colourspace(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	vi := decodeJpegVips(t, "benchmark_images/1.jpg", BENCHMARK_IMAGE_1_BOUNDS, &DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}})
	defer vi.Free()
	if BENCHMARK_IMAGE_1_BOUNDS != vi.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi.Bounds())
	}
	if vi.Interpretation() != VIPS_INTERPRETATION_sRGB {
		t.Fatalf("Invalid interpretation: %v", vi.Interpretation())
	}
	if vi.Bands() != 3 {
		t.Fatalf("Invalid bands: %v", vi.Bands())
	}
	vi2, err := Colourspace(vi, VIPS_INTERPRETATION_LAB, nil)
	checkError(t, err)
	defer vi2.Free()
	if BENCHMARK_IMAGE_1_BOUNDS != vi2.Bounds() {
		t.Fatalf("Invalid bounds: %v", vi2.Bounds())
	}
	if vi2.Interpretation() != VIPS_INTERPRETATION_LAB {
		t.Fatalf("Invalid interpretation: %v", vi2.Interpretation())
	}
	if vi2.Bands() != 3 {
		t.Fatalf("Invalid bands: %v", vi2.Bands())
	}
}

func Test_ColourspaceIsSupported(t *testing.T) {
	err := Initialize()
	defer ThreadShutdown()
	defer checkErrorBuffer(t)
	checkError(t, err)
	runTest := func(path string, expectedInterpretation VipsInterpretation, expectedBands int, expectedSupport bool) {
		vi := decodeJpegVips(t, path, BENCHMARK_IMAGE_1_BOUNDS, &DecodeJpegOptions{DecodeOptions: DecodeOptions{Access: VIPS_ACCESS_SEQUENTIAL}})
		defer vi.Free()
		if BENCHMARK_IMAGE_1_BOUNDS != vi.Bounds() {
			t.Fatalf("Invalid bounds for %s: %v", path, vi.Bounds())
		}
		if vi.Interpretation() != expectedInterpretation {
			t.Fatalf("Invalid interpretation for %s: %v", path, vi.Interpretation())
		}
		if vi.Bands() != expectedBands {
			t.Fatalf("Invalid bands for %s: %v", path, vi.Bands())
		}
		if expectedSupport != ColourspaceIsSupported(vi) {
			t.Fatal("Colourspace support was incorrectly reported.")
		}
	}
	runTest("benchmark_images/1.jpg", VIPS_INTERPRETATION_sRGB, 3, true)
	runTest("benchmark_images/1_bw.jpg", VIPS_INTERPRETATION_B_W, 1, true)
	runTest("benchmark_images/1_cmyk.jpg", VIPS_INTERPRETATION_CMYK, 4, false)
}
