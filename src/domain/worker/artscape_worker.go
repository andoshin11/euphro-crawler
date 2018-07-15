package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/andoshin11/artscape-crawler/src/domain/fetcher"
	"github.com/andoshin11/artscape-crawler/src/domain/uploader"
	"github.com/andoshin11/artscape-crawler/src/repository"
	"github.com/andoshin11/artscape-crawler/src/types"
	"googlemaps.github.io/maps"
)

const collection = "artscape"

// ArtscapeWorker interdace
type ArtscapeWorker interface {
	Crawl(ctx context.Context)
	CrawlDetail(ctx context.Context)
}

type artscapeWorker struct {
	Client *firestore.Client
}

// NewArtscapeWorker return struct
func NewArtscapeWorker(Client *firestore.Client) ArtscapeWorker {
	return &artscapeWorker{Client}
}

func (w *artscapeWorker) Crawl(ctx context.Context) {
	museumRepository := repository.NewMuseumRepository(w.Client)
	uploader := uploader.NewArtscapeUploader(museumRepository)

	// worker count
	wc := 0
	uploadWc := 0

	chs := types.NewChannels()

	// 47都道府県の各エリア
	for i := 1; i <= 47; i++ {
		wc++
		id := strconv.Itoa(i)
		url := "http://artscape.jp/mdb/mdb_result.php?area=" + id
		go fetcher.ArtscapeItemsFetcher(url, chs)
	}

	done := false
	for !done {
		select {
		case res := <-chs.FetcherResult:
			link := res.URL
			Identifier := link[23 : len(link)-10]
			fmt.Printf("Success %#v\n", Identifier)

			// Check if the data exists on database
			uploadWc++
			go uploader.RegisterArtscapeMuseum(ctx, Identifier, chs)
		case <-chs.FetcherDone:
			wc--
			if wc == 0 && uploadWc == 0 {
				done = true
			}
		case <-chs.UploaderDone:
			uploadWc--
			if wc == 0 && uploadWc == 0 {
				done = true
			}
		}
	}
}

func (w *artscapeWorker) CrawlDetail(ctx context.Context) {
	museumRepository := repository.NewMuseumRepository(w.Client)
	uploader := uploader.NewArtscapeUploader(museumRepository)

	// worker count
	fetcherWc := 0
	uploadWc := 0

	chs := types.NewDetailChannels()

	museums, err := museumRepository.GetAll(ctx, collection)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAP_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	for _, museum := range museums {
		if museum != nil {
			fetcherWc++
			log.Println("start fetching")
			go fetcher.ArtscapeMuseumDetailFetcher(ctx, museum.Identifier, c, chs)
		}
	}

	done := false
	for !done {
		select {
		case res := <-chs.FetcherResult:
			fmt.Printf("Success %#v\n", res.Item)

			// Check if the data exists on database
			uploadWc++
			go uploader.UpdateArtscapeMuseum(ctx, res.ID, res.Item, chs)
		case <-chs.FetcherDone:
			fetcherWc--
			if fetcherWc == 0 && uploadWc == 0 {
				done = true
			}
		case <-chs.UploaderDone:
			uploadWc--
			if fetcherWc == 0 && uploadWc == 0 {
				done = true
			}
		}
	}
}
