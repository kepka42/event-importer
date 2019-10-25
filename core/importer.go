package core

type Importer interface {
	Init(token string) error
	Upload(lat float64, long float64, radius int) (interface{}, error)
}

type Pin struct {

}