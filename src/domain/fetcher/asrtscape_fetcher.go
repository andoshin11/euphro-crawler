package fetcher

import (
	"context"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/andoshin11/artscape-crawler/src/domain/parser"
	"github.com/andoshin11/artscape-crawler/src/types"
	"googlemaps.github.io/maps"
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

// ArtscapeMuseumDetailFetcher function
func ArtscapeMuseumDetailFetcher(ctx context.Context, ID string, client *maps.Client, ch *types.DetailChannels) (err error) {
	defer func() { ch.FetcherDone <- 0 }()
	path := "http://artscape.jp/mdb/" + ID + "_1900.html"
	resp, err := http.Get(path)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	museum := parser.ArtscapeMuseumDetailParser(doc)

	result := &maps.GeocodingRequest{
		Address: museum.Address,
	}

	res, err := client.Geocode(ctx, result)

	if len(res) >= 1 {
		Lat := res[0].Geometry.Location.Lat
		Lng := res[0].Geometry.Location.Lng

		museum.Lat = Lat
		museum.Lng = Lng
	}

	ch.FetcherResult <- types.DetailFetcherResult{
		ID:   ID,
		Item: museum,
	}

	return
}
