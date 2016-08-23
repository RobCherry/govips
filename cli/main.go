package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/RobCherry/govips"
	"golang.org/x/image/draw"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var scalerByName = map[string]draw.Scaler{
	"NearestNeighbor": draw.NearestNeighbor,
	"ApproxBiLinear":  draw.ApproxBiLinear,
	"BiLinear":        draw.BiLinear,
	"CatmullRom":      draw.CatmullRom,
}

var (
	fileLocation   string
	outputLocation string

	resize        string
	resizePoint   image.Point
	fastResize    bool
	crop          string
	cropRectangle image.Rectangle
	quality       int
	blur          uint
	vips          bool

	scalerName = "ApproxBiLinear"

	err error
)

func init() {
	var b bytes.Buffer
	for name := range scalerByName {
		if b.Len() > 0 {
			b.WriteString(", ")
		}
		b.WriteString(name)
	}
	availableScalers := b.String()

	flag.StringVar(&resize, "r", "", "Resize.  Example: 300x300")
	flag.BoolVar(&fastResize, "fast-resize", false, "Do a nearest neighbor partial resize. (Native)")
	flag.StringVar(&crop, "c", "", "Crop (before resize).  Example: 0x50y400w300h")
	flag.IntVar(&quality, "q", 85, "Quality (1-100)")
	flag.UintVar(&blur, "b", 0, "Blur")
	flag.StringVar(&scalerName, "s", scalerName, "Scaler.  One of: "+availableScalers)
	flag.BoolVar(&vips, "v", false, "VIPS")

	flag.Parse()

	fileLocation = flag.Arg(0)
	if len(fileLocation) == 0 {
		fmt.Println("Please provide a file or URL")
		os.Exit(1)
	}

	outputLocation = flag.Arg(1)
	if len(outputLocation) == 0 {
		fmt.Println("Please provide an output file location")
		os.Exit(1)
	}

	if len(resize) > 0 {
		parts := strings.Split(resize, "x")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("Invalid resize: %s", resize))
		}
		resizePoint.X, err = strconv.Atoi(parts[0])
		checkErr(err)
		resizePoint.Y, err = strconv.Atoi(parts[1])
		checkErr(err)
	}

	if len(crop) > 0 {
		parts := strings.Split(crop, "x")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("Invalid crop: %s", crop))
		}
		cropRectangle.Min.X, err = strconv.Atoi(parts[0])
		checkErr(err)

		parts = strings.Split(parts[1], "y")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("Invalid crop: %s", crop))
		}
		cropRectangle.Min.Y, err = strconv.Atoi(parts[0])
		checkErr(err)

		parts = strings.Split(parts[1], "w")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("Invalid crop: %s", crop))
		}
		cropRectangle.Max.X, err = strconv.Atoi(parts[0])
		checkErr(err)

		parts = strings.Split(parts[1], "h")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("Invalid crop: %s", crop))
		}
		cropRectangle.Max.Y, err = strconv.Atoi(parts[0])
		checkErr(err)
	}
}

