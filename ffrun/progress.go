package ffrun

import (
	"bufio"
	"io"
	"strconv"
	"time"

	"github.com/ernado/ff/ffmpeg"
)

type progressReader struct {
	r io.Reader
	f func(p *ffmpeg.Progress)
}

func (r progressReader) Run() error {
	s := bufio.NewScanner(r.r)
	pr := &ffmpeg.Progress{}
	for s.Scan() {
		var line ffmpeg.ProgressLine
		if err := line.Parse(s.Text()); err != nil {
			continue
		}
		switch line.Key {
		case "out_time_us":
			n, err := strconv.ParseInt(line.Value, 10, 64)
			if err != nil {
				continue
			}
			pr.OutTime = time.Duration(n) * time.Microsecond
		case "progress":
			pr.Done = line.Value == "done"
			r.f(pr)
			pr = &ffmpeg.Progress{}
		}
	}
	return s.Err()
}

func newProgressReader(r io.Reader, p func(p *ffmpeg.Progress)) progressReader {
	return progressReader{
		r: r,
		f: p,
	}
}
