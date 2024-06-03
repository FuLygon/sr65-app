package internal

import (
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sr65-app/logger"
)

// false for debugging ffmpeg
const ffmpegSilent = true

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
}

func ConvertDynamic(inputPath, outputDir, outputExt string) {
	var ffmpegInput = inputPath

	// convert gif to avi if input is gif
	if filepath.Ext(inputPath) == ".gif" {
		// create temporary directory for storing avi file
		tmpDir, err := os.MkdirTemp("", "sr65-app-*")
		if err != nil {
			logger.Fatal("error creating temporary directory for gif conversion", err)
		}
		defer os.RemoveAll(tmpDir)

		// generate output avi path
		outputFileAvi := generateOutput(inputPath, tmpDir, "avi")

		// convert gif to avi
		err = ffmpeg.Input(inputPath).Output(outputFileAvi,
			ffmpeg.KwArgs{
				"f":   "gif",
				"vf":  "fps=31,scale=128:128:flags=lanczos",
				"q:v": 1,
			}).
			OverWriteOutput().
			Silent(ffmpegSilent).
			Run()
		if err != nil {
			logger.Fatal("error converting gif to avi", err)
		}

		// set ffmpegInput to avi path
		ffmpegInput = outputFileAvi
	}

	// generate output path
	outputFile := generateOutput(inputPath, outputDir, outputExt)

	// convert
	err := ffmpeg.Input(ffmpegInput).Output(outputFile,
		ffmpeg.KwArgs{
			"vf":  "fps=31,scale=128:128:flags=lanczos",
			"q:v": 1,
		}).
		OverWriteOutput().
		Silent(ffmpegSilent).
		Run()
	if err != nil {
		logger.Fatal("error converting video to "+outputExt, err)
	}
}

func generateOutput(inputPath, outputDir, outputExt string) string {
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputPath)
	inputName := inputBase[0 : len(inputBase)-len(inputExt)]
	return filepath.Join(outputDir, fmt.Sprintf("%s.%s", inputName, outputExt))
}
