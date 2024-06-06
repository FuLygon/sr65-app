package internal

import (
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"image"
	"image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sr65-app/logger"
)

// false for debugging ffmpeg
const ffmpegSilent = true

func ConvertStatic(inputPath, outputDir, outputExt string, jpegQuality int) error {
	// open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Error("error closing input file", err)
		}
	}(inputFile)

	// decode image
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// resizing image
	img = imaging.Resize(img, 128, 128, imaging.Lanczos)

	// create output file
	outputFile, err := os.Create(generateOutput(inputPath, outputDir, outputExt))
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Error("error closing output file", err)
		}
	}(outputFile)

	// encode img to output file
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: jpegQuality})
	if err != nil {
		return fmt.Errorf("error encoding image: %w", err)
	}

	return nil
}

func ConvertDynamic(inputPath, outputDir, outputExt, tmpDir string) error {
	var ffmpegInput = inputPath

	// convert gif to avi if input is gif
	if filepath.Ext(inputPath) == ".gif" {
		// generate output avi path
		outputFileAvi := generateOutput(inputPath, tmpDir, "avi")

		// convert gif to avi
		err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"f": "gif"}).Output(outputFileAvi,
			ffmpeg.KwArgs{
				"vf":  "fps=31,scale=128:128:flags=lanczos",
				"q:v": 1,
			}).
			OverWriteOutput().
			Silent(ffmpegSilent).
			Run()
		if err != nil {
			return fmt.Errorf("error converting gif to avi: %w", err)
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
		return fmt.Errorf("error converting video to %s: %w", outputExt, err)
	}

	return nil
}

func ConvertGif(inputPath, outputDir, outputExt string, jpegQuality int) error {
	// open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Error("error closing input file", err)
		}
	}(inputFile)

	// decode gif
	gifData, err := gif.DecodeAll(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding gif: %w", err)
	}

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

	// loop through gif frames
	for _, frame := range gifData.Image {
		// resizing frame
		img := imaging.Resize(frame, 128, 128, imaging.Lanczos)

		// encode img to output file
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: jpegQuality})
		if err != nil {
			return fmt.Errorf("error encoding image: %w", err)
		}
	}

	return nil
}

func generateOutput(inputPath, outputDir, outputExt string) string {
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputPath)
	inputName := inputBase[0 : len(inputBase)-len(inputExt)]
	return filepath.Join(outputDir, fmt.Sprintf("%s.%s", inputName, outputExt))
}
