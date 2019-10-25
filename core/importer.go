package core

import (
	"event-importer/models"
)

type Importer interface {
	Init(token string) error
	Download(location *models.Location) ([]models.Point, error)
	Type() string
}
