package core

import (
	"event-importer/core/database"
	"event-importer/models"
	"time"
)

type Query struct {
	LocationID int
	CityID     int
}

type Manager struct {
	importers []Importer
	query     Query
	db        *database.Database
}

func (m *Manager) Init(importers []Importer, dbConn string, query Query) error {
	m.importers = importers
	m.db = &database.Database{}
	m.query = query
	err := m.db.Init(dbConn)

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Run() error {
	var err error
	locations := make([]*models.Location, 0)
	if m.query.CityID != 0 {
		locations, err = m.db.GetLocationsByCityID(m.query.CityID)
		if err != nil {
			return err
		}
	} else if m.query.LocationID != 0 {
		location, err := m.db.GetLocationById(m.query.LocationID)
		if err != nil {
			return err
		}

		locations = append(locations, location)
	} else {
		locations, err = m.db.GetLocations()
		if err != nil {
			return err
		}
	}

	ids := make([]int, 0)
	for _, location := range locations {
		pins := make([]models.Point, 0)
		for _, imp := range m.importers {
			res, err := imp.Download(location)

			if err != nil {
				return err
			}

			pins = append(pins, res...)
		}

		err := m.db.SavePoints(location, pins)
		if err != nil {
			return err
		}

		ids = append(ids, location.ID)
	}

	now := time.Now()
	err = m.db.UpdateStartFrom(ids, &now)
	if err != nil {
		return err
	}

	return nil
}
