package database

import (
	"event-importer/models"
)

func (d *Database) GetLocations() ([]models.Location, error) {
	rows, err := d.db.Query("select city_id, latitude, longitude from locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]models.Location, 0)
	for rows.Next() {
		loc := models.Location{}
		err := rows.Scan(&loc.CityID, &loc.Long, &loc.Lat)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}