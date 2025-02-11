package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-faster/errors"

	"github.com/ernado/ff/ffprobe"
	"github.com/ernado/ff/ffrun"
)

func run() error {
	ctx := context.Background()
	i := ffrun.New(ffrun.Options{})
	filePath := filepath.Join("ffrun", "_testdata", "bbb.mp4")

	// Probe.
	probe, err := i.Probe(ctx, filePath)
	if err != nil {
		return errors.Wrap(err, "probe")
	}

	summary, err := ffprobe.ParseSummary(probe)
	if err != nil {
		return errors.Wrap(err, "summary")
	}

	fmt.Println("Duration:", summary.Duration)

	// Encode.
	tempDir, err := os.MkdirTemp("", "output")
	if err != nil {
		return errors.Wrap(err, "temp dir")
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
		fmt.Println("cleaned up")
	}()

	outputPath := filepath.Join(tempDir, "output.mp4")
	if err := i.Run(ctx, ffrun.RunOptions{
		Input:  filePath,
		Output: outputPath,
		Progress: func(p ffrun.Progress) {
			fmt.Println("progress:", p.Complete)
		},
		Args: []string{
			"-t", "5",
			"-ac", "2",
		},
	}); err != nil {
		return errors.Wrap(err, "run")
	}
	fmt.Println("done")
	return nil
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}
