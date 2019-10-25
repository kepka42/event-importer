package core

import (
	"event-importer/core/database"
	"event-importer/models"
)

type Manager struct {
	importers []Importer
	db *database.Database
}

func (m *Manager) Init(importers []Importer, dbConn string) error {
	m.importers = importers
	m.db = &database.Database{}
	err := m.db.Init(dbConn)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Run() error {
	locations, err := m.db.GetLocations()

	if err != nil {
		return err
	}

	for _, location := range locations {
		pins := make([]models.Point, 0)
		for _, imp := range m.importers {
			res, err := imp.Upload(location.Lat, location.Long, location.Radius)

			if err != nil {
				return err
			}

			pins = append(pins, res...)
		}

		m.db.SavePoints(location, pins)
	}

	return nil
}