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
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Logf("error closing test image: %v", err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			t.Logf("error removing test image: %v", err)
		}
	}(testFile)

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
	err = ConvertStatic(testFile.Name(), outputDir, outputStaticExt, outputStaticQuality)
	if err != nil {
		t.Errorf("error converting image: %v", err)
		t.FailNow()
	}

	// check if the output file exist
	outputPath := generateOutput(testFile.Name(), outputDir, outputStaticExt)
	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			t.Logf("error removing output directory: %v", err)
		}
	}(outputDir)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
		t.FailNow()
	}

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
		t.FailNow()
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Logf("error closing output file: %v", err)
		}
	}(outputFile)

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
		t.FailNow()
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
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Logf("error closing test gif: %v", err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			t.Logf("error removing test gif: %v", err)
		}
	}(testFile)

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
	err = ConvertDynamic(testFile.Name(), outputDir, outputDynamicExt)
	if err != nil {
		t.Errorf("error converting gif: %v", err)
		t.FailNow()
	}

	// check if the output file exist
	outputPath := generateOutput(testVideo, outputDir, outputDynamicExt)
	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			t.Logf("error removing output directory: %v", err)
		}
	}(outputDir)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
		t.FailNow()
	}

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
		t.FailNow()
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Logf("error closing output file: %v", err)
			t.FailNow()
		}
	}(outputFile)

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
		t.FailNow()
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
	defer func(path string) {
		err = os.Remove(path)
		if err != nil {
			t.Logf("error removing test video: %v", err)
		}
	}(testVideo)

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
	err = ConvertDynamic(testVideoPath, outputDir, outputDynamicExt)
	if err != nil {
		t.Errorf("error converting video: %v", err)
		t.FailNow()
	}

	// check if the output file exist
	outputPath := generateOutput(testVideo, outputDir, outputDynamicExt)
	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			t.Logf("error removing output directory: %v", err)
		}
	}(outputDir)
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("output file was not created")
		t.FailNow()
	}

	// open output file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Errorf("error opening output file: %v", err)
		t.FailNow()
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			t.Logf("error closing output file: %v", err)
		}
	}(outputFile)

	// decode output file
	config, _, err := image.DecodeConfig(outputFile)
	if err != nil {
		t.Errorf("error decoding output file: %v", err)
		t.FailNow()
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
