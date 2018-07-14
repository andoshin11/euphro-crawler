package handler

import (
	"context"
	"log"

	"github.com/andoshin11/artscape-crawler/src/client"
	"github.com/andoshin11/artscape-crawler/src/domain/worker"
)

// CrawlMuseumListHandler is the request handler
func CrawlMuseumListHandler(ctx context.Context) {
	client, err := client.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	worker := worker.NewArtscapeWorker(client)
	worker.Crawl(ctx)
	return
}
