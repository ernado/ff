package ffrun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/zctx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ernado/ff/ffmpeg"
	"github.com/ernado/ff/ffprobe"
)

type Instance struct {
	binary      string // ffmpeg
	binaryProbe string // ffprobe
	trace       trace.Tracer
}

type Progress struct {
	Speed    float64 // 10.8x
	Complete float64 // 0..1
}

type RunOptions struct {
	Input          string // input file path
	Output         string // output file path
	Progress       func(p Progress)
	ProgressPeriod time.Duration
	InputArgs      []string
	Args           []string
	Probe          *ffmpeg.Probe
}

type logBuffer struct {
	limit int
	lines []string
}

func (b *logBuffer) Write(p []byte) (n int, err error) {
	for _, s := range strings.Split(string(p), "\n") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		b.lines = append(b.lines, s)
	}
	if len(b.lines) >= b.limit {
		b.lines = b.lines[1:]
	}

	return len(p), nil
}

type Error struct {
	execErr error
	buffer  *logBuffer
}

func IsInvalidInput(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	for _, s := range []string{
		"Invalid data found when processing input",
		"Non-monotonous DTS in output stream",
	} {
		if e.Contains(s) {
			return true
		}
	}
	return false
}

func (e *Error) Unwrap() error {
	return e.execErr
}

func (e *Error) Contains(s string) bool {
	for _, line := range e.Lines() {
		if strings.Contains(line, s) {
			return true
		}
	}
	return false
}

func (e *Error) Lines() []string {
	if e == nil || e.buffer == nil {
		return nil
	}
	return e.buffer.lines
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil!>"
	}
	if e.buffer == nil || len(e.buffer.lines) == 0 {
		if e.execErr == nil {
			return "<nil???>"
		}
		return e.execErr.Error()
	}
	return fmt.Sprintf("%v: %v", e.execErr, e.Lines())
}

// Run performs requested ffmpeg operation like video encoding.
func (i *Instance) Run(ctx context.Context, opt RunOptions) error {
	ctx, span := i.trace.Start(ctx, "Run")
	defer span.End()

	if opt.Probe == nil {
		probe, err := i.Probe(ctx, opt.Input)
		if err != nil {
			return fmt.Errorf("probe: %w", err)
		}
		opt.Probe = probe
	}
	if opt.ProgressPeriod == 0 {
		opt.ProgressPeriod = time.Second
	}
	summary, err := ffprobe.ParseSummary(opt.Probe)
	if err != nil {
		return fmt.Errorf("summary: %w", err)
	}
	pr, pw := io.Pipe()
	progress := newProgressReader(pr, func(p *ffmpeg.Progress) {
		if opt.Progress == nil {
			return
		}
		var progress Progress
		progress.Complete = p.OutTime.Seconds() / summary.Duration.Seconds()
		opt.Progress(progress)
	})

	args := []string{
		fHideBanner, fOverwrite,
		fVerbose, verboseError,
		fXError,
		fNoStdin,

		fProgress, "pipe:1",
		fStatsPeriod, strconv.FormatFloat(opt.ProgressPeriod.Seconds(), 'f', -1, 64),
	}

	// Set seekable only on http input.
	if strings.HasPrefix(opt.Input, prefixHTTP) {
		args = append(args, fSeekable, "1")
	}

	// Some arguments should be set before input file.
	if len(opt.InputArgs) > 0 {
		args = append(args, opt.InputArgs...)
	}

	args = append(args, fInput, opt.Input)

	// Custom arguments.
	args = append(args, opt.Args...)

	// Output file.
	args = append(args, opt.Output)

	zctx.From(ctx).Info("Running ffmpeg",
		zap.String("ffrun.binary", i.binary),
		zap.Strings("args", args),
	)
	span.SetAttributes(
		attribute.String("ffrun.binary", i.binary),
		attribute.String("input", opt.Input),
		attribute.String("output", opt.Output),
		attribute.StringSlice("args", args),
	)
	cmd := exec.CommandContext(ctx, i.binary, args...)
	logs := &logBuffer{
		limit: 10,
	}
	cmd.Stderr = logs
	cmd.Stdout = pw

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := progress.Run(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				return nil
			}
			return errors.Wrap(err, "progress.Run")
		}
		return nil
	})

	done := make(chan struct{})
	g.Go(func() error {
		select {
		case <-done:
		case <-ctx.Done():
		}
		return pr.Close()
	})
	g.Go(func() error {
		defer close(done)
		return cmd.Run()
	})

	if err := g.Wait(); err != nil {
		span.AddEvent("Error",
			trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("logs", strings.Join(logs.lines, "\n")),
			),
		)
		span.SetStatus(codes.Error, err.Error())
		return &Error{
			execErr: err,
			buffer:  logs,
		}
	}

	return nil
}

