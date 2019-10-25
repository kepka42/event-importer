package core

type Importer interface {
	Init(token string) error
	Upload(lat float64, long float64, radius int) ([]Pin, error)
}

type Pin struct {
	URL string
	Gender string
	Age int
	HasChildren bool
	Lat float64
	Long float64
}

type Params struct {
	SocialID string
	CityID int
	PointID int
	SocialType string
}