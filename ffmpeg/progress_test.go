package ffmpeg

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ProgressLine_Parse(t *testing.T) {
	for _, tt := range []struct {
		name    string
		input   string
		output  ProgressLine
		wantErr bool
	}{
		{
			name:  "progress",
			input: "progress=continue",
			output: ProgressLine{
				Key:   "progress",
				Value: "continue",
			},
		},
		{
			name:    "blank",
			input:   "",
			wantErr: true,
		},
		{
			name:    "malformed",
			input:   "foobar:baz",
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var l ProgressLine
			if err := l.Parse(tt.input); tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.output, l)
			}
		})
	}
	t.Run("multiline", func(t *testing.T) {
		in := `progress frame=1040
progress fps=342.68
progress stream_0_0_q=22.0
progress bitrate=3206.8kbits/s
progress total_size=14417964
progress out_time_us=35968000
progress out_time_ms=35968000
progress out_time=00:00:35.968000
progress dup_frames=3
progress drop_frames=0
progress speed=11.9x
progress progress=continue`
		var l ProgressLine
		s := bufio.NewScanner(strings.NewReader(in))
		for s.Scan() {
			if err := l.Parse(s.Text()); err != nil {
				t.Fatal(err)
			}
		}
		require.NoError(t, s.Err())
	})
}
