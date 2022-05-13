package regions

var azureRegionAbbreviations = map[string]string{
	// naming convention is region code (e.g. n as in north, one letter) and country code (e.g. eu as in europe, two letter)
	"northeurope":   "neu",
	"westeurope":    "weu",
	"francecentral": "cfr",
	"eastus":        "eus",
	"westus":        "wus",
	"centralus":     "cus",
	"canadaeast":    "eca",
}

func GetAbbreviatedRegion(region string) string {
	abbreviatedRegion, regionSupported := azureRegionAbbreviations[region]

	if !regionSupported {
		panic("This region (" + region + ") is not supported in our naming convention mapper yet. Please extend the mapper.")
	}

	return abbreviatedRegion
}
