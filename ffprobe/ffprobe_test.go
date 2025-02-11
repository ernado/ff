package ffprobe

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ernado/ff/ffmpeg"
)

//go:embed _testdata/sample_probe.json
var sampleProbe []byte

func TestNewSummary(t *testing.T) {
	var probe ffmpeg.Probe
	require.NoError(t, json.Unmarshal(sampleProbe, &probe))

	summary, err := ParseSummary(&probe)
	require.NoError(t, err)

	t.Logf("Summary: %+v", summary)
}
