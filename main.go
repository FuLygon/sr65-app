package main

import (
	"errors"
	"github.com/ncruces/zenity"
	"os"
	"path/filepath"
	"sr65-app/embed"
	"sr65-app/internal"
	"sr65-app/logger"
	"strings"
)

const (
	outputDir           = "outputs"
	outputExtStatic     = "jpg"
	outputExtDynamic    = "mjpeg"
	outputQualityStatic = 95
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
		zenity.Title("SR65 App"),
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
			{
				Name:     "Media files",
				Patterns: []string{"*.png", "*.jpg", "*.gif", "*.mp4"},
				CaseFold: true,
			},
		},
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			logger.Fatal("user cancelled", err)
		}
		logger.Fatal("an error occurred", err)
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		logger.Fatal("error creating output directory", err)
	}

	// convert input file
	logger.Info("converting file")
	switch strings.ToLower(filepath.Ext(inputPath)) {
	case ".jpg", ".png":
		internal.ConvertStatic(inputPath, outputDir, outputExtStatic, outputQualityStatic)
	case ".gif", ".mp4":
		internal.ConvertDynamic(inputPath, outputDir, outputExtDynamic)
	default:
		logger.Fatal("unsupported file format")
	}

	logger.Info(`converted successfully, output saved in "outputs" directory. Exiting...`)
}
