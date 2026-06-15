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
	NonDiegetic     int `json:"non_diegetic"`
	Captions        int `json:"captions"`
	Descriptions    int `json:"descriptions"`
	Metadata        int `json:"metadata"`
	Dependent       int `json:"dependent"`
	StillImage      int `json:"still_image"`
	Multilayer      int `json:"multilayer"`
}

// ProbeSideData is a single entry of a stream's side_data_list.
type ProbeSideData struct {
	SideDataType  string `json:"side_data_type"`
	DisplayMatrix string `json:"displaymatrix,omitempty"`
	Rotation      int    `json:"rotation,omitempty"`
}

type ProbeStream struct {
	Index              int              `json:"index"`
	CodecName          string           `json:"codec_name,omitempty"`
	CodecLongName      string           `json:"codec_long_name,omitempty"`
	Profile            string           `json:"profile,omitempty"`
	CodecType          string           `json:"codec_type"`
	CodecTimeBase      string           `json:"codec_time_base,omitempty"`
	CodecTagString     string           `json:"codec_tag_string"`
	CodecTag           string           `json:"codec_tag"`
	Width              int              `json:"width,omitempty"`
	Height             int              `json:"height,omitempty"`
	CodedWidth         int              `json:"coded_width,omitempty"`
	CodedHeight        int              `json:"coded_height,omitempty"`
	ClosedCaptions     int              `json:"closed_captions,omitempty"`
	FilmGrain          int              `json:"film_grain,omitempty"`
	HasBFrames         int              `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string           `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string           `json:"display_aspect_ratio,omitempty"`
	PixFmt             string           `json:"pix_fmt,omitempty"`
	Level              int              `json:"level,omitempty"`
	ColorRange         string           `json:"color_range,omitempty"`
	ColorSpace         string           `json:"color_space,omitempty"`
	ColorTransfer      string           `json:"color_transfer,omitempty"`
	ColorPrimaries     string           `json:"color_primaries,omitempty"`
	ChromaLocation     string           `json:"chroma_location,omitempty"`
	FieldOrder         string           `json:"field_order,omitempty"`
	Refs               int              `json:"refs,omitempty"`
	IsAVC              string           `json:"is_avc,omitempty"`
	NalLengthSize      string           `json:"nal_length_size,omitempty"`
	ID                 string           `json:"id,omitempty"`
	RFrameRate         string           `json:"r_frame_rate"`
	AvgFrameRate       string           `json:"avg_frame_rate"`
	TimeBase           string           `json:"time_base"`
	StartPts           int64            `json:"start_pts"`
	StartTime          string           `json:"start_time"`
	DurationTs         int64            `json:"duration_ts,omitempty"`
	Duration           string           `json:"duration"`
	BitRate            string           `json:"bit_rate,omitempty"`
	BitsPerRawSample   string           `json:"bits_per_raw_sample,omitempty"`
	NbFrames           string           `json:"nb_frames,omitempty"`
	ExtradataSize      int              `json:"extradata_size,omitempty"`
	Disposition        ProbeDisposition `json:"disposition"`
	Tags               ProbeTags        `json:"tags"`
	SideDataList       []ProbeSideData  `json:"side_data_list,omitempty"`
	SampleFmt          string           `json:"sample_fmt,omitempty"`
	SampleRate         string           `json:"sample_rate,omitempty"`
	Channels           int              `json:"channels,omitempty"`
	ChannelLayout      string           `json:"channel_layout,omitempty"`
	BitsPerSample      int              `json:"bits_per_sample,omitempty"`
	InitialPadding     int              `json:"initial_padding,omitempty"`
}

type ProbeFormat struct {
	Filename       string    `json:"filename"`
	NbStreams      int       `json:"nb_streams"`
	NbPrograms     int       `json:"nb_programs"`
	NbStreamGroups int       `json:"nb_stream_groups"`
	FormatName     string    `json:"format_name"`
	FormatLongName string    `json:"format_long_name"`
	StartTime      string    `json:"start_time"`
	Duration       string    `json:"duration"`
	Size           string    `json:"size"`
	BitRate        string    `json:"bit_rate"`
	ProbeScore     int       `json:"probe_score"`
	Tags           ProbeTags `json:"tags"`
}

// ProbeTags holds metadata tags reported for a stream or container format.
//
// ffprobe emits an open-ended set of tags; the fields below cover the common
// ones shared across formats (notably MP4/MOV). Unmodeled tags are ignored.
type ProbeTags struct {
	Title            string `json:"title,omitempty"`
	Encoder          string `json:"encoder,omitempty"`
	CreationTime     string `json:"creation_time,omitempty"`
	Language         string `json:"language,omitempty"`
	HandlerName      string `json:"handler_name,omitempty"`
	VendorID         string `json:"vendor_id,omitempty"`
	Timecode         string `json:"timecode,omitempty"`
	MajorBrand       string `json:"major_brand,omitempty"`
	MinorVersion     string `json:"minor_version,omitempty"`
	CompatibleBrands string `json:"compatible_brands,omitempty"`
}

type Probe struct {
	Streams []ProbeStream `json:"streams"`
	Format  ProbeFormat   `json:"format"`

	// Raw value returned by ffmpeg.
	Raw []byte `json:"-"`
}
