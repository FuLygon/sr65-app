package internal

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

var (
	testImage           = "test.png"
	outputDir           = "outputs_test"
	outputStaticExt     = "jpg"
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
	f, err := os.Create(testImage)
	if err != nil {
		panic(err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	// encode test image
	if err = png.Encode(f, img); err != nil {
		panic(err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		panic(err)
	}

	// convert static media
	ConvertStatic(f.Name(), outputDir, outputStaticExt, outputStaticQuality)

	// check if the output file exist
	outputPath := generateOutput(f.Name(), outputDir, outputStaticExt)
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

func TestGenerateOutputStatic(t *testing.T) {
	outputPath := generateOutput(testImage, outputDir, outputStaticExt)

	// check output path
	expectedOutputPath := filepath.Join(outputDir, "test."+outputStaticExt)
	if outputPath != expectedOutputPath {
		t.Errorf("expected %s, got %s", expectedOutputPath, outputPath)
	}
}
