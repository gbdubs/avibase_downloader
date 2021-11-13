package avibase_downloader

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gbdubs/attributions"
)

type AvibaseDownloaderInput struct {
	RegionCode string
}

type AvibaseDownloaderOutput struct {
	Entries     []AvibaseEntry
	Attribution attributions.Attribution
}

type AvibaseEntry struct {
	EnglishName string
	LatinName   string
	Url         string
}

const rootUrl = "http://avibase.bsc-eoc.org/"
const checklistUrl = rootUrl + "checklist.jsp"

func (input AvibaseDownloaderInput) Execute() (AvibaseDownloaderOutput, error) {
	output := AvibaseDownloaderOutput{}
	req, err := http.NewRequest("GET", checklistUrl, nil)
	if err != nil {
		return output, err
	}
	q := req.URL.Query()
	q.Add("region", input.RegionCode)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return output, errors.New(fmt.Sprintf("Request Failed: %d, %s", resp.StatusCode, resp.Status))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return output, err
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("td").Length() == 3 {
			entry := AvibaseEntry{}
			s.Find("td").Each(func(j int, ss *goquery.Selection) {
				if j == 0 {
					entry.EnglishName = ss.Text()
				} else if j == 1 {
					entry.LatinName = ss.Find("i").Text()
					partialUrl, _ := ss.Find("a").Attr("href")
					entry.Url = rootUrl + partialUrl
				}
			})
			output.Entries = append(output.Entries, entry)
		}
	})
	output.Attribution = attributions.Attribution{
		OriginUrl:           req.URL.String(),
		CollectedAt:         time.Now(),
		OriginalTitle:       doc.Find("title").Text(),
		Author:              "Avibase - Denis LePage",
		AuthorUrl:           rootUrl,
		ScrapingMethodology: "github.com/gbdubs/avibase_downloader",
		Context:             []string{"Scraped the Avibase Website to list the set of birds that can be found in a given region."},
	}
	return output, nil
}
