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
	UpdateArtscapeMuseum(ctx context.Context, identifier string, museum *types.Museum, ch *types.DetailChannels)
}

type artscapeUploader struct {
	MuseumRepository repository.MuseumRepository
}

// NewArtscapeUploader returns struct
func NewArtscapeUploader(MuseumRepository repository.MuseumRepository) ArtscapeUploader {
	return &artscapeUploader{MuseumRepository}
}

func (u *artscapeUploader) RegisterArtscapeMuseum(ctx context.Context, Identifier string, ch *types.Channels) {
	_, err := u.MuseumRepository.GetByID(ctx, collection, Identifier)
	if err != nil {
		log.Println("New Item: " + Identifier)
		err = u.MuseumRepository.AddMuseum(ctx, collection, Identifier)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer func() { ch.UploaderDone <- 0 }()
}

func (u *artscapeUploader) UpdateArtscapeMuseum(ctx context.Context, identifier string, museum *types.Museum, ch *types.DetailChannels) {
	err := u.MuseumRepository.UpdateMuseum(ctx, collection, identifier, museum)
	if err != nil {
		log.Println(err)
	}

	defer func() { ch.UploaderDone <- 0 }()
}
