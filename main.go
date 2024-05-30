package main

import (
	"errors"
	"fmt"
	"github.com/ncruces/zenity"
	"sr65-software/embed"
	"sr65-software/logger"
)

func main() {
	// embedding binaries
	cleanup, err := embed.ExtractBinaries()
	if err != nil {
		logger.Warn("error extracting embedded binaries, falling back to system ffmpeg")
	}
	defer cleanup()

	// question dialog
	err = zenity.Question(
		"",
		zenity.Title("SR65 Software"),
		zenity.OKLabel("Choose file..."),
		zenity.CancelLabel("Close"),
	)
	if err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			logger.Fatal("user cancelled", err)
		}
		logger.Fatal("an error occurred", err)
	}

	// file picker
	filename, err := zenity.SelectFile(
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
	fmt.Println("Selected file path:", filename)
}