func main() {
	output, err := os.Create(outputLocation)
	if os.IsExist(err) {
		fmt.Printf("%s already exists\n", outputLocation)
		os.Exit(1)
	} else {
		checkErr(err)
	}
	defer output.Close()

	if strings.HasPrefix(fileLocation, "http") {
		response, err := http.Get(fileLocation)
		checkErr(err)
		if response.StatusCode != 200 {
			checkErr(fmt.Errorf("%s not found\n", fileLocation))
		}
		fileLocation = outputLocation
		ext := filepath.Ext(fileLocation)
		if len(ext) > 0 {
			fileLocation = fileLocation[0 : len(fileLocation)-len(ext)]
		}
		fileLocation = fileLocation + "_original" + ext

		file, err := os.Create(fileLocation)
		checkErr(err)
		_, err = io.Copy(file, response.Body)
		checkErr(err)
		file.Close()
	}

	var decodeDuration, resizeDuration, cropDuration, blurDuration, encodeDuration, totalDuration time.Duration
	startTime := time.Now()

	imageReader, err := os.Open(fileLocation)
	checkErr(err)
	defer imageReader.Close()

	if vips {
		i, format, err := decodeVips(imageReader)
		checkErr(err)
		decodeDuration = time.Since(startTime)
		if cropRectangle != image.ZR {
			localStartTime := time.Now()
			i2, err := cropVips(i, cropRectangle.Min.X, cropRectangle.Min.Y, cropRectangle.Max.X, cropRectangle.Max.Y)
			checkErr(err)
			i.Free()
			i = i2
			cropDuration = time.Since(localStartTime)
		}
		if resizePoint != image.ZP {
			localStartTime := time.Now()
			i, err = resizeVips(i, resizePoint.X, resizePoint.Y, fastResize)
			checkErr(err)
			resizeDuration = time.Since(localStartTime)
		}
		if blur > 0 {
			localStartTime := time.Now()
			i2, err := govips.Blur(i, float64(blur), nil)
			checkErr(err)
			i.Free()
			i = i2
			blurDuration = time.Since(localStartTime)
		}
		localStartTime := time.Now()
		switch format {
		case "jpeg":
			_, err = govips.EncodeJpegWriter(output, i, &govips.EncodeJpegOptions{Q: quality, OptimizeCoding: true, Strip: true, NoSubsample: true})
		case "gif", "png":
			_, err = govips.EncodePngWriter(output, i, &govips.EncodePngOptions{Compression: 6})
		case "webp":
			_, err = govips.EncodeWebpWriter(output, i, &govips.EncodeWebpOptions{Q: quality})
		default:
			err = fmt.Errorf("Invalid image format: %s\n", format)
		}
		encodeDuration = time.Since(localStartTime)
		checkErr(err)
		i.Free()
	} else {
		i, format, err := image.Decode(imageReader)
		checkErr(err)
		decodeDuration = time.Since(startTime)
		if cropRectangle != image.ZR {
			localStartTime := time.Now()
			i = cropNative(i, cropRectangle.Min.X, cropRectangle.Min.Y, cropRectangle.Max.X, cropRectangle.Max.Y)
			cropDuration = time.Since(localStartTime)
		}
		if resizePoint != image.ZP {
			localStartTime := time.Now()
			i = resizeNative(i, resizePoint.X, resizePoint.Y, scalerByName[scalerName], fastResize)
			resizeDuration = time.Since(localStartTime)
		}
		if blur > 0 {
			fmt.Println("Native blur not supported...")
		}
		localStartTime := time.Now()
		switch format {
		case "gif":
			err = gif.Encode(output, i, &gif.Options{256, nil, nil})
		case "jpeg":
			err = jpeg.Encode(output, i, &jpeg.Options{quality})
		case "png":
			err = png.Encode(output, i)
		default:
			err = fmt.Errorf("Invalid image format: %s\n", format)
		}
		encodeDuration = time.Since(localStartTime)
		checkErr(err)
	}
	totalDuration = time.Since(startTime)

	fmt.Println("Timing Data:")
	fmt.Printf("  Decode: %v\n", decodeDuration)
	fmt.Printf("  Crop: %v\n", cropDuration)
	fmt.Printf("  Resize: %v\n", resizeDuration)
	fmt.Printf("  Blur: %v\n", blurDuration)
	fmt.Printf("  Encode: %v\n", encodeDuration)
	fmt.Printf("  Total: %v\n", totalDuration)
}

