package helpers

type OpenttdColour uint8

const (
	ColourDarkBlue OpenttdColour = iota
	ColourPaleGreen
	ColourPink
	ColourYellow
	ColourRed
	ColourLightBlue
	ColourGreen
	ColourDarkGreen
	ColourBlue
	ColourCream
	ColourMauve
	ColourPurple
	ColourOrange
	ColourBrown
	ColourGrey
	ColourWhite
)

func (colour OpenttdColour) String() string {
	names := [...]string{
		"Dark Blue",
		"Pale Green",
		"Pink",
		"Yellow",
		"Red",
		"Light Blue",
		"Green",
		"Dark Green",
		"Blue",
		"Cream",
		"Mauve",
		"Purple",
		"Orange",
		"Brown",
		"Grey",
		"White",
	}
	// prevent panics for out of range lookups
	if colour < ColourDarkBlue || colour > ColourWhite {
		return "Unknown"
	}
	return names[colour]
}
