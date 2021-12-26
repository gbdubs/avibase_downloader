package bird_region_rosters

import (
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
	"github.com/gbdubs/verbose"
)

type Input struct {
	RegionCodes []string
	IncludeRare bool
	ForceReload bool
	verbose.Verbose
}

type Output struct {
	Entries      []bird.BirdName            `xml:"entries"`
	Attributions []attributions.Attribution `xml:"attribution"`
}
