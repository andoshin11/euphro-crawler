package uploader

import (
	"context"
	"log"

	"github.com/andoshin11/artscape-crawler/src/repository"
	"github.com/andoshin11/artscape-crawler/src/types"
)

const collection = "artscape"

// ArtscapeUploader interface
type ArtscapeUploader interface {
	RegisterArtscapeMuseum(ctx context.Context, id string, ch *types.Channels)
}

type artscapeUploader struct {
	MuseumRepository repository.MuseumRepository
}

// NewArtscapeUploader returns struct
func NewArtscapeUploader(MuseumRepository repository.MuseumRepository) ArtscapeUploader {
	return &artscapeUploader{MuseumRepository}
}

func (u *artscapeUploader) RegisterArtscapeMuseum(ctx context.Context, id string, ch *types.Channels) {
	_, err := u.MuseumRepository.GetByID(ctx, collection, id)
	if err != nil {
		log.Println("New Item: " + id)
		err = u.MuseumRepository.AddMuseum(ctx, collection, id)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer func() { ch.UploaderDone <- 0 }()
}
