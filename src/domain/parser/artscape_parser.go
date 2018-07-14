package parser

import (
	"github.com/PuerkitoBio/goquery"
)

// ArtscapeAreaItemsParser return the list of item urls
func ArtscapeAreaItemsParser(doc *goquery.Document) (urls []string) {
	urls = make([]string, 0)
	doc.Find(".exhiInfo").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Find(".headH3D > a").Attr("href")
		if exists {
			urls = append(urls, href)
		}
	})
	return
}
