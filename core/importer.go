package core

import "event-importer/models"

type Importer interface {
	Init(token string) error
	Download(lat float64, long float64, radius int) ([]models.Point, error)
	Type() string
}
