package core

import (
	"event-importer/core/database"
	"event-importer/logger"
	"event-importer/models"
	"strconv"
	"time"
)

type Query struct {
	LocationID int
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
	if m.query.LocationID != 0 {
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

	if len(locations) == 0 {
		logger.LogError("empty locations")
		return nil
	}

	for _, location := range locations {
		pins := make([]models.Point, 0)
		for _, imp := range m.importers {
			res, err := imp.Download(location)
			if err != nil {
				return err
			}

			pins = append(pins, res...)
		}

		count_pins := len(pins)
		logger.Log("downloaded " + strconv.Itoa(count_pins) + " points")

		err := m.db.SavePoints(location, pins)
		if err != nil {
			return err
		}

		logger.Log("saved points " + strconv.Itoa(len(pins)))

		max_date := time.Now()
		if len(pins) != 0 {
			max_date = maxDate(pins)
		}

		logger.Log("max date is " + max_date.Format(time.RFC3339) + " for location " + strconv.Itoa(location.ID))

		ids := make([]int, 0)
		ids = append(ids, location.ID)
		err = m.db.UpdateStartFrom(ids, &max_date)
		if err != nil {
			return err
		}

		logger.Log("updated location")
	}

	return nil
}

func maxDate(pins []models.Point) time.Time {
	max := int64(0)

	for _, p := range pins {
		current := p.Date
		if current > max {
			max = current
		}
	}

	return time.Unix(max, 0)
}
