package parser

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/andoshin11/artscape-crawler/src/types"
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

// ArtscapeMuseumDetailParser return the museum struct
func ArtscapeMuseumDetailParser(doc *goquery.Document) (museum *types.Museum) {
	Name := doc.Find(".mainColHeading > h2").Text()
	Address := doc.Find(".address").Text()
	Img, exists := doc.Find(".imageArea > p > img").Attr("src")
	if exists {
		Img = "http://artscape.jp" + Img
	}
	Entry := doc.Find(".entryArea > p").Text()
	SiteURL, _ := doc.Find(".boxLinkC > li").First().Find("a").Attr("href")

	fmt.Println(Name)
	fmt.Println(Address)
	fmt.Println(Img)
	fmt.Println(Entry)
	fmt.Println(SiteURL)

	museum = &types.Museum{
		Name:    Name,
		Address: Address,
		Img:     Img,
		Entry:   Entry,
		SiteURL: SiteURL,
	}

	return
}
