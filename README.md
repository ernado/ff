# ff [![Go Reference](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/ernado/ff#section-documentation) [![codecov](https://img.shields.io/codecov/c/github/ernado/ff?label=cover)](https://codecov.io/gh/ernado/ff) [![experimental](https://img.shields.io/badge/-experimental-blueviolet)](https://ernado.org/docs/projects/status#experimental)


Helpers for running ffmpeg commands in Go.

## Example

You need ffmpeg in PATH to use this package.
Alternatively, you can provide path to binary in `ffrun.Options`.

```go
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

```
