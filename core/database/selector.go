package database

import (
	"event-importer/models"
)

func (d *Database) GetLocations() ([]models.Location, error) {
	rows, err := d.db.Query("select city_id, latitude, longitude, radius from locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]models.Location, 0)
	for rows.Next() {
		loc := models.Location{}
		err := rows.Scan(&loc.CityID, &loc.Lat, &loc.Long, &loc.Radius)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}