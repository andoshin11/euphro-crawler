package fetcher

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/andoshin11/artscape-crawler/src/domain/parser"
	"github.com/andoshin11/artscape-crawler/src/types"
)

// ArtscapeItemsFetcher returns the list of child page urls
func ArtscapeItemsFetcher(u string, ch *types.Channels) (err error) {
	defer func() { ch.FetcherDone <- 0 }()
	baseURL, err := url.Parse(u)
	if err != nil {
		return
	}

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	urls := parser.ArtscapeAreaItemsParser(doc)

	for _, item := range urls {
		itemURL, err := baseURL.Parse(item)
		if err == nil {
			ch.FetcherResult <- types.FetcherResult{
				URL: itemURL.String(),
			}
		}
	}
	return
}