func decodeVips(reader io.ReadSeeker) (*govips.VipsImage, string, error) {
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, format, err
	}
	reader.Seek(0, 0)
	switch format {
	case "gif":
		vipsImage, err := govips.DecodeMagickReader(reader, nil)
		return vipsImage, format, err
	case "jpeg":
		vipsImage, err := govips.DecodeJpegReader(reader, nil)
		return vipsImage, format, err
	case "png":
		vipsImage, err := govips.DecodePngReader(reader, nil)
		return vipsImage, format, err
	case "webp":
		vipsImage, err := govips.DecodeWebpReader(reader, nil)
		return vipsImage, format, err
	}
	return nil, format, fmt.Errorf("Unhandled format: %s", format)
}

func resizeVips(i *govips.VipsImage, width, height int, useFastScale bool) (*govips.VipsImage, error) {
	scale := math.Min(float64(width)/float64(i.Bounds().Dx()), float64(height)/float64(i.Bounds().Dy()))
	if scale < 1 {
		//i2, err := govips.Resize(i, scale, scale, govips.VIPS_KERNEL_LANCZOS3)
		//checkErr(err)
		//i.Free()
		//i = i2
		//return i, nil

		// Perform a fast shrink ...
		if useFastScale {
			shrink := math.Max(1, math.Floor(1/(scale*2)))
			if shrink > 1 {
				i2, err := govips.Shrink(i, shrink, shrink)
				checkErr(err)
				i.Free()
				i = i2
			}
		}
		// Recompute scale...
		scale = math.Min(float64(width)/float64(i.Bounds().Dx()), float64(height)/float64(i.Bounds().Dy()))
		if scale < 1 {
			i2, err := govips.Reduce(i, 1/scale, 1/scale, govips.VIPS_KERNEL_LANCZOS3)
			checkErr(err)
			i.Free()
			i = i2
		}
	}
	return i, nil
}

func cropVips(i *govips.VipsImage, x, y, width, height int) (*govips.VipsImage, error) {
	return govips.Crop(i, x, y, width, height)
}

func resizeNative(i image.Image, width, height int, scaler draw.Scaler, useFastScale bool) image.Image {
	if scaler == nil {
		scaler = draw.BiLinear
	}
	imageBounds := i.Bounds()
	imageWidth := float64(imageBounds.Dx())
	imageHeight := float64(imageBounds.Dy())
	scaleFactor := math.Min(float64(width)/imageWidth, float64(height)/imageHeight)

	if scaleFactor >= 1.0 {
		return i
	}

	scaledWidth := int(roundFloat64(imageWidth * scaleFactor))
	scaledHeight := int(roundFloat64(imageHeight * scaleFactor))

	if useFastScale {
		nearestNeighborScale := math.Max(1.0, math.Floor(1.0/(scaleFactor*2)))
		//fmt.Printf("Fast Scale Factor: %v\n", nearestNeighborScale)
		if scaler != draw.NearestNeighbor && nearestNeighborScale > 1.0 {
			r := image.NewNRGBA(image.Rect(0, 0, roundInt(imageWidth/nearestNeighborScale), roundInt(imageHeight/nearestNeighborScale)))
			//fmt.Printf("Fast Scale Dimensions: %vx%v\n", r.Bounds().Dx(), r.Bounds().Dy())
			draw.NearestNeighbor.Scale(r, r.Bounds(), i, i.Bounds(), draw.Src, nil)
			i = r
		}
	}

	r := image.NewNRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))
	scaler.Scale(r, r.Bounds(), i, i.Bounds(), draw.Src, nil)
	return r
}

func cropNative(i image.Image, x, y, width, height int) image.Image {
	return i.(SubImager).SubImage(image.Rect(x, y, x+width, y+height))
}

// SubImager is a utility interface for an image.Image that can extract a sub-image.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func roundInt(value float64) int {
	return int(roundFloat64(value))
}

func roundFloat64(value float64) float64 {
	if value < 0.0 {
		return value - 0.5
	}
	return value + 0.5
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		if vips {
			vipsError := govips.ErrorBuffer()
			if vipsError != nil {
				fmt.Println(vipsError)
			}
		}
		os.Exit(1)
	}
}
