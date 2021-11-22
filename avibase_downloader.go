package avibase_downloader

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gbdubs/attributions"
)

type AvibaseDownloaderInput struct {
	RegionCode string
}

type AvibaseDownloaderOutput struct {
	Entries     []AvibaseEntry           `xml:"entries"`
	Attribution attributions.Attribution `xml:"attribution"`
}

type AvibaseEntry struct {
	EnglishName string `xml:"english-name"`
	LatinName   string `xml:"latin-name"`
	URL         string `xml:"URL"`
}

const rootUrl = "http://avibase.bsc-eoc.org/"
const checklistUrl = rootUrl + "checklist.jsp"

func (input *AvibaseDownloaderInput) memoizedFileName() string {
	return "/tmp/avibase_downloader/" + input.RegionCode + ".xml"
}

func (input *AvibaseDownloaderInput) readMemoized() (*AvibaseDownloaderOutput, error) {
	output := &AvibaseDownloaderOutput{}
	asBytes, err := ioutil.ReadFile(input.memoizedFileName())
	if err != nil {
		return output, err
	}
	err = xml.Unmarshal(asBytes, output)
	return output, err
}

func (input *AvibaseDownloaderInput) writeMemoized(output *AvibaseDownloaderOutput) error {
	err := os.MkdirAll(filepath.Dir(input.memoizedFileName()), 0777)
	if err != nil {
		return err
	}
	asBytes, err := xml.MarshalIndent(*output, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(input.memoizedFileName(), asBytes, 0777)
}

func (input *AvibaseDownloaderInput) Execute() (*AvibaseDownloaderOutput, error) {
	output, err := input.readMemoized()
	if err == nil {
		return output, err
	}
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
					entry.URL = rootUrl + partialUrl
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
	err = input.writeMemoized(output)
	if err != nil {
		return output, err
	}
	return output, nil
}
