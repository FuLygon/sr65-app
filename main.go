package main

import (
	"errors"
	"github.com/ncruces/zenity"
	_ "image/png"
	"os"
	"path/filepath"
	"sr65-software/embed"
	"sr65-software/internal"
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
		internal.ConvertStatic(inputPath, outputDir, outputStaticExt, outputStaticQuality)
	case ".gif", ".mp4":
		internal.ConvertDynamic()
	default:
		logger.Fatal("unsupported file format")
	}

	logger.Info(`converted successfully, output saved in "outputs" directory. Exiting...`)
}
