package main

import (
	"errors"
	"fmt"
	"github.com/ncruces/zenity"
	"os"
	"os/exec"
	"path/filepath"
	"sr65-app/embed"
	"sr65-app/internal"
	"sr65-app/logger"
	"strings"
)

const (
	// output
	outputDir           = "outputs"
	outputExtStatic     = "jpg"
	outputExtDynamic    = "mjpeg"
	outputQualityStatic = 95

	// zenity
	zenityTitle = "SR65 App"
	zenityWidth = 400
)

var ffmpegInstalled bool

func init() {
	// check if ffmpeg is installed
	_, err := exec.LookPath("ffmpeg")
	if err == nil {
		ffmpegInstalled = true
	}
}

func main() {
	// embedding binaries
	tmpDir, err := embed.ExtractBinaries()
	if err != nil {
		logger.Warn("error extracting embedded binaries, falling back to system binaries", err)
	}
	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			logger.Error("error removing temporary directory", err)
		}
	}(tmpDir)

	// show warning dialog if ffmpeg is not installed
	if !ffmpegInstalled {
		err = zenity.Warning(
			"ffmpeg command not found, certain features will be unavailable.",
			zenity.Title(zenityTitle),
			zenity.Width(zenityWidth),
		)
		if err != nil {
			handleZenityCancelErr(err)
			return
		}
	}

	// show question dialog
	err = zenity.Question(
		"Choose an image or video file. Supported file formats:"+"\n"+
			"- Image: png, jpg/jpeg, gif"+"\n"+
			"- Video: mp4",
		zenity.Width(zenityWidth),
		zenity.Title(zenityTitle),
		zenity.OKLabel("Open File..."),
		zenity.CancelLabel("Close"),
	)
	if err != nil {
		handleZenityCancelErr(err)
		return
	}

	// show file picker
	inputPath, err := zenity.SelectFile(
		zenity.Title("Select a file"),
		zenity.FileFilters{
			{
				Name:     "Media files",
				Patterns: []string{"*.png", "*.jpg", "*.jpeg", "*.gif", "*.mp4"},
				CaseFold: true,
			},
		},
	)
	if err != nil {
		handleZenityCancelErr(err)
		return
	}

	// create output directory
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		logger.Error("error creating output directory", err)
		return
	}

	// convert input file
	switch strings.ToLower(filepath.Ext(inputPath)) {
	case ".jpg", ".jpeg", ".png":
		err = internal.ConvertStatic(inputPath, outputDir, outputExtStatic, outputQualityStatic)
		if err != nil {
			logger.Error("error converting static media", err)
		}
	case ".gif", ".mp4":
		err = internal.ConvertDynamic(inputPath, outputDir, outputExtDynamic, tmpDir)
		if err != nil {
			logger.Error("error converting dynamic media", err)
		}
	default:
		err = fmt.Errorf("unsupported format: %s", filepath.Ext(inputPath))
		logger.Error("unsupported format", err)
	}
	if err != nil {
		return
	}

	logger.Info("converted successfully, output saved in 'outputs' directory. Exiting...")
}

func handleZenityCancelErr(err error) {
	if errors.Is(err, zenity.ErrCanceled) {
		logger.Warn("user cancelled", err)
		return
	}
	logger.Error("an error occurred", err)
}
