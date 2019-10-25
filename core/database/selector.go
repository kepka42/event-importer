package database

import (
	"database/sql"
	"event-importer/models"
	"strings"
	"time"
)

func (d *Database) GetLocationById(ID int) (*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius, start_from from locations where id = ?", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	loc := new(models.Location)

	var t *time.Time
	err = rows.Scan(&loc.ID, &loc.CityID, &loc.Lat, &loc.Long, &loc.Radius, &t)
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

	cities, err := d.getSocialCities([]int{loc.CityID})

	if err != nil {
		return nil, err
	}

	if len(cities) > 0 {
		loc.CitySocials = cities
	}

	return loc, nil
}

func (d *Database) GetLocationsByCityID(cityID int) ([]*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius, start_from from locations where city_id = ?", cityID)
	if err != nil {
		return nil, err
	}

	cities, err := d.getSocialCities([]int{cityID})

	if err != nil {
		return nil, err
	}

	return d.formatLocations(rows, cities)
}

func (d *Database) GetLocations() ([]*models.Location, error) {
	rows, err := d.db.Query("select id, city_id, latitude, longitude, radius, start_from from locations")
	if err != nil {
		return nil, err
	}

	cities, err := d.getSocialCities([]int{})
	if err != nil {
		return nil, err
	}

	return d.formatLocations(rows, cities)
}

func (d *Database) getSocialCities(ids []int) ([]*models.CitySocial, error) {
	var rows *sql.Rows
	var err error

	if len(ids) > 0 {
		args := make([]interface{}, len(ids))
		for i, id := range ids {
			args[i] = id
		}

		rows, err = d.db.Query("select city_id, social_type, social_id from city_socials where city_id in (?"+strings.Repeat(",?", len(args)-1)+")", args...)
	} else {
		rows, err = d.db.Query("select city_id, social_type, social_id from city_socials")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	citySocials := make([]*models.CitySocial, 0)
	for rows.Next() {
		city := new(models.CitySocial)
		err := rows.Scan(&city.ID, &city.SocialType, &city.SocialID)
		if err != nil {
			return nil, err
		}
		citySocials = append(citySocials, city)
	}

	return citySocials, nil
}

func (d *Database) formatLocations(rows *sql.Rows, cities []*models.CitySocial) ([]*models.Location, error) {
	defer rows.Close()
	locations := make([]*models.Location, 0)
	for rows.Next() {
		loc := new(models.Location)

		var t *time.Time
		err := rows.Scan(&loc.ID, &loc.CityID, &loc.Lat, &loc.Long, &loc.Radius, &t)
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

		for _, city := range cities {
			if city.ID == loc.CityID {
				loc.CitySocials = append(loc.CitySocials, city)
			}
		}
	}

	return locations, nil
}
