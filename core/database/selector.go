package database

import (
	"database/sql"
	"event-importer/models"
	"time"
)

func (d *Database) GetLocationById(ID int) (*models.Location, error) {
	rows, err := d.db.Query("select id, ST_X(coordinates) latitude, ST_X(coordinates) longitude, radius, start_from from locations where id = ?", ID)
	if err != nil {
		return nil, err
	}

	locations, err := d.formatLocations(rows)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, nil
	}

	return locations[0], nil
}

func (d *Database) GetLocations() ([]*models.Location, error) {
	rows, err := d.db.Query("select id, ST_X(coordinates) latitude, ST_Y(coordinates) longitude, radius, start_from from locations")
	if err != nil {
		return nil, err
	}

	return d.formatLocations(rows)
}

func (d *Database) formatLocations(rows *sql.Rows) ([]*models.Location, error) {
	defer rows.Close()
	locations := make([]*models.Location, 0)
	for rows.Next() {
		loc := new(models.Location)

		var t *time.Time
		coordinates := models.PointDB{}
		err := rows.Scan(&loc.ID, &coordinates.Lat, &coordinates.Long, &loc.Radius, &t)

		loc.Coordinates = coordinates
		if err != nil {
			return nil, err
		}

		if t == nil {
			loc.StartFrom.Int64 = 0
			loc.StartFrom.Valid = false
		} else {
			loc.StartFrom.Int64 = t.Unix()
			loc.StartFrom.Valid = true
		}
		locations = append(locations, loc)
	}

	return locations, nil
}
