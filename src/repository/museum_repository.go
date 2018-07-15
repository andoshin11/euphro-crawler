package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/andoshin11/artscape-crawler/src/types"
	"google.golang.org/api/iterator"
)

// MuseumRepository interface
type MuseumRepository interface {
	GetAll(ctx context.Context, collection string) ([]*types.Museum, error)
	GetByID(ctx context.Context, collection, id string) (*types.Museum, error)
	AddMuseum(ctx context.Context, collection, id string) error
	UpdateMuseum(ctx context.Context, collection, identifier string, museum *types.Museum) (err error)
}

type museumRepository struct {
	Client *firestore.Client
}

// NewMuseumRepository return struct
func NewMuseumRepository(Client *firestore.Client) MuseumRepository {
	return &museumRepository{Client}
}

func (r *museumRepository) GetAll(ctx context.Context, collection string) ([]*types.Museum, error) {
	iter := r.Client.Collection(collection).Documents(ctx)
	fmt.Println("get all")

	museums := make([]*types.Museum, 10)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		museum := types.Museum{}
		doc.DataTo(&museum)
		museums = append(museums, &museum)
	}

	return museums, nil
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

func (r *museumRepository) AddMuseum(ctx context.Context, collection, identifier string) (err error) {
	t := time.Now()

	_, err = r.Client.Collection(collection).Doc(identifier).Set(ctx, map[string]interface{}{
		"CreatedAt":  t,
		"Identifier": identifier,
	})

	return
}

func (r *museumRepository) UpdateMuseum(ctx context.Context, collection, identifier string, museum *types.Museum) (err error) {
	t := time.Now()

	_, err = r.Client.Collection(collection).Doc(identifier).Set(ctx, map[string]interface{}{
		"Name":      museum.Name,
		"Address":   museum.Address,
		"Img":       museum.Img,
		"Entry":     museum.Entry,
		"SiteURL":   museum.SiteURL,
		"UpdatedAt": t,
		"Lat":       museum.Lat,
		"Lng":       museum.Lng,
	}, firestore.MergeAll)

	return
}
