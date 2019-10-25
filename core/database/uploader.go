package database

import (
	"fmt"
	"event-importer/models"
)

func (d *Database) SavePoints(location models.Location, points []models.Point) error {
	sqlStr := "REPLACE INTO points(social_id, social_type, description, photo, gender, age, has_children, latitude, longitude) VALUES"
	for k, v := range points {
		if k == 0 {
			sqlStr += fmt.Sprintf("('%v', '%v', '%v', '%v', '%v', %v, %v, %v, %v)", v.ID, v.SocialType, v.Text, v.URL, v.Gender, v.Age, v.HasChildren, v.Lat, v.Long)
		} else {
			sqlStr += fmt.Sprintf(", ('%v', '%v', '%v', '%v', '%v', %v, %v, %v, %v)", v.ID, v.SocialType, v.Text, v.URL, v.Gender, v.Age, v.HasChildren, v.Lat, v.Long)
		}
	}

	_, err := d.db.Exec(sqlStr)

	fmt.Println(sqlStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}