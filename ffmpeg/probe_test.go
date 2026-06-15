package ffmpeg

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// movH264Probe is the output of:
//
//	ffprobe -hide_banner -v error -print_format json -show_format -show_streams \
//		big_buck_bunny_480p_h264.mov
//
// produced by ffmpeg 7.1.
//
//go:embed _testdata/mov_h264_probe.json
var movH264Probe []byte

// strictDecode ensures every JSON key in data is covered by a struct tag in v.
func strictDecode(t *testing.T, data []byte, v any) {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	require.NoError(t, dec.Decode(v))
}

func TestProbeDecode(t *testing.T) {
	var probe Probe
	require.NoError(t, json.Unmarshal(movH264Probe, &probe))

	require.Len(t, probe.Streams, 3)

	// Every key emitted for the streams must be covered by struct tags: the
	// raw stream objects are re-decoded with DisallowUnknownFields. (The
	// container's open-ended tags map holds proprietary com.apple.* keys and
	// is intentionally left to normal decoding.)
	var raw struct {
		Streams []json.RawMessage `json:"streams"`
	}
	require.NoError(t, json.Unmarshal(movH264Probe, &raw))
	require.Len(t, raw.Streams, 3)
	for _, s := range raw.Streams {
		var st ProbeStream
		strictDecode(t, s, &st)
	}

	video := probe.Streams[0]
	require.Equal(t, "video", video.CodecType)
	require.Equal(t, "h264", video.CodecName)
	require.Equal(t, "Main", video.Profile)
	require.Equal(t, 853, video.Width)
	require.Equal(t, 480, video.Height)
	require.Equal(t, "yuv420p", video.PixFmt)
	require.Equal(t, "bt709", video.ColorSpace)
	require.Equal(t, "bt709", video.ColorTransfer)
	require.Equal(t, "bt709", video.ColorPrimaries)
	require.Equal(t, "tv", video.ColorRange)
	require.Equal(t, "topleft", video.ChromaLocation)
	require.Equal(t, "true", video.IsAVC)
	require.Equal(t, "596.458333", video.Duration)
	require.Equal(t, int64(1431500), video.DurationTs)
	require.Equal(t, "2899884", video.BitRate)
	require.Equal(t, "14315", video.NbFrames)
	require.Equal(t, "8", video.BitsPerRawSample)
	require.Equal(t, "0x1", video.ID)
	require.Equal(t, 1, video.Disposition.Default)

	// Stream-level tags.
	require.Equal(t, "eng", video.Tags.Language)
	require.Equal(t, "Apple Video Media Handler", video.Tags.HandlerName)
	require.Equal(t, "appl", video.Tags.VendorID)
	require.Equal(t, "2008-05-27T18:32:32.000000Z", video.Tags.CreationTime)

	// Side data list (display matrix / rotation).
	require.Len(t, video.SideDataList, 1)
	require.Equal(t, "Display Matrix", video.SideDataList[0].SideDataType)

	audio := probe.Streams[2]
	require.Equal(t, "audio", audio.CodecType)
	require.Equal(t, "aac", audio.CodecName)
	require.Equal(t, 6, audio.Channels)
	require.Equal(t, "5.1", audio.ChannelLayout)
	require.Equal(t, "48000", audio.SampleRate)

	// Format.
	require.Equal(t, 3, probe.Format.NbStreams)
	require.Equal(t, 0, probe.Format.NbStreamGroups)
	require.Equal(t, "mov,mp4,m4a,3gp,3g2,mj2", probe.Format.FormatName)
	require.Equal(t, "596.461667", probe.Format.Duration)
	require.Equal(t, "qt  ", probe.Format.Tags.MajorBrand)
}
