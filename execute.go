package avibase_downloader

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gbdubs/attributions"
)

const rootUrl = "http://avibase.bsc-eoc.org/"
const checklistUrl = rootUrl + "checklist.jsp"

func (input *Input) Execute() (output *Output, err error) {
	output = &Output{}
	if !input.ForceReload {
		output, err = input.readMemoized()
		if err == nil {
			return
		}
	}
	for _, regionCode := range input.RegionCodes {
		e, a, er := executeForRegion(regionCode, input.IncludeRare)
		if er != nil {
			err = er
			return
		}
		output.Entries = append(output.Entries, e...)
		output.Attributions = append(output.Attributions, a)
	}
	err = input.writeMemoized(output)
	if err != nil {
		err = fmt.Errorf("memoization failed: %v", err)
	}
	return
}

func executeForRegion(regionCode string, includeRare bool) (entries []AvibaseEntry, attribution attributions.Attribution, err error) {
	req, err := http.NewRequest("GET", checklistUrl, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("region", regionCode)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("request failed: %v", err)
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request failed: %d %s", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		err = fmt.Errorf("parse document failed: %v", err)
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td").Length() == 3 {
			entry := AvibaseEntry{}
			isRare := false
			s.Find("td").Each(func(j int, ss *goquery.Selection) {
				if j == 0 {
					entry.EnglishName = ss.Text()
				} else if j == 1 {
					entry.LatinName = ss.Find("i").Text()
					partialUrl, _ := ss.Find("a").Attr("href")
					entry.URL = rootUrl + partialUrl
				} else if j == 2 {
					if strings.Contains(ss.Text(), "Rare") {
						isRare = true
					}
				}
			})
			if !isRare || includeRare {
				entries = append(entries, entry)
			}
		}
	})
	attribution = attributions.Attribution{
		OriginUrl:           req.URL.String(),
		CollectedAt:         time.Now(),
		OriginalTitle:       doc.Find("title").Text(),
		Author:              "Avibase - Denis LePage",
		AuthorUrl:           rootUrl,
		ScrapingMethodology: "github.com/gbdubs/avibase_downloader",
		Context:             []string{"Scraped the Avibase Website to list the set of birds that can be found in a given region."},
	}
	return
}
