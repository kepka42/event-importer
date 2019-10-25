package models

import "time"

type Point struct {
	ID          int
	URL         string
	Gender      *string
	Age         *int
	HasChildren bool
	IsTourist   *bool
	Lat         float64
	Long        float64
	Text        string
	SocialType  string
	UserID      int
	CreatedAT   time.Time
	UpdatedAT   time.Time
}

type Location struct {
	ID          int
	CityID      int
	Lat         float64
	Long        float64
	Radius      int
	CitySocials []*CitySocial
}

type CitySocial struct {
	ID         int
	SocialID   int
	SocialType string
}
