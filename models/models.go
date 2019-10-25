package models

type Point struct {
	ID int
	URL string
	Gender string
	Age int
	HasChildren bool
	Lat float64
	Long float64
	Text string
	SocialType string
}

type Location struct {
	CityID int
	Lat float64
	Long float64
	Radius int
}