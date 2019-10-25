package database

import (
	"database/sql"
	"event-importer/models"
)

func (d *Database) GetLocationById(ID int) (*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius from locations where id = ?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	loc := new(models.Location)
	err = rows.Scan(&loc.ID, &loc.CityID, &loc.Lat, &loc.Long, &loc.Radius)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

func (d *Database) GetLocationsByCityID(cityID int) ([]*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius from locations where city_id = ?", cityID)
	if err != nil {
		return nil, err
	}

	return d.formatLocations(rows)
}

func (d *Database) GetLocations() ([]*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius from locations")
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
		err := rows.Scan(&loc.ID, &loc.CityID, &loc.Lat, &loc.Long, &loc.Radius)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}
