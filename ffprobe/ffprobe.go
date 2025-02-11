// Package ffprobe wraps helpers for working with ffprobe.
package ffprobe

import (
	"math"
	"strconv"
	"time"

	"github.com/ernado/ff/ffmpeg"
)

// Summary of probe result.
type Summary struct {
	Duration time.Duration
	HasVideo bool
	HasAudio bool

	Width  int
	Height int
}

func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		// No duration.
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	sec, dec := math.Modf(f)
	return time.Second*time.Duration(sec) + time.Nanosecond*time.Duration(dec*(1e9)), nil
}

// ParseSummary parses Summary from raw ffmpeg.Probe.
func ParseSummary(probe *ffmpeg.Probe) (*Summary, error) {
	var s Summary
	for _, st := range probe.Streams {
		if st.Duration != "" {
			duration, err := parseDuration(st.Duration)
			if err != nil {
				return nil, err
			}
			if duration > s.Duration {
				s.Duration = duration
			}
		}
		switch st.CodecType {
		case "video":
			s.HasVideo = true
			s.Width = st.Width
			s.Height = st.Height
		case "audio":
			s.HasAudio = true
		}
	}
	if s.Duration == 0 {
		duration, err := parseDuration(probe.Format.Duration)
		if err != nil {
			return nil, err
		}
		s.Duration = duration
	}

	return &s, nil
}
