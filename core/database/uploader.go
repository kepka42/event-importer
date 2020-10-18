package database

import (
	"event-importer/logger"
	"event-importer/models"
	"strconv"
	"strings"
	"time"
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
		pointStr := v.Coordinates.ToMySQLString()
		result, err := tx.Exec(
			`INSERT INTO points(location_id, photo, gender, age, has_children, coordinates, is_tourist, vk_user_id, user_city) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE location_id = location_id`,
			location.ID, v.URL, v.Gender, v.Age, v.HasChildren, pointStr, v.IsTourist, v.VkUserID, v.UserCity,
			)
		if err != nil {
			logger.LogError("can not pass point: " + err.Error())
			continue
		}

		id, err := result.LastInsertId()
		if err != nil {
			logger.LogError("can not get last insert id: " + err.Error())
			continue
		}
		ids = append(ids, id)

		logger.Log("saved point: " + strconv.Itoa(v.ID))
	}

	return nil
}

func (d *Database) UpdateStartFrom(locIds []int, t *time.Time) error {
	args := make([]interface{}, len(locIds)+1)
	args[0] = t
	for i, id := range locIds {
		args[i+1] = id
	}

	_, err := d.db.Exec(`UPDATE locations SET start_from = ? where id in (?`+strings.Repeat(`,?`, len(args)-2)+`)`, args...)
	if err != nil {
		return nil
	}

	return nil
}
