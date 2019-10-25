package database

import (
	"event-importer/models"
	"fmt"
)

func (d *Database) SavePoints(location *models.Location, points []models.Point) error {
	if len(points) == 0 {
		return nil
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	ids := make([]int64, 0)

	for _, v := range points {
		result, err := tx.Exec(`INSERT INTO points(social_id, social_type, description, photo, gender, age, has_children, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE social_id = social_id`, v.ID, v.SocialType, v.Text, v.URL, v.Gender, v.Age, v.HasChildren, v.Lat, v.Long)
		if err != nil {
			continue
		}

		id, err := result.LastInsertId()
		if err != nil {
			continue
		}
		ids = append(ids, id)

		fmt.Println("Added point:", v.ID, v.SocialType)
	}

	for _, v := range ids {
		tx.Exec(`INSERT INTO location_point(location_id, point_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE location_id = location_id`, location.ID, v)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
