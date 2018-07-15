package types

import "time"

// Museum type
type Museum struct {
	Identifier string    `firestore:"Identifier"`
	CreatedAt  time.Time `firestore:"CreatedAt"`
	UpdatedAt  time.Time
	Name       string
	Address    string
	Img        string
	Entry      string
	SiteURL    string
	Lat        float64
	Lng        float64
}
