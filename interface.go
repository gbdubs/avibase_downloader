package avibase_downloader

import (
	"github.com/gbdubs/attributions"
)

type Input struct {
	RegionCodes []string
	IncludeRare bool
	ForceReload bool
}

type Output struct {
	Entries      []AvibaseEntry             `xml:"entries"`
	Attributions []attributions.Attribution `xml:"attribution"`
}

type AvibaseEntry struct {
	EnglishName string `xml:"english-name"`
	LatinName   string `xml:"latin-name"`
	URL         string `xml:"URL"`
}
