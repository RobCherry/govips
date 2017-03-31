package govips

/*
#cgo pkg-config: vips
#include "vips.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Special constants used to signify a zero value instead of the default value.
const (
	INT_ZERO    = -1
	FLOAT_ZERO  = -1.0
	STRING_ZERO = "GOVIPS_STRING_ZERO"

	DEFAULT_CONCURRENCY      = 0
	DEFAULT_CACHE_MAX        = 1000
	DEFAULT_CACHE_MAX_FILES  = 100
	DEFAULT_CACHE_MAX_MEMORY = 100 * 1024 * 1024
)

var (
	ErrInitialize = errors.New("Failed to initialize libvips")
	ErrConfigure  = errors.New("Failed to configure libvips")

	ErrLoad = errors.New("Failed to load image")
	ErrSave = errors.New("Failed to save image")

	ErrEmbed        = errors.New("Failed to embed image")
	ErrCrop         = errors.New("Failed to crop image")
	ErrShrink       = errors.New("Failed to shrink image")
	ErrReduce       = errors.New("Failed to reduce image")
	ErrResize       = errors.New("Failed to resize image")
	ErrAffine       = errors.New("Failed to affine image")
	ErrBlur         = errors.New("Failed to blur image")
	ErrSharpen      = errors.New("Failed to sharpen image")
	ErrFlatten      = errors.New("Failed to flatten image")
	ErrColourspace  = errors.New("Failed to convert colourspace of image")
	ErrICCTransform = errors.New("Failed to transform colourspace of image")
)

var (
	VIPS_BACKGROUND_BLACK = []float64{0}
	VIPS_BACKGROUND_WHITE = []float64{255}
)

type Config struct {
	Concurrency    int
	CacheMax       int
	CacheMaxFiles  int
	CacheMaxMemory int
}

var (
	initializeLock sync.Mutex
	initialized    uint32
)

func Initialize() error {
	initializeLock.Lock()
	defer initializeLock.Unlock()
	if atomic.LoadUint32(&initialized) == 0 {
		if err := C.vips_init(cApplicationNane); err != 0 {
			C.vips_shutdown()
			return ErrInitialize
		}
		atomic.StoreUint32(&initialized, 1)
	}
	return nil
}

func InitializeWithConfig(config Config) error {
	err := Initialize()
	if err != nil {
		return err
	}
	return Configure(config)
}

func Configure(config Config) error {
	initializeLock.Lock()
	defer initializeLock.Unlock()
	if atomic.LoadUint32(&initialized) == 0 {
		return ErrConfigure
	}
	if config.Concurrency > 0 {
		C.vips_concurrency_set(C.int(config.Concurrency))
	} else {
		C.vips_concurrency_set(C.int(DEFAULT_CONCURRENCY))
	}
	if config.CacheMax > 0 {
		C.vips_cache_set_max(C.int(config.CacheMax))
	} else {
		C.vips_cache_set_max(C.int(DEFAULT_CACHE_MAX))
	}
	if config.CacheMaxFiles > 0 {
		C.vips_cache_set_max_files(C.int(config.CacheMaxFiles))
	} else {
		C.vips_cache_set_max_files(C.int(DEFAULT_CACHE_MAX_FILES))
	}
	if config.CacheMaxMemory > 0 {
		C.vips_cache_set_max_mem(C.size_t(config.CacheMaxMemory))
	} else {
		C.vips_cache_set_max_mem(C.size_t(DEFAULT_CACHE_MAX_MEMORY))
	}
	return nil
}

func Shutdown() {
	initializeLock.Lock()
	defer initializeLock.Unlock()
	if atomic.LoadUint32(&initialized) == 1 {
		C.vips_shutdown()
		atomic.StoreUint32(&initialized, 0)
	}
}

func ThreadShutdown() {
	C.vips_thread_shutdown()
}

func ErrorBuffer() error {
	C.vips_error_freeze()
	defer C.vips_error_thaw()
	errorBuffer := C.GoString(C.vips_error_buffer())
	if len(errorBuffer) == 0 {
		return nil
	}
	C.vips_error_clear()
	return errors.New(errorBuffer)
}

// Constants

var (
	cApplicationNane    = C.CString("govips")
	cVIPS_META_ICC_NAME = C.CString(C.VIPS_META_ICC_NAME)
)

type VipsInterpretation int

func (i VipsInterpretation) toC() C.VipsInterpretation {
	return C.VipsInterpretation(i)
}

const (
	VIPS_INTERPRETATION_ERROR     VipsInterpretation = C.VIPS_INTERPRETATION_ERROR
	VIPS_INTERPRETATION_MULTIBAND VipsInterpretation = C.VIPS_INTERPRETATION_MULTIBAND
	VIPS_INTERPRETATION_B_W       VipsInterpretation = C.VIPS_INTERPRETATION_B_W
	VIPS_INTERPRETATION_HISTOGRAM VipsInterpretation = C.VIPS_INTERPRETATION_HISTOGRAM
	VIPS_INTERPRETATION_XYZ       VipsInterpretation = C.VIPS_INTERPRETATION_XYZ
	VIPS_INTERPRETATION_LAB       VipsInterpretation = C.VIPS_INTERPRETATION_LAB
	VIPS_INTERPRETATION_CMYK      VipsInterpretation = C.VIPS_INTERPRETATION_CMYK
	VIPS_INTERPRETATION_LABQ      VipsInterpretation = C.VIPS_INTERPRETATION_LABQ
	VIPS_INTERPRETATION_RGB       VipsInterpretation = C.VIPS_INTERPRETATION_RGB
	VIPS_INTERPRETATION_CMC       VipsInterpretation = C.VIPS_INTERPRETATION_CMC
	VIPS_INTERPRETATION_LCH       VipsInterpretation = C.VIPS_INTERPRETATION_LCH
	VIPS_INTERPRETATION_LABS      VipsInterpretation = C.VIPS_INTERPRETATION_LABS
	VIPS_INTERPRETATION_sRGB      VipsInterpretation = C.VIPS_INTERPRETATION_sRGB
	VIPS_INTERPRETATION_YXY       VipsInterpretation = C.VIPS_INTERPRETATION_YXY
	VIPS_INTERPRETATION_FOURIER   VipsInterpretation = C.VIPS_INTERPRETATION_FOURIER
	VIPS_INTERPRETATION_RGB16     VipsInterpretation = C.VIPS_INTERPRETATION_RGB16
	VIPS_INTERPRETATION_GREY16    VipsInterpretation = C.VIPS_INTERPRETATION_GREY16
	VIPS_INTERPRETATION_MATRIX    VipsInterpretation = C.VIPS_INTERPRETATION_MATRIX
	VIPS_INTERPRETATION_scRGB     VipsInterpretation = C.VIPS_INTERPRETATION_scRGB
	VIPS_INTERPRETATION_HSV       VipsInterpretation = C.VIPS_INTERPRETATION_HSV
	VIPS_INTERPRETATION_LAST      VipsInterpretation = C.VIPS_INTERPRETATION_LAST
)

type VipsAccess int

func (a VipsAccess) toC() C.VipsAccess {
	switch a {
	case VIPS_ACCESS_RANDOM:
		return C.VIPS_ACCESS_RANDOM
	case VIPS_ACCESS_SEQUENTIAL:
		return C.VIPS_ACCESS_SEQUENTIAL
	case VIPS_ACCESS_SEQUENTIAL_UNBUFFERED:
		return C.VIPS_ACCESS_SEQUENTIAL_UNBUFFERED
	case VIPS_ACCESS_LAST:
		return C.VIPS_ACCESS_LAST
	default:
		return C.VIPS_ACCESS_RANDOM
	}
}

const (
	VIPS_ACCESS_RANDOM VipsAccess = iota
	VIPS_ACCESS_SEQUENTIAL
	VIPS_ACCESS_SEQUENTIAL_UNBUFFERED
	VIPS_ACCESS_LAST
)

const (
	JPEG_QUANTIZATION_TABLE_DEFAULT int = 0
	JPEG_QUANTIZATION_TABLE_FLAT
	JPEG_QUANTIZATION_TABLE_MSSIM
	JPEG_QUANTIZATION_TABLE_IMAGEMAGICK
	JPEG_QUANTIZATION_TABLE_PSNR_HVS_M
)

type PngFilter int

func (p PngFilter) toC() C.VipsForeignPngFilter {
	switch p {
	case VIPS_PNG_FILTER_ALL:
		return C.VIPS_FOREIGN_PNG_FILTER_ALL
	case VIPS_PNG_FILTER_NONE:
		return C.VIPS_FOREIGN_PNG_FILTER_NONE
	case VIPS_PNG_FILTER_SUB:
		return C.VIPS_FOREIGN_PNG_FILTER_SUB
	case VIPS_PNG_FILTER_UP:
		return C.VIPS_FOREIGN_PNG_FILTER_UP
	case VIPS_PNG_FILTER_AVG:
		return C.VIPS_FOREIGN_PNG_FILTER_AVG
	case VIPS_PNG_FILTER_PAETH:
		return C.VIPS_FOREIGN_PNG_FILTER_PAETH
	default:
		return C.VIPS_FOREIGN_PNG_FILTER_ALL
	}
}

const (
	VIPS_PNG_FILTER_DEFAULT PngFilter = iota
	VIPS_PNG_FILTER_NONE
	VIPS_PNG_FILTER_SUB
	VIPS_PNG_FILTER_UP
	VIPS_PNG_FILTER_AVG
	VIPS_PNG_FILTER_PAETH
	VIPS_PNG_FILTER_ALL
)

type WebpPreset int

func (p WebpPreset) toC() C.VipsForeignWebpPreset {
	switch p {
	case VIPS_WEBP_PRESET_DEFAULT:
		return C.VIPS_FOREIGN_WEBP_PRESET_DEFAULT
	case VIPS_WEBP_PRESET_PICTURE:
		return C.VIPS_FOREIGN_WEBP_PRESET_PICTURE
	case VIPS_WEBP_PRESET_PHOTO:
		return C.VIPS_FOREIGN_WEBP_PRESET_PHOTO
	case VIPS_WEBP_PRESET_DRAWING:
		return C.VIPS_FOREIGN_WEBP_PRESET_DRAWING
	case VIPS_WEBP_PRESET_ICON:
		return C.VIPS_FOREIGN_WEBP_PRESET_ICON
	case VIPS_WEBP_PRESET_TEXT:
		return C.VIPS_FOREIGN_WEBP_PRESET_TEXT
	default:
		return C.VIPS_FOREIGN_WEBP_PRESET_LAST
	}
}

const (
	VIPS_WEBP_PRESET_DEFAULT WebpPreset = iota
	VIPS_WEBP_PRESET_PICTURE
	VIPS_WEBP_PRESET_PHOTO
	VIPS_WEBP_PRESET_DRAWING
	VIPS_WEBP_PRESET_ICON
	VIPS_WEBP_PRESET_TEXT
	VIPS_WEBP_PRESET_LAST
)

type VipsExtend int

func (e VipsExtend) toC() C.VipsExtend {
	switch e {
	case VIPS_EXTEND_BLACK:
		return C.VIPS_EXTEND_BLACK
	case VIPS_EXTEND_COPY:
		return C.VIPS_EXTEND_COPY
	case VIPS_EXTEND_REPEAT:
		return C.VIPS_EXTEND_REPEAT
	case VIPS_EXTEND_MIRROR:
		return C.VIPS_EXTEND_MIRROR
	case VIPS_EXTEND_WHITE:
		return C.VIPS_EXTEND_WHITE
	case VIPS_EXTEND_BACKGROUND:
		return C.VIPS_EXTEND_BACKGROUND
	default:
		return C.VIPS_EXTEND_BLACK
	}
}

const (
	VIPS_EXTEND_BLACK VipsExtend = iota
	VIPS_EXTEND_COPY
	VIPS_EXTEND_REPEAT
	VIPS_EXTEND_MIRROR
	VIPS_EXTEND_WHITE
	VIPS_EXTEND_BACKGROUND
)

type VipsKernel int

func (k VipsKernel) toC() C.VipsKernel {
	switch k {
	case VIPS_KERNEL_NEAREST:
		return C.VIPS_KERNEL_NEAREST
	case VIPS_KERNEL_LINEAR:
		return C.VIPS_KERNEL_LINEAR
	case VIPS_KERNEL_CUBIC:
		return C.VIPS_KERNEL_CUBIC
	case VIPS_KERNEL_LANCZOS2:
		return C.VIPS_KERNEL_LANCZOS2
	case VIPS_KERNEL_LANCZOS3:
		return C.VIPS_KERNEL_LANCZOS3
	case VIPS_KERNEL_LAST:
		return C.VIPS_KERNEL_LAST
	default:
		return C.VIPS_KERNEL_LANCZOS3
	}
}

const (
	VIPS_KERNEL_NEAREST VipsKernel = iota
	VIPS_KERNEL_LINEAR
	VIPS_KERNEL_CUBIC
	VIPS_KERNEL_LANCZOS2
	VIPS_KERNEL_LANCZOS3
	VIPS_KERNEL_LAST
)

type VipsPrecision int

func (p VipsPrecision) toC() C.VipsPrecision {
	switch p {
	case VIPS_PRECISION_INTEGER:
		return C.VIPS_PRECISION_INTEGER
	case VIPS_PRECISION_FLOAT:
		return C.VIPS_PRECISION_FLOAT
	case VIPS_PRECISION_APPROXIMATE:
		return C.VIPS_PRECISION_APPROXIMATE
	case VIPS_PRECISION_LAST:
		return C.VIPS_PRECISION_LAST
	default:
		return C.VIPS_PRECISION_INTEGER
	}
}

const (
	VIPS_PRECISION_INTEGER VipsPrecision = iota
	VIPS_PRECISION_FLOAT
	VIPS_PRECISION_APPROXIMATE
	VIPS_PRECISION_LAST
)

type VipsIntent int

func (i VipsIntent) toC() C.VipsIntent {
	switch i {
	case VIPS_INTENT_PERCEPTUAL:
		return C.VIPS_INTENT_PERCEPTUAL
	case VIPS_INTENT_RELATIVE:
		return C.VIPS_INTENT_RELATIVE
	case VIPS_INTENT_SATURATION:
		return C.VIPS_INTENT_SATURATION
	case VIPS_INTENT_ABSOLUTE:
		return C.VIPS_INTENT_ABSOLUTE
	case VIPS_INTENT_LAST:
		return C.VIPS_INTENT_LAST
	default:
		return C.VIPS_INTENT_PERCEPTUAL
	}
}

const (
	VIPS_INTENT_PERCEPTUAL = iota
	VIPS_INTENT_RELATIVE
	VIPS_INTENT_SATURATION
	VIPS_INTENT_ABSOLUTE
	VIPS_INTENT_LAST
)

// Image

type VipsImage struct {
	cVipsImage *C.struct__VipsImage
	goBytes    []byte
}

func (v *VipsImage) Bounds() image.Rectangle {
	if v.cVipsImage == nil {
		return image.ZR
	}
	return image.Rect(0, 0, int(C.vips_image_get_width(v.cVipsImage)), int(C.vips_image_get_height(v.cVipsImage)))
}

func (v *VipsImage) Interpretation() VipsInterpretation {
	if v.cVipsImage == nil {
		return VIPS_INTERPRETATION_ERROR
	}
	return VipsInterpretation(C.vips_image_guess_interpretation(v.cVipsImage))
}

func (v *VipsImage) Bands() int {
	if v.cVipsImage == nil {
		return 0
	}
	return int(v.cVipsImage.Bands)
}

func (v *VipsImage) HasProfile() bool {
	if v.cVipsImage == nil {
		return false
	}
	return C.vips_image_get_typeof(v.cVipsImage, cVIPS_META_ICC_NAME) != 0
}

func (v *VipsImage) RemoveProfile() {
	C.vips_image_remove(v.cVipsImage, cVIPS_META_ICC_NAME)
}

func (v *VipsImage) Free() {
	if v.cVipsImage != nil {
		C.g_object_unref(C.gpointer(v.cVipsImage))
		v.cVipsImage = nil
	}
	if v.goBytes != nil {
		v.goBytes = nil
	}
}

func newVipsImage(i *C.struct__VipsImage, b []byte) *VipsImage {
	return &VipsImage{cVipsImage: i, goBytes: b}
}

// Decode

type DecodeOptions struct {
	Access VipsAccess
	Disc   bool
}

func (o DecodeOptions) toC() cDecodeOptions {
	return cDecodeOptions{
		Access: o.Access.toC(),
		Disc:   toGBool(o.Disc),
	}
}

type cDecodeOptions struct {
	Access C.VipsAccess
	Disc   C.gboolean
}

func (c *cDecodeOptions) Free() {
}

type DecodeGifOptions struct {
	DecodeOptions
	Page int
}

func (o DecodeGifOptions) toC() cDecodeGifOptions {
	return cDecodeGifOptions{
		cDecodeOptions: o.DecodeOptions.toC(),
		Page:           C.gint(o.Page),
	}
}

type cDecodeGifOptions struct {
	cDecodeOptions
	Page C.gint
}

func (c *cDecodeGifOptions) Free() {
	c.cDecodeOptions.Free()
}

type DecodeJpegOptions struct {
	DecodeOptions
	Shrink     int
	Fail       bool
	Autorotate bool
}

func (o DecodeJpegOptions) toC() cDecodeJpegOptions {
	if o.Shrink >= 8 {
		o.Shrink = 8
	} else if o.Shrink >= 4 {
		o.Shrink = 4
	} else if o.Shrink >= 2 {
		o.Shrink = 2
	} else {
		o.Shrink = 1
	}
	return cDecodeJpegOptions{
		cDecodeOptions: o.DecodeOptions.toC(),
		Shrink:         C.gint(o.Shrink),
		Fail:           toGBool(o.Fail),
		Autorotate:     toGBool(o.Autorotate),
	}
}

type cDecodeJpegOptions struct {
	cDecodeOptions
	Shrink     C.gint
	Fail       C.gboolean
	Autorotate C.gboolean
}

func (c *cDecodeJpegOptions) Free() {
	c.cDecodeOptions.Free()
}

type DecodeMagickOptions struct {
	DecodeOptions
	AllFrames bool
	Density   string
	Page      int
}

func (o DecodeMagickOptions) toC() cDecodeMagickOptions {
	var density *C.char
	if o.Density == STRING_ZERO {
		density = C.CString("")
	} else if o.Density != "" {
		density = C.CString(o.Density)
	}
	return cDecodeMagickOptions{
		cDecodeOptions: o.DecodeOptions.toC(),
		AllFrames:      toGBool(o.AllFrames),
		Density:        density,
		Page:           C.gint(o.Page),
	}
}

type cDecodeMagickOptions struct {
	cDecodeOptions
	AllFrames C.gboolean
	Density   *C.char
	Page      C.gint
}

func (c *cDecodeMagickOptions) Free() {
	c.cDecodeOptions.Free()
	if c.Density != nil {
		C.free(unsafe.Pointer(c.Density))
		c.Density = nil
	}
}

type DecodeWebpOptions struct {
	DecodeOptions
	Shrink int
}

func (o DecodeWebpOptions) toC() cDecodeWebpOptions {
	if o.Shrink < 1 {
		o.Shrink = 1
	}
	return cDecodeWebpOptions{
		cDecodeOptions: o.DecodeOptions.toC(),
		Shrink:         C.gint(o.Shrink),
	}
}

type cDecodeWebpOptions struct {
	cDecodeOptions
	Shrink C.gint
}

func (c *cDecodeWebpOptions) Free() {
	c.cDecodeOptions.Free()
}

func DecodeGifReader(r io.Reader, options *DecodeGifOptions) (*VipsImage, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeGifBytes(b, options)
}

func DecodeGifBytes(b []byte, options *DecodeGifOptions) (*VipsImage, error) {
	if options == nil {
		options = &DecodeGifOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_gifload_buffer(unsafe.Pointer(&b[0]), C.size_t(len(b)), &i, cOptions.Page, cOptions.Access, cOptions.Disc) != 0 {
		return nil, ErrLoad
	}
	return newVipsImage(i, b), nil
}

func DecodeJpegReader(r io.Reader, options *DecodeJpegOptions) (*VipsImage, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeJpegBytes(b, options)
}

func DecodeJpegBytes(b []byte, options *DecodeJpegOptions) (*VipsImage, error) {
	if options == nil {
		options = &DecodeJpegOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_jpegload_buffer(unsafe.Pointer(&b[0]), C.size_t(len(b)), &i, cOptions.Shrink, cOptions.Fail, cOptions.Autorotate, cOptions.Access, cOptions.Disc) != 0 {
		return nil, ErrLoad
	}
	return newVipsImage(i, b), nil
}

func DecodeMagickReader(r io.Reader, options *DecodeMagickOptions) (*VipsImage, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeMagickBytes(b, options)
}

func DecodeMagickBytes(b []byte, options *DecodeMagickOptions) (*VipsImage, error) {
	if options == nil {
		options = &DecodeMagickOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_magickload_buffer(unsafe.Pointer(&b[0]), C.size_t(len(b)), &i, cOptions.AllFrames, cOptions.Density, cOptions.Page, cOptions.Access, cOptions.Disc) != 0 {
		return nil, ErrLoad
	}
	return newVipsImage(i, b), nil
}

func DecodePngReader(r io.Reader, options *DecodeOptions) (*VipsImage, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodePngBytes(b, options)
}

func DecodePngBytes(b []byte, options *DecodeOptions) (*VipsImage, error) {
	if options == nil {
		options = &DecodeOptions{}
	}
	cOptions := options.toC()
	var i *C.struct__VipsImage
	if C.govips_pngload_buffer(unsafe.Pointer(&b[0]), C.size_t(len(b)), &i, cOptions.Access, cOptions.Disc) != 0 {
		return nil, ErrLoad
	}
	return newVipsImage(i, b), nil
}

func DecodeWebpReader(r io.Reader, options *DecodeWebpOptions) (*VipsImage, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeWebpBytes(b, options)
}

func DecodeWebpBytes(b []byte, options *DecodeWebpOptions) (*VipsImage, error) {
	if options == nil {
		options = &DecodeWebpOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_webpload_buffer(unsafe.Pointer(&b[0]), C.size_t(len(b)), &i, cOptions.Shrink, cOptions.Access, cOptions.Disc) != 0 {
		return nil, ErrLoad
	}
	return newVipsImage(i, b), nil
}

// Encode

type EncodeJpegOptions struct {
	Q                   int
	Profile             string
	OptimizeCoding      bool
	Interlace           bool
	Strip               bool
	NoSubsample         bool
	TrellisQuantization bool
	OvershootDeringing  bool
	OptimizeScans       bool
	QuantizationTable   int
}

func (o EncodeJpegOptions) toC() cEncodeJpegOptions {
	if o.Q == 0 {
		o.Q = 75
	} else if o.Q == INT_ZERO {
		o.Q = 0
	}
	var profile *C.char
	if o.Profile == STRING_ZERO {
		profile = C.CString("")
	} else if o.Profile != "" {
		profile = C.CString(o.Profile)
	}
	return cEncodeJpegOptions{
		Q:                   C.gint(o.Q),
		Profile:             profile,
		OptimizeCoding:      toGBool(o.OptimizeCoding),
		Interlace:           toGBool(o.Interlace),
		Strip:               toGBool(o.Strip),
		NoSubsample:         toGBool(o.NoSubsample),
		TrellisQuantization: toGBool(o.TrellisQuantization),
		OvershootDeringing:  toGBool(o.OvershootDeringing),
		OptimizeScans:       toGBool(o.OptimizeScans),
		QuantizationTable:   C.gint(o.QuantizationTable),
	}
}

type cEncodeJpegOptions struct {
	Q                   C.gint
	Profile             *C.char
	OptimizeCoding      C.gboolean
	Interlace           C.gboolean
	Strip               C.gboolean
	NoSubsample         C.gboolean
	TrellisQuantization C.gboolean
	OvershootDeringing  C.gboolean
	OptimizeScans       C.gboolean
	QuantizationTable   C.gint
}

func (c *cEncodeJpegOptions) Free() {
	if c.Profile != nil {
		C.free(unsafe.Pointer(c.Profile))
		c.Profile = nil
	}
}

func EncodeJpegFile(i *VipsImage, file *os.File, options *EncodeJpegOptions) error {
	if options == nil {
		options = &EncodeJpegOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	cFileName := C.CString(file.Name())
	defer C.free(unsafe.Pointer(cFileName))
	if C.govips_jpegsave(i.cVipsImage, cFileName, cOptions.Q, cOptions.Profile, cOptions.OptimizeCoding, cOptions.Interlace, cOptions.Strip, cOptions.NoSubsample, cOptions.TrellisQuantization, cOptions.OvershootDeringing, cOptions.OptimizeScans, cOptions.QuantizationTable) != 0 {
		return ErrSave
	}
	return nil
}

func EncodeJpegBytes(i *VipsImage, options *EncodeJpegOptions) ([]byte, error) {
	if options == nil {
		options = &EncodeJpegOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var obuf unsafe.Pointer
	olen := C.size_t(0)
	if C.govips_jpegsave_buffer(i.cVipsImage, &obuf, &olen, cOptions.Q, cOptions.Profile, cOptions.OptimizeCoding, cOptions.Interlace, cOptions.Strip, cOptions.NoSubsample, cOptions.TrellisQuantization, cOptions.OvershootDeringing, cOptions.OptimizeScans, cOptions.QuantizationTable) != 0 {
		return nil, ErrSave
	}
	defer C.g_free(C.gpointer(obuf))
	bytes := C.GoBytes(obuf, C.int(olen))
	return bytes, nil
}

type EncodePngOptions struct {
	Compression int
	Interlace   bool
	Profile     string
	Filter      PngFilter
}

func (o EncodePngOptions) toC() cEncodePngOptions {
	if o.Compression == 0 {
		o.Compression = 6
	} else if o.Compression == INT_ZERO {
		o.Compression = 0
	}
	var profile *C.char
	if o.Profile == STRING_ZERO {
		profile = C.CString("")
	} else if o.Profile != "" {
		profile = C.CString(o.Profile)
	}
	return cEncodePngOptions{
		Compression: C.gint(o.Compression),
		Interlace:   toGBool(o.Interlace),
		Profile:     profile,
		Filter:      o.Filter.toC(),
	}
}

type cEncodePngOptions struct {
	Compression C.gint
	Interlace   C.gboolean
	Profile     *C.char
	Filter      C.VipsForeignPngFilter
}

func (c *cEncodePngOptions) Free() {
	if c.Profile != nil {
		C.free(unsafe.Pointer(c.Profile))
		c.Profile = nil
	}
}

func EncodePngFile(i *VipsImage, file *os.File, options *EncodePngOptions) error {
	if options == nil {
		options = &EncodePngOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	cFileName := C.CString(file.Name())
	defer C.free(unsafe.Pointer(cFileName))
	if C.govips_pngsave(i.cVipsImage, cFileName, cOptions.Compression, cOptions.Interlace, cOptions.Profile, cOptions.Filter) != 0 {
		return ErrSave
	}
	return nil
}

func EncodePngBytes(i *VipsImage, options *EncodePngOptions) ([]byte, error) {
	if options == nil {
		options = &EncodePngOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var obuf unsafe.Pointer
	olen := C.size_t(0)
	if C.govips_pngsave_buffer(i.cVipsImage, &obuf, &olen, cOptions.Compression, cOptions.Interlace, cOptions.Profile, cOptions.Filter) != 0 {
		return nil, ErrSave
	}
	defer C.g_free(C.gpointer(obuf))
	bytes := C.GoBytes(obuf, C.int(olen))
	return bytes, nil
}

type EncodeWebpOptions struct {
	Q              int
	Lossless       bool
	Preset         WebpPreset
	SmartSubsample bool
	NearLossless   bool
	AlphaQ         int
}

func (o EncodeWebpOptions) toC() cEncodeWebpOptions {
	if o.Q == 0 {
		o.Q = 75
	} else if o.Q == INT_ZERO {
		o.Q = 0
	}
	if o.AlphaQ == 0 {
		o.AlphaQ = 100
	}
	return cEncodeWebpOptions{
		Q:              C.gint(o.Q),
		Lossless:       toGBool(o.Lossless),
		Preset:         o.Preset.toC(),
		SmartSubsample: toGBool(o.SmartSubsample),
		NearLossless:   toGBool(o.NearLossless),
		AlphaQ:         C.gint(o.AlphaQ),
	}
}

type cEncodeWebpOptions struct {
	Q              C.gint
	Lossless       C.gboolean
	Preset         C.VipsForeignWebpPreset
	SmartSubsample C.gboolean
	NearLossless   C.gboolean
	AlphaQ         C.gint
}

func (c *cEncodeWebpOptions) Free() {
}

func EncodeWebpFile(i *VipsImage, file *os.File, options *EncodeWebpOptions) error {
	if options == nil {
		options = &EncodeWebpOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	cFileName := C.CString(file.Name())
	defer C.free(unsafe.Pointer(cFileName))
	if C.govips_webpsave(i.cVipsImage, cFileName, cOptions.Q, cOptions.Lossless, cOptions.Preset, cOptions.SmartSubsample, cOptions.NearLossless, cOptions.AlphaQ) != 0 {
		return ErrSave
	}
	return nil
}

func EncodeWebpBytes(i *VipsImage, options *EncodeWebpOptions) ([]byte, error) {
	if options == nil {
		options = &EncodeWebpOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var obuf unsafe.Pointer
	olen := C.size_t(0)
	if C.govips_webpsave_buffer(i.cVipsImage, &obuf, &olen, cOptions.Q, cOptions.Lossless, cOptions.Preset, cOptions.SmartSubsample, cOptions.NearLossless, cOptions.AlphaQ) != 0 {
		return nil, ErrSave
	}
	defer C.g_free(C.gpointer(obuf))
	bytes := C.GoBytes(obuf, C.int(olen))
	return bytes, nil
}

// Operations

type EmbedOptions struct {
	Extend     VipsExtend
	Background []float64
}

func (o EmbedOptions) toC() cEmbedOptions {
	var Background *C.struct__VipsArrayDouble
	if o.Background != nil && len(o.Background) > 0 {
		Background = newVipsArrayDouble(o.Background)
	}
	return cEmbedOptions{
		Extend:     o.Extend.toC(),
		Background: Background,
	}
}

type cEmbedOptions struct {
	Extend     C.VipsExtend
	Background *C.struct__VipsArrayDouble
}

func (c *cEmbedOptions) Free() {
	if c.Background != nil {
		vipsArrayDoubleUnref(c.Background)
		c.Background = nil
	}
}

func Embed(v *VipsImage, x, y, width, height int, options *EmbedOptions) (*VipsImage, error) {
	if options == nil {
		options = &EmbedOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_embed(v.cVipsImage, &i, C.int(x), C.int(y), C.int(width), C.int(height), cOptions.Extend, cOptions.Background) != 0 {
		return nil, ErrEmbed
	}
	return newVipsImage(i, v.goBytes), nil
}

func ExtractArea(v *VipsImage, left, top, width, height int) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_extract_area(v.cVipsImage, &i, C.int(left), C.int(top), C.int(width), C.int(height)) != 0 {
		return nil, ErrCrop
	}
	return newVipsImage(i, v.goBytes), nil
}

func Crop(v *VipsImage, left, top, width, height int) (*VipsImage, error) {
	return ExtractArea(v, left, top, width, height)
}

func Shrink(v *VipsImage, xshrink, yshrink float64) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_shrink(v.cVipsImage, &i, C.double(xshrink), C.double(yshrink)) != 0 {
		return nil, ErrShrink
	}
	return newVipsImage(i, v.goBytes), nil
}

func ShrinkH(v *VipsImage, xshrink float64) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_shrinkh(v.cVipsImage, &i, C.double(xshrink)) != 0 {
		return nil, ErrShrink
	}
	return newVipsImage(i, v.goBytes), nil
}

func ShrinkV(v *VipsImage, yshrink float64) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_shrinkv(v.cVipsImage, &i, C.double(yshrink)) != 0 {
		return nil, ErrShrink
	}
	return newVipsImage(i, v.goBytes), nil
}

func Reduce(v *VipsImage, xshrink, yshrink float64, kernel VipsKernel) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_reduce(v.cVipsImage, &i, C.double(xshrink), C.double(yshrink), C.VipsKernel(kernel)) != 0 {
		return nil, ErrReduce
	}
	return newVipsImage(i, v.goBytes), nil
}

func ReduceH(v *VipsImage, xshrink float64, kernel VipsKernel) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_reduceh(v.cVipsImage, &i, C.double(xshrink), C.VipsKernel(kernel)) != 0 {
		return nil, ErrReduce
	}
	return newVipsImage(i, v.goBytes), nil
}

func ReduceV(v *VipsImage, yshrink float64, kernel VipsKernel) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_reducev(v.cVipsImage, &i, C.double(yshrink), C.VipsKernel(kernel)) != 0 {
		return nil, ErrReduce
	}
	return newVipsImage(i, v.goBytes), nil
}

func Resize(v *VipsImage, scale, vscale float64, kernel VipsKernel) (*VipsImage, error) {
	var i *C.struct__VipsImage
	if C.govips_resize(v.cVipsImage, &i, C.double(scale), C.double(vscale), C.VipsKernel(kernel)) != 0 {
		return nil, ErrResize
	}
	return newVipsImage(i, v.goBytes), nil
}

type SimilarityOptions struct {
	Scale       float64
	Angle       float64
	Interpolate *VipsInterpolate
	Idx         float64
	Idy         float64
	Odx         float64
	Ody         float64
}

func (o SimilarityOptions) toC() cSimilarityOptions {
	if o.Scale == 0 {
		o.Scale = 1
	} else if o.Scale == FLOAT_ZERO {
		o.Scale = 0
	}
	if o.Scale > 1 && o.Angle == 0 {
		o.Angle = 360
	}
	var interpolate *VipsInterpolate
	if o.Interpolate == nil {
		interpolate = NewBilinearVipsInterpolator()
	}
	return cSimilarityOptions{
		Scale:       C.gdouble(o.Scale),
		Angle:       C.gdouble(o.Angle),
		Interpolate: interpolate.cVipsInterpolate,
		Idx:         C.gdouble(o.Idx),
		Idy:         C.gdouble(o.Idy),
		Odx:         C.gdouble(o.Odx),
		Ody:         C.gdouble(o.Ody),
	}
}

type cSimilarityOptions struct {
	Scale       C.gdouble
	Angle       C.gdouble
	Interpolate *C.struct__VipsInterpolate
	Idx         C.gdouble
	Idy         C.gdouble
	Odx         C.gdouble
	Ody         C.gdouble
}

func (c *cSimilarityOptions) Free() {
	if c.Interpolate != nil {
		C.g_object_unref(C.gpointer(c.Interpolate))
		c.Interpolate = nil
	}
}

func Similarity(v *VipsImage, options *SimilarityOptions) (*VipsImage, error) {
	if options == nil {
		options = &SimilarityOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_similarity(v.cVipsImage, &i, cOptions.Scale, cOptions.Angle, cOptions.Interpolate, cOptions.Idx, cOptions.Idy, cOptions.Odx, cOptions.Ody) != 0 {
		return nil, ErrAffine
	}
	return newVipsImage(i, v.goBytes), nil
}

type AffineOptions struct {
	Interpolate *VipsInterpolate
	OArea       []int
	Idx         float64
	Idy         float64
	Odx         float64
	Ody         float64
}

func (o AffineOptions) toC() cAffineOptions {
	var interpolate *VipsInterpolate
	if o.Interpolate == nil {
		interpolate = NewBilinearVipsInterpolator()
	}
	var oArea *C.struct__VipsArrayInt
	if o.OArea != nil && len(o.OArea) > 0 {
		oArea = newVipsArrayInt(o.OArea)
	}
	return cAffineOptions{
		Interpolate: interpolate.cVipsInterpolate,
		OArea:       oArea,
		Idx:         C.gdouble(o.Idx),
		Idy:         C.gdouble(o.Idy),
		Odx:         C.gdouble(o.Odx),
		Ody:         C.gdouble(o.Ody),
	}
}

type cAffineOptions struct {
	Interpolate *C.struct__VipsInterpolate
	OArea       *C.struct__VipsArrayInt
	Idx         C.gdouble
	Idy         C.gdouble
	Odx         C.gdouble
	Ody         C.gdouble
}

func (c *cAffineOptions) Free() {
	if c.Interpolate != nil {
		C.g_object_unref(C.gpointer(c.Interpolate))
		c.Interpolate = nil
	}
	if c.OArea != nil {
		vipsArrayIntUnref(c.OArea)
		c.OArea = nil
	}
}

func Affine(v *VipsImage, a, b, c, d float64, options *AffineOptions) (*VipsImage, error) {
	if options == nil {
		options = &AffineOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_affine(v.cVipsImage, &i, C.double(a), C.double(b), C.double(c), C.double(d), cOptions.Interpolate, cOptions.OArea, cOptions.Idx, cOptions.Idy, cOptions.Odx, cOptions.Ody) != 0 {
		return nil, ErrAffine
	}
	return newVipsImage(i, v.goBytes), nil
}

type BlurOptions struct {
	Precision        VipsPrecision
	MinimumAmplitude float64
}

func (o BlurOptions) toC() cBlurOptions {
	if o.MinimumAmplitude == 0 {
		o.MinimumAmplitude = 0.2
	}
	return cBlurOptions{
		Precision:        C.VipsPrecision(o.Precision),
		MinimumAmplitude: C.double(o.MinimumAmplitude),
	}
}

type cBlurOptions struct {
	Precision        C.VipsPrecision
	MinimumAmplitude C.double
}

func (c *cBlurOptions) Free() {
}

func Blur(v *VipsImage, sigma float64, options *BlurOptions) (*VipsImage, error) {
	if options == nil {
		options = &BlurOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_gaussblur(v.cVipsImage, &i, C.double(sigma), cOptions.Precision, cOptions.MinimumAmplitude) != 0 {
		return nil, ErrBlur
	}
	return newVipsImage(i, v.goBytes), nil
}

type SharpenOptions struct {
	Sigma float64
	X1    float64
	Y2    float64
	Y3    float64
	M1    float64
	M2    float64
}

func (o SharpenOptions) toC() cSharpenOptions {
	if o.Sigma == 0 {
		o.Sigma = 0.5
	}
	if o.X1 == 0 {
		o.X1 = 2.0
	} else if o.X1 == FLOAT_ZERO {
		o.X1 = 0
	}
	if o.Y2 == 0 {
		o.Y2 = 10.0
	} else if o.Y2 == FLOAT_ZERO {
		o.Y2 = 0
	}
	if o.Y3 == 0 {
		o.Y3 = 20.0
	} else if o.Y3 == FLOAT_ZERO {
		o.Y3 = 0
	}
	if o.M1 == FLOAT_ZERO {
		o.M1 = 0
	}
	if o.M2 == 0 {
		o.M2 = 3.0
	} else if o.M2 == FLOAT_ZERO {
		o.Sigma = 0
	}
	return cSharpenOptions{
		Sigma: C.double(o.Sigma),
		X1:    C.double(o.X1),
		Y2:    C.double(o.Y2),
		Y3:    C.double(o.Y3),
		M1:    C.double(o.M1),
		M2:    C.double(o.M2),
	}
}

type cSharpenOptions struct {
	Sigma C.double
	X1    C.double
	Y2    C.double
	Y3    C.double
	M1    C.double
	M2    C.double
}

func (c *cSharpenOptions) Free() {
}

func Sharpen(v *VipsImage, options *SharpenOptions) (*VipsImage, error) {
	if options == nil {
		options = &SharpenOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_sharpen(v.cVipsImage, &i, cOptions.Sigma, cOptions.X1, cOptions.Y2, cOptions.Y3, cOptions.M1, cOptions.M2) != 0 {
		return nil, ErrSharpen
	}
	return newVipsImage(i, v.goBytes), nil
}

type FlattenOptions struct {
	Background []float64
	MaxAlpha   float64
}

func (o FlattenOptions) toC() cFlattenOptions {
	var Background *C.struct__VipsArrayDouble
	if o.Background != nil && len(o.Background) > 0 {
		Background = newVipsArrayDouble(o.Background)
	}
	if o.MaxAlpha == 0 {
		o.MaxAlpha = 255
	} else if o.MaxAlpha == FLOAT_ZERO {
		o.MaxAlpha = 0
	}
	return cFlattenOptions{
		Background: Background,
		MaxAlpha:   C.double(o.MaxAlpha),
	}
}

type cFlattenOptions struct {
	Background *C.struct__VipsArrayDouble
	MaxAlpha   C.double
}

func (c *cFlattenOptions) Free() {
	if c.Background != nil {
		vipsArrayDoubleUnref(c.Background)
		c.Background = nil
	}
}

func Flatten(v *VipsImage, options *FlattenOptions) (*VipsImage, error) {
	if options == nil {
		options = &FlattenOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_flatten(v.cVipsImage, &i, cOptions.Background, cOptions.MaxAlpha) != 0 {
		return nil, ErrFlatten
	}
	return newVipsImage(i, v.goBytes), nil
}

type ColourspaceOptions struct {
	SourceSpace VipsInterpretation
}

func (o ColourspaceOptions) toC() cColourspaceOptions {
	return cColourspaceOptions{
		SourceSpace: o.SourceSpace.toC(),
	}
}

type cColourspaceOptions struct {
	SourceSpace C.VipsInterpretation
}

func (c *cColourspaceOptions) Free() {
}

func Colourspace(v *VipsImage, space VipsInterpretation, options *ColourspaceOptions) (*VipsImage, error) {
	if options == nil {
		options = &ColourspaceOptions{
			SourceSpace: v.Interpretation(),
		}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	var i *C.struct__VipsImage
	if C.govips_colourspace(v.cVipsImage, &i, space.toC(), cOptions.SourceSpace) != 0 {
		return nil, ErrColourspace
	}
	return newVipsImage(i, v.goBytes), nil
}

func ColourspaceIsSupported(v *VipsImage) bool {
	return fromGBool(C.vips_colourspace_issupported(v.cVipsImage))
}

type ICCTransformOptions struct {
	InputProfile string
	Intent       VipsIntent
	Depth        int
	Embedded     bool
}

func (o ICCTransformOptions) toC() cICCTransformOptions {
	var inputProfile *C.char
	if o.InputProfile == STRING_ZERO {
		inputProfile = C.CString("")
	} else if o.InputProfile != "" {
		inputProfile = C.CString(o.InputProfile)
	}
	if o.Depth == 0 {
		o.Depth = 8
	} else if o.Depth == INT_ZERO {
		o.Depth = 0
	}
	return cICCTransformOptions{
		InputProfile: inputProfile,
		Intent:       o.Intent.toC(),
		Depth:        C.int(o.Depth),
		Embedded:     toGBool(o.Embedded),
	}
}

type cICCTransformOptions struct {
	InputProfile *C.char
	Intent       C.VipsIntent
	Depth        C.int
	Embedded     C.gboolean
}

func (c *cICCTransformOptions) Free() {
	if c.InputProfile != nil {
		C.free(unsafe.Pointer(c.InputProfile))
		c.InputProfile = nil
	}
}

func ICCTransform(v *VipsImage, outputProfile string, options *ICCTransformOptions) (*VipsImage, error) {
	if options == nil {
		options = &ICCTransformOptions{}
	}
	cOptions := options.toC()
	defer cOptions.Free()
	p := C.CString(outputProfile)
	defer C.free(unsafe.Pointer(p))
	var i *C.struct__VipsImage
	if C.govips_icc_transform(v.cVipsImage, &i, p, cOptions.InputProfile, cOptions.Intent, cOptions.Depth, cOptions.Embedded) != 0 {
		return nil, ErrICCTransform
	}
	return newVipsImage(i, v.goBytes), nil
}

// Interpolators

type VipsInterpolate struct {
	Nickname         string
	cVipsInterpolate *C.struct__VipsInterpolate
}

func (i *VipsInterpolate) Free() {
	if i.cVipsInterpolate != nil {
		i.Nickname = fmt.Sprintf("%s (freed)", i.Nickname)
		C.g_object_unref(C.gpointer(i.cVipsInterpolate))
		i.cVipsInterpolate = nil
	}
}

func NewVipsInterpolator(interpolator string) (*VipsInterpolate, error) {
	s := C.CString(interpolator)
	defer C.free(unsafe.Pointer(s))
	cVipsInterpolate := C.vips_interpolate_new(s)
	if cVipsInterpolate == nil {
		return nil, fmt.Errorf("Failed to create interpolator for: %s", interpolator)
	}
	return &VipsInterpolate{interpolator, cVipsInterpolate}, nil
}

func NewNearestVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("nearest")
	return result
}

func NewBilinearVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("bilinear")
	return result
}

func NewBicubicVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("bicubic")
	return result
}

func NewLBBVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("lbb")
	return result
}

func NewNohaloVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("nohalo")
	return result
}

func NewVSQBSVipsInterpolator() *VipsInterpolate {
	result, _ := NewVipsInterpolator("vsqbs")
	return result
}

// Helpers...

func newVipsRegion(i *VipsImage, bounds image.Rectangle) *C.VipsRegion {
	vipsRegion := C.vips_region_new(i.cVipsImage)
	rect := C.govips_rect_new(C.int(bounds.Min.X), C.int(bounds.Min.Y), C.int(bounds.Dx()), C.int(bounds.Dy()))
	C.vips_region_prepare(vipsRegion, &rect)
	return vipsRegion
}

type NRGBAVipsImage struct {
	*VipsImage
	cVipsRegion *C.VipsRegion
	alphaBand   bool
}

func (v *NRGBAVipsImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (v *NRGBAVipsImage) At(x, y int) color.Color {
	if v.cVipsRegion == nil {
		v.cVipsRegion = newVipsRegion(v.VipsImage, v.Bounds())
		v.alphaBand = C.govips_vips_region_n_elements(v.cVipsRegion) == 4
	}
	vipsPel := C.govips_region_addr(v.cVipsRegion, C.int(x), C.int(y))
	red := uint8(*vipsPel)
	green := uint8(*C.govips_pel_band(vipsPel, 1))
	blue := uint8(*C.govips_pel_band(vipsPel, 2))
	alpha := uint8(255)
	if v.alphaBand {
		alpha = uint8(*C.govips_pel_band(vipsPel, 3))
	}
	return &color.NRGBA{red, green, blue, alpha}
}

func (v *NRGBAVipsImage) Free() {
	if v.cVipsRegion != nil {
		C.g_object_unref(C.gpointer(v.cVipsRegion))
		v.cVipsRegion = nil
	}
	v.VipsImage.Free()
}

func NewNRGBAVipsImage(vi *VipsImage) (*NRGBAVipsImage, error) {
	bands := vi.Bands()
	if !(bands == 3 || bands == 4) {
		return nil, fmt.Errorf("Invalid number of bands: %d", bands)
	}
	format := C.vips_image_get_format(vi.cVipsImage)
	if C.vips_image_get_format(vi.cVipsImage) != C.VIPS_FORMAT_UCHAR {
		return nil, fmt.Errorf("Invalid band format: %v", format)
	}
	if vi.Interpretation() != VIPS_INTERPRETATION_sRGB {
		return nil, fmt.Errorf("Invalid interpretation: %v", vi.Interpretation())
	}
	return &NRGBAVipsImage{
		VipsImage: vi,
	}, nil
}

type CMYKVipsImage struct {
	*VipsImage
	cVipsRegion *C.VipsRegion
}

func (v *CMYKVipsImage) ColorModel() color.Model {
	return color.CMYKModel
}

func (v *CMYKVipsImage) At(x, y int) color.Color {
	if v.cVipsRegion == nil {
		v.cVipsRegion = newVipsRegion(v.VipsImage, v.Bounds())
	}
	vipsPel := C.govips_region_addr(v.cVipsRegion, C.int(x), C.int(y))
	cyan := uint8(*vipsPel)
	magenta := uint8(*C.govips_pel_band(vipsPel, 1))
	yellow := uint8(*C.govips_pel_band(vipsPel, 2))
	black := uint8(*C.govips_pel_band(vipsPel, 3))
	return &color.CMYK{cyan, magenta, yellow, black}
}

func (v *CMYKVipsImage) Free() {
	if v.cVipsRegion != nil {
		C.g_object_unref(C.gpointer(v.cVipsRegion))
		v.cVipsRegion = nil
	}
	v.VipsImage.Free()
}

func NewCMYKVipsImage(vi *VipsImage) (*CMYKVipsImage, error) {
	bands := vi.Bands()
	if bands != 4 {
		return nil, fmt.Errorf("Invalid number of bands: %d", bands)
	}
	format := C.vips_image_get_format(vi.cVipsImage)
	if C.vips_image_get_format(vi.cVipsImage) != C.VIPS_FORMAT_UCHAR {
		return nil, fmt.Errorf("Invalid band format: %v", format)
	}
	if vi.Interpretation() != VIPS_INTERPRETATION_CMYK {
		return nil, fmt.Errorf("Invalid interpretation: %v", vi.Interpretation())
	}
	return &CMYKVipsImage{
		VipsImage: vi,
	}, nil
}

type GrayVipsImage struct {
	*VipsImage
	cVipsRegion *C.VipsRegion
}

func (v *GrayVipsImage) ColorModel() color.Model {
	return color.GrayModel
}

func (v *GrayVipsImage) At(x, y int) color.Color {
	if v.cVipsRegion == nil {
		v.cVipsRegion = newVipsRegion(v.VipsImage, v.Bounds())
	}
	vipsPel := C.govips_region_addr(v.cVipsRegion, C.int(x), C.int(y))
	return &color.Gray{uint8(*vipsPel)}
}

func (v *GrayVipsImage) Free() {
	if v.cVipsRegion != nil {
		C.g_object_unref(C.gpointer(v.cVipsRegion))
		v.cVipsRegion = nil
	}
	v.VipsImage.Free()
}

func NewGrayVipsImage(vi *VipsImage) (*GrayVipsImage, error) {
	bands := vi.Bands()
	if bands != 1 {
		return nil, fmt.Errorf("Invalid number of bands: %d", bands)
	}
	format := C.vips_image_get_format(vi.cVipsImage)
	if C.vips_image_get_format(vi.cVipsImage) != C.VIPS_FORMAT_UCHAR {
		return nil, fmt.Errorf("Invalid band format: %v", format)
	}
	if vi.Interpretation() != VIPS_INTERPRETATION_B_W {
		return nil, fmt.Errorf("Invalid interpretation: %v", vi.Interpretation())
	}
	return &GrayVipsImage{
		VipsImage: vi,
	}, nil
}

// Utilities...

func newVipsArrayInt(slice []int) *C.struct__VipsArrayInt {
	return C.vips_array_int_new((*C.int)(unsafe.Pointer(&slice[0])), C.int(len(slice)))
}

func newVipsArrayDouble(slice []float64) *C.struct__VipsArrayDouble {
	return C.vips_array_double_new((*C.double)(unsafe.Pointer(&slice[0])), C.int(len(slice)))
}

func vipsArrayDoubleUnref(i *C.struct__VipsArrayDouble) {
	vipsAreaUnref((*C.struct__VipsArea)(unsafe.Pointer(i)))
}

func vipsArrayIntUnref(i *C.struct__VipsArrayInt) {
	vipsAreaUnref((*C.struct__VipsArea)(unsafe.Pointer(i)))
}

func vipsAreaUnref(i *C.struct__VipsArea) {
	C.vips_area_unref(i)
}

func toGBool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func fromGBool(b C.gboolean) bool {
	return int(b) == 1
}
