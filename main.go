package main

import (
	"context"

	"github.com/andoshin11/artscape-crawler/src/handler"
)

func main() {
	ctx := context.Background()
	handler.CrawlMuseumListHandler(ctx)
}
