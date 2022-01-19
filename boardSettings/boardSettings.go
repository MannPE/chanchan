package boardSettings

// Note: All keywords have to be lowercase - text is not case sensitive yet
func GetSFW4chanBoards() map[string][]string {
	general := [] string { 
		"wallpapers",
		"kep1er",
		"konosuba",
		"twice",
	}

	animated := [] string {
		"tengen",
		"asuka",
		"steins",
		"poke",
		"emblem",
	}

	real := []string {
		"mexico",
		"throwback",
		"citiscape",
		"flower",
	}


	defaultSettings := map[string] []string {
		"global": general,
		"p": real,
		"w": animated,
		"v": animated,
	}

	return defaultSettings
}