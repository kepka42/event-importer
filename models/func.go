package models

import "strconv"

func MakePointDB(lat float64, lng float64) PointDB {
	return PointDB{
		Lat: lat,
		Lng: lng,
	}
}

func (p PointDB) ToMySQLString() string {
	strLat := strconv.FormatFloat(p.Lat, 'f', 6, 64)
	strLng := strconv.FormatFloat(p.Lng, 'f', 6, 64)
	return "Point(" + strLat + ", " + strLng + ")"
}