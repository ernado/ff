package ffmpeg

type ProbeDisposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

type ProbeStream struct {
	Index              int              `json:"index"`
	CodecName          string           `json:"codec_name"`
	CodecLongName      string           `json:"codec_long_name"`
	Profile            string           `json:"profile,omitempty"`
	CodecType          string           `json:"codec_type"`
	CodecTimeBase      string           `json:"codec_time_base"`
	CodecTagString     string           `json:"codec_tag_string"`
	CodecTag           string           `json:"codec_tag"`
	Width              int              `json:"width,omitempty"`
	Height             int              `json:"height,omitempty"`
	CodedWidth         int              `json:"coded_width,omitempty"`
	CodedHeight        int              `json:"coded_height,omitempty"`
	HasBFrames         int              `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string           `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string           `json:"display_aspect_ratio,omitempty"`
	PixFmt             string           `json:"pix_fmt,omitempty"`
	Level              int              `json:"level,omitempty"`
	FieldOrder         string           `json:"field_order,omitempty"`
	Refs               int              `json:"refs,omitempty"`
	RFrameRate         string           `json:"r_frame_rate"`
	AvgFrameRate       string           `json:"avg_frame_rate"`
	TimeBase           string           `json:"time_base"`
	StartPts           int              `json:"start_pts"`
	StartTime          string           `json:"start_time"`
	Duration           string           `json:"duration"`
	Disposition        ProbeDisposition `json:"disposition"`
	SampleFmt          string           `json:"sample_fmt,omitempty"`
	SampleRate         string           `json:"sample_rate,omitempty"`
	Channels           int              `json:"channels,omitempty"`
	ChannelLayout      string           `json:"channel_layout,omitempty"`
	BitsPerSample      int              `json:"bits_per_sample,omitempty"`
	Tags               StreamTags       `json:"tags,omitempty"`
}

type StreamTags struct {
	Langeuate string `json:"langeuage,omitempty"`
	Duration  string `json:"DURATION,omitempty"`
}
 
type ProbeFormat struct {
	Filename       string    `json:"filename"`
	NbStreams      int       `json:"nb_streams"`
	NbPrograms     int       `json:"nb_programs"`
	FormatName     string    `json:"format_name"`
	FormatLongName string    `json:"format_long_name"`
	StartTime      string    `json:"start_time"`
	Duration       string    `json:"duration"`
	Size           string    `json:"size"`
	BitRate        string    `json:"bit_rate"`
	ProbeScore     int       `json:"probe_score"`
	Tags           ProbeTags `json:"tags"`
}

type ProbeTags struct {
	Title   string `json:"title"`
	Encoder string `json:"encoder"`
}

type Probe struct {
	Streams []ProbeStream `json:"streams"`
	Format  ProbeFormat   `json:"format"`

	// Raw value returned by ffmpeg.
	Raw []byte `json:"-"`
}
