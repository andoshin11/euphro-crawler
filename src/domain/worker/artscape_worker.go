package worker

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/andoshin11/artscape-crawler/src/domain/fetcher"
	"github.com/andoshin11/artscape-crawler/src/domain/uploader"
	"github.com/andoshin11/artscape-crawler/src/repository"
	"github.com/andoshin11/artscape-crawler/src/types"
)

const collection = "collection"

// ArtscapeWorker interdace
type ArtscapeWorker interface {
	Crawl(ctx context.Context)
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
			id := link[23 : len(link)-10]
			fmt.Printf("Success %#v\n", id)

			// Check if the data exists on database
			uploadWc++
			go uploader.RegisterArtscapeMuseum(ctx, id, chs)
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
