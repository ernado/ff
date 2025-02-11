package ffargs

// Variant is a ffmpeg argument set variant.
type Variant struct {
	Name string
	Args []string
}

// All returns all variants.
func All() []Variant {
	return []Variant{
		{
			Name: "default",
			Args: Default(),
		},
		{
			Name: "blank",
			Args: []string{},
		},
	}
}
