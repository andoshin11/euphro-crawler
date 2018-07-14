package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/andoshin11/artscape-crawler/src/types"
)

// MuseumRepository interface
type MuseumRepository interface {
	GetByID(ctx context.Context, collection, id string) (*types.Museum, error)
	AddMuseum(ctx context.Context, collection, id string) error
}

type museumRepository struct {
	Client *firestore.Client
}

// NewMuseumRepository return struct
func NewMuseumRepository(Client *firestore.Client) MuseumRepository {
	return &museumRepository{Client}
}

func (r *museumRepository) GetByID(ctx context.Context, collection, id string) (*types.Museum, error) {
	doc, err := r.Client.Collection(collection).Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	museum := types.Museum{}
	doc.DataTo(&museum)

	return &museum, err
}

func (r *museumRepository) AddMuseum(ctx context.Context, collection, id string) (err error) {
	t := time.Now()

	_, err = r.Client.Collection(collection).Doc(id).Set(ctx, map[string]interface{}{
		"createdAt": t,
		"id":        id,
	})

	return
}
