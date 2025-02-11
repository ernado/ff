package ffmpeg

import (
	"errors"
	"strings"
	"time"
)

type Progress struct {
	Done    bool
	Speed   float64 // 10.8x
	OutTime time.Duration
}

// ProgressLine is `k=v` line of ffmpeg progress report.
type ProgressLine struct {
	Key   string
	Value string
}

func (l *ProgressLine) Parse(s string) error {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return errors.New("bad format")
	}
	l.Key = parts[0]
	l.Value = parts[1]

	return nil
}
