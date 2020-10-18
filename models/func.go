package models

import "strconv"

func MakePointDB(lat float64, long float64) PointDB {
	return PointDB{
		Lat: lat,
		Long: long,
	}
}

func (p PointDB) ToMySQLString() string {
	strLat := strconv.FormatFloat(p.Lat, 'f', 6, 64)
	strLong := strconv.FormatFloat(p.Long, 'f', 6, 64)
	return "Point(" + strLat + " " + strLong + ")"
}