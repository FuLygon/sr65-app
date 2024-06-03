package internal

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

const (
	testImage           = "test.png"
	testGif             = "test.gif"
	testVideo           = "test.mp4"
	outputDir           = "outputs_test"
	outputStaticExt     = "jpg"
	outputDynamicExt    = "mjpeg"
	outputStaticQuality = 95
)

func TestConvertStatic(t *testing.T) {
	// prepare test file
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.RGBA{})
		}
	}

	// create test file
	testFile, err := os.Create(testImage)
	if err != nil {
		panic(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// encode test image
	if err = png.Encode(testFile, img); err != nil {
		panic(err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	// convert static media
	ConvertStatic(testFile.Name(), outputDir, outputStaticExt, outputStaticQuality)

	// check if the output file exist
	outputPath := generateOutput(testFile.Name(), outputDir, outputStaticExt)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
	}
	defer os.RemoveAll(outputDir)

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
	}
	defer outputFile.Close()

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
	}

	// check output file dimensions
	if config.Width != 128 || config.Height != 128 {
		t.Errorf("output image dimensions is incorrect, got: %dx%d, want: 128x128", config.Width, config.Height)
	}
}

func TestConvertDynamicGif(t *testing.T) {
	// prepare test file
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.RGBA{})
		}
	}

	// create test file
	testFile, err := os.Create(testGif)
	if err != nil {
		panic(err)
	}
	defer func() {
		testFile.Close()
		os.Remove(testFile.Name())
	}()

	// encode test gif
	if err = gif.Encode(testFile, img, nil); err != nil {
		panic(err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	// convert gif
	ConvertDynamic(testFile.Name(), outputDir, outputDynamicExt)

	// check if the output file exist
	outputPath := generateOutput(testVideo, outputDir, outputDynamicExt)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
	}
	defer os.RemoveAll(outputDir)

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
	}
	defer outputFile.Close()

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
	}

	// check output file dimensions
	if config.Width != 128 || config.Height != 128 {
		t.Errorf("output gif dimensions is incorrect, got: %dx%d, want: 128x128", config.Width, config.Height)
	}
}

func TestConvertDynamicVideo(t *testing.T) {
	// generate blank video
	err := ffmpeg.Input("color=c=black:s=512x512:d=5", ffmpeg.KwArgs{"f": "lavfi"}).
		Output(testVideo, ffmpeg.KwArgs{"c:v": "libx264", "pix_fmt": "yuv420p", "t": "5"}).
		OverWriteOutput().
		Silent(ffmpegSilent).
		Run()
	if err != nil {
		panic(err)
	}
	defer os.Remove(testVideo)

	// get test file path
	testVideoPath, err := filepath.Abs(testVideo)
	if err != nil {
		panic(err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	// convert mp4
	ConvertDynamic(testVideoPath, outputDir, outputDynamicExt)

	// check if the output file exist
	outputPath := generateOutput(testVideo, outputDir, outputDynamicExt)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
	}
	defer os.RemoveAll(outputDir)

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
	}
	defer outputFile.Close()

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
	}

	// check output file dimensions
	if config.Width != 128 || config.Height != 128 {
		t.Errorf("output video dimensions is incorrect, got: %dx%d, want: 128x128", config.Width, config.Height)
	}
}

func TestGenerateOutput(t *testing.T) {
	outputPath := generateOutput(testImage, outputDir, outputStaticExt)

	// check output path
	expectedOutputPath := filepath.Join(outputDir, "test."+outputStaticExt)
	if outputPath != expectedOutputPath {
		t.Errorf("expected %s, got %s", expectedOutputPath, outputPath)
	}
}
