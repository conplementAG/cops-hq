package regions

var azureRegionAbbreviations = map[Region]string{
	// naming convention is region code (e.g. n as in north, one letter) and country code (e.g. eu as in europe, two letter)
	NorthEurope:   "neu",
	WestEurope:    "weu",
	FranceCentral: "cfr",
	EastUs:        "eus",
	WestUs:        "wus",
	CentralUs:     "cus",
	CanadaEast:    "eca",
}

func GetAbbreviatedRegion(region Region) string {
	abbreviatedRegion, regionSupported := azureRegionAbbreviations[region]

	if !regionSupported {
		panic("This region (" + region + ") is not supported in our naming convention mapper yet. Please extend the mapper.")
	}

	return abbreviatedRegion
}
