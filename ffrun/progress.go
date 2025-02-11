package ffrun

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ernado/ff/ffmpeg"
)

type ProgressListener struct {
	path     string
	listener net.Listener
	mux      *http.ServeMux
	srv      *http.Server
	f        func(p *ffmpeg.Progress)
}

func (p *ProgressListener) Addr() string {
	u := &url.URL{
		Scheme: "http",
		Path:   p.path,
		Host:   p.listener.Addr().String(),
	}
	return u.String()
}

func (p *ProgressListener) registerHandler() {
	p.mux.HandleFunc(p.path, func(w http.ResponseWriter, r *http.Request) {
		s := bufio.NewScanner(r.Body)
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
			case "speed":
				n, err := strconv.ParseFloat(line.Value[:len(line.Value)-1], 64)
				if err != nil {
					continue
				}
				pr.Speed = n
			case "progress":
				pr.Done = line.Value == "done"
				p.f(pr)
			}
		}
	})
}

func (p *ProgressListener) Stop() error {
	defer func() {
		_ = p.listener.Close()
	}()
	return p.srv.Close()
}

func (p *ProgressListener) Run() error {
	if err := p.srv.Serve(p.listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewProgressListener(f func(p *ffmpeg.Progress)) (*ProgressListener, error) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	mux := http.NewServeMux()
	p := &ProgressListener{
		path:     "/progress",
		listener: ln,
		mux:      mux,
		f:        f,

		srv: &http.Server{
			Handler: mux,
		},
	}
	p.registerHandler()
	return p, nil
}
