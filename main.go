package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/ncruces/zenity"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sr65-software/embed"
	"sr65-software/logger"
	"strings"
)

const (
	outputDir           = "outputs"
	outputStaticExt     = "jpg"
	outputStaticQuality = 95
)

func main() {
	// embedding binaries
	cleanup, err := embed.ExtractBinaries()
	if err != nil {
		logger.Warn("error extracting embedded binaries, falling back to system ffmpeg")
	}
	defer cleanup()

	// show question dialog
	err = zenity.Question(
		"Choose an image or video file. Supported file formats:"+"\n"+
			"- Image: png, jpg, gif"+"\n"+
			"- Video: mp4",
		zenity.Width(400),
		zenity.Title("SR65 Software"),
		zenity.OKLabel("Open File..."),
		zenity.CancelLabel("Close"),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			logger.Fatal("user cancelled", err)
		}
		logger.Fatal("an error occurred", err)
	}

	// show file picker
	inputPath, err := zenity.SelectFile(
		zenity.Title("Select a file"),
		zenity.FileFilters{
			{"Media files", []string{"*.png", "*.jpg", "*.gif", "*.mp4"}, true},
		},
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			logger.Fatal("user cancelled", err)
		}
		logger.Fatal("an error occurred", err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0644)
	if err != nil {
		logger.Fatal("error creating output directory", err)
	}

	// convert input file
	logger.Info("converting file")
	switch strings.ToLower(filepath.Ext(inputPath)) {
	case ".jpg", ".png":
		convertStatic(inputPath)
	case ".gif", ".mp4":
		convertDynamic()
	default:
		logger.Fatal("unsupported file format")
	}

	logger.Info(`converted successfully, output saved in "outputs" directory. Exiting...`)
}

func convertStatic(inputPath string) {
	// open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		logger.Fatal("error opening input file", err)
	}
	defer func(inputFile *os.File) {
		err = inputFile.Close()
		if err != nil {
			logger.Error("error closing input file", err)
		}
	}(inputFile)

	// decode image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		logger.Fatal("error decoding image", err)
	}

	// resizing image
	img = imaging.Resize(img, 128, 128, imaging.Lanczos)

	// create output file
	outputFile, err := os.Create(generateOutput(inputPath, outputStaticExt))
	if err != nil {
		logger.Fatal("error creating output file", err)
	}
	defer func(outputFile *os.File) {
		err = outputFile.Close()
		if err != nil {
			logger.Error("error closing output file", err)
		}
	}(outputFile)

	// encode img to output file
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: outputStaticQuality})
	if err != nil {
		logger.Fatal("error encoding image", err)
	}

	return
}

func convertDynamic() {
}

func generateOutput(inputPath, outputExt string) string {
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputPath)
	inputName := inputBase[0 : len(inputBase)-len(inputExt)]
	return fmt.Sprintf("%s/%s.%s", outputDir, inputName, outputExt)
}
