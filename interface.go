package avibase_downloader

import (
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
)

type Input struct {
	RegionCodes []string
	IncludeRare bool
	ForceReload bool
}

type Output struct {
	Entries      []bird.BirdName            `xml:"entries"`
	Attributions []attributions.Attribution `xml:"attribution"`
}
