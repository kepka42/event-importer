package core

import "event-importer/models"

type Importer interface {
	Init(token string) error
	Upload(lat float64, long float64, radius int) ([]models.Point, error)
	Type() string
}
