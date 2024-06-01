package internal

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"sr65-software/logger"
)

func ConvertStatic(inputPath, outputDir, outputExt string, jpegQuality int) {
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
	outputFile, err := os.Create(generateOutput(inputPath, outputDir, outputExt))
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
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: jpegQuality})
	if err != nil {
		logger.Fatal("error encoding image", err)
	}

	return
}

func ConvertDynamic() {
}

func generateOutput(inputPath, outputDir, outputExt string) string {
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputPath)
	inputName := inputBase[0 : len(inputBase)-len(inputExt)]
	return fmt.Sprintf("%s/%s.%s", outputDir, inputName, outputExt)
}
