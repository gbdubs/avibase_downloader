package bird_region_rosters

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
)

const (
	avibaseSiteKey                 = "avibase_synonyms"
	avibaseMaxConcurrentRequests   = 2
	wikipediaSiteKey               = "wikipedia_synonyms"
	wikipediaMaxConcurrentRequests = 4
	avibaseSynonymsPage            = "synonyms"
)

func (a avibaseEntry) getSynonymsRequests() []*amass.GetRequest {
	avibaseRequest := &amass.GetRequest{
		Site:                      avibaseSiteKey,
		RequestKey:                a.LatinName,
		URL:                       fmt.Sprintf("https://avibase.bsc-eoc.org/species.jsp?lang=EN&sec=%s&avibaseid=%s", avibaseSynonymsPage, a.AvibaseId),
		SiteMaxConcurrentRequests: avibaseMaxConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "Avibase - Denis LePage",
			AuthorUrl:           "https://avibase.bsc-eoc.org/",
			ScrapingMethodology: "github.com/gbdubs/avibase_downloader",
		},
	}
	avibaseRequest.SetRoundTripData(a)

	v := url.Values{}
	v.Add("action", "query")
	v.Add("prop", "redirects")
	v.Add("titles", a.EnglishName)
	v.Add("redirects", "1") // Means we will get back the canonical page if we give a non-dominant spelling/casing.
	v.Add("rdprop", "title")
	v.Add("rdlimit", "500")
	v.Add("format", "xml")
	wurl := "https://en.wikipedia.org/w/api.php?" + v.Encode()
	wikipediaRequest := &amass.GetRequest{
		Site:                      wikipediaSiteKey,
		RequestKey:                a.EnglishName,
		URL:                       wurl,
		SiteMaxConcurrentRequests: wikipediaMaxConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "Wikipedia Foundation, Inc.",
			AuthorUrl:           "https://wikipedia.org",
			License:             "Creative Commons Attribution-ShareAlike 3.0 Unported License (CC BY-SA)",
			LicenseUrl:          "https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License",
			ScrapingMethodology: "github.com/gbdubs/avibase_downloader/wikipedia",
		},
	}
	wikipediaRequest.SetRoundTripData(a)
	return []*amass.GetRequest{avibaseRequest, wikipediaRequest}
}

func processGetResponses(resps []*amass.GetResponse) []bird.BirdName {
	m := make(map[string]*bird.BirdName)
	for _, resp := range resps {
		a := &avibaseEntry{}
		resp.GetRoundTripData(a)
		bn, ok := m[a.LatinName]
		if !ok {
			bn = bird.Name(a.EnglishName, a.LatinName)
			m[bn.Latin] = bn
		}
		if resp.Site == avibaseSiteKey {
			doc := resp.AsDocument()
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				synonym := s.Find("td").Eq(1).Text()
				synType := s.Find("td").Eq(2).Text()
				if shouldRecordSynonymType(synType) {
					bn.AddLatinSynonym(synonym)
				}
			})
		}
	}
	for _, resp := range resps {
		if resp.Site == wikipediaSiteKey {
			a := &avibaseEntry{}
			resp.GetRoundTripData(a)
			bn, _ := m[a.LatinName]
			w := &wikipediaApi{}
			resp.AsXMLObject(w)
			for _, r := range w.Query.Pages.Page.Redirects.Rd {
				t := r.Title
				if shouldRecordEnglishSynonym(bn, t) {
					bn.AddInformalSynonym(t)
				}
			}
		}
	}
	i := 0
	result := make([]bird.BirdName, len(m))
	for _, bn := range m {
		result[i] = *bn
		i++
	}
	return result
}

func shouldRecordSynonymType(t string) bool {
	return strings.Contains(t, "currently in use") || strings.Contains(t, "protonym")
}

func shouldRecordEnglishSynonym(n *bird.BirdName, e string) bool {
	if strings.Contains(e, ":") {
		return false
	}
	return true
}

// From my favorite: https://www.onlinetool.io/xmltogo/
type wikipediaApi struct {
	XMLName       xml.Name `xml:"api"`
	Text          string   `xml:",chardata"`
	Batchcomplete string   `xml:"batchcomplete,attr"`
	Query         struct {
		Text  string `xml:",chardata"`
		Pages struct {
			Text string `xml:",chardata"`
			Page struct {
				Text                 string `xml:",chardata"`
				Idx                  string `xml:"_idx,attr"`
				Pageid               string `xml:"pageid,attr"`
				Ns                   string `xml:"ns,attr"`
				Title                string `xml:"title,attr"`
				Contentmodel         string `xml:"contentmodel,attr"`
				Pagelanguage         string `xml:"pagelanguage,attr"`
				Pagelanguagehtmlcode string `xml:"pagelanguagehtmlcode,attr"`
				Pagelanguagedir      string `xml:"pagelanguagedir,attr"`
				Touched              string `xml:"touched,attr"`
				Lastrevid            string `xml:"lastrevid,attr"`
				Length               string `xml:"length,attr"`
				Fullurl              string `xml:"fullurl,attr"`
				Editurl              string `xml:"editurl,attr"`
				Canonicalurl         string `xml:"canonicalurl,attr"`
				Redirects            struct {
					Text string `xml:",chardata"`
					Rd   []struct {
						Text   string `xml:",chardata"`
						Pageid string `xml:"pageid,attr"`
						Ns     string `xml:"ns,attr"`
						Title  string `xml:"title,attr"`
					} `xml:"rd"`
				} `xml:"redirects"`
			} `xml:"page"`
		} `xml:"pages"`
	} `xml:"query"`
}