// Command line arguments for ffmpeg and ffprobe.
//
// See https://ffmpeg.org/ffprobe-all.html for reference.
const (
	fHideBanner  = "-hide_banner"
	fVerbose     = "-v"
	fInput       = "-i"
	fOverwrite   = "-y"
	fProgress    = "-progress"
	fStatsPeriod = "-stats_period"
	fSeekable    = "-seekable"
	fXError      = "-xerror"
	fNoStdin     = "-nostdin"

	fProbePrintFormat = "-print_format"
	fProbeShowFormat  = "-show_format"
	fProbeShowStreams = "-show_streams"

	printFormatJSON = "json"

	verboseQuiet   = "quiet"
	verbosePanic   = "panic"
	verboseFatal   = "fatal"
	verboseError   = "error"
	verboseWarning = "warning"
	verboseInfo    = "info"
	verboseVerbose = "verbose"
	verboseDebug   = "debug"
	verboseTrace   = "trace"
)

const (
	prefixHTTP = "http://"
)

// Probe performs ffprobe on filePath and returns typed result.
func (i *Instance) Probe(ctx context.Context, filePath string) (*ffmpeg.Probe, error) {
	ctx, span := i.trace.Start(ctx, "Probe")
	defer span.End()

	args := []string{
		fHideBanner,
		fVerbose, verboseError,

		fProbePrintFormat, printFormatJSON,
		fProbeShowFormat,
		fProbeShowStreams,
	}
	if strings.HasPrefix(filePath, prefixHTTP) {
		args = append(args, fSeekable, "1")
	}
	args = append(args, filePath)

	span.SetAttributes(
		attribute.String("ffprobe.binary", i.binaryProbe),
		attribute.String("file.path", filePath),
		attribute.StringSlice("args", args),
	)
	cmd := exec.CommandContext(ctx, i.binaryProbe, args...)

	probeRaw := new(bytes.Buffer)
	cmd.Stdout = probeRaw

	logs := &logBuffer{
		limit: 10,
	}
	cmd.Stderr = logs

	if err := cmd.Run(); err != nil {
		return nil, &Error{
			execErr: err,
			buffer:  logs,
		}
	}

	var probe ffmpeg.Probe
	if err := json.NewDecoder(probeRaw).Decode(&probe); err != nil {
		return nil, errors.Wrap(err, "decode")
	}

	if summary, err := ffprobe.ParseSummary(&probe); err != nil {
		span.AddEvent("ParseSummary",
			trace.WithAttributes(
				attribute.String("error", err.Error()),
			),
		)
	} else {
		span.AddEvent("ParseSummary",
			trace.WithAttributes(
				attribute.Int("streams", len(probe.Streams)),
				attribute.Int("duration_sec", int(summary.Duration.Seconds())),
				attribute.String("duration", summary.Duration.String()),
			),
		)
	}

	return &probe, nil
}

type Options struct {
	Binary         string
	BinaryProbe    string
	TracerProvider trace.TracerProvider
}

func (o *Options) setDefaults() {
	if o.Binary == "" {
		o.Binary = "ffmpeg"
	}
	if o.BinaryProbe == "" {
		o.BinaryProbe = "ffprobe"
	}
	if o.TracerProvider == nil {
		o.TracerProvider = noop.NewTracerProvider()
	}
}

func New(opt Options) *Instance {
	opt.setDefaults()

	i := &Instance{
		binary:      opt.Binary,
		binaryProbe: opt.BinaryProbe,
		trace:       opt.TracerProvider.Tracer("ffrun"),
	}

	return i
}
