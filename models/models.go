package models

import (
	"database/sql"
)

type Point struct {
	ID          int
	URL         string
	Gender      *string
	Age         *int
	HasChildren bool
	IsTourist   *bool
	Coordinates PointDB
	Text        string
	VkUserID    int
	UserCity  	*string
	UserCityId  int
	Date 		int64
}

type Location struct {
	ID          int
	Coordinates PointDB
	Radius      int
	StartFrom   sql.NullInt64
}

type PointDB struct {
	Lat float64
	Lng float64
}

