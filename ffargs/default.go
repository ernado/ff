package ffargs

// Default ffmpeg args.
func Default() []string {
	return []string{
		"-y", // replace
		"-bsf:a", "aac_adtstoasc",
		"-c", "copy",
	}
}
