package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
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
	f, err := os.Create("test.png")
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

	// convert static media
	convertStatic(f.Name())

	// check if the output file exist
	outputPath := generateOutputStatic(f.Name())
	defer os.Remove(outputPath)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
	}

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
	// Prepare a test image file
	testImagePath := "test.png"

	// Call the function
	outputPath := generateOutputStatic(testImagePath)

	// Check if the output path is correct
	expectedOutputPath := filepath.Join(outputDir, "test."+outputStaticExt)
	if outputPath != expectedOutputPath {
		t.Errorf("expected %s, got %s", expectedOutputPath, outputPath)
	}
}
