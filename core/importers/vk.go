package importers

import (
	"encoding/json"
	"event-importer/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type VK struct {
	url   string
	token string
}

type RootPhotosResponse struct {
	Response VKPhotosResponse `json:"response"`
}

type VKPhotosResponse struct {
	Count int    `json:"count"`
	Items []Item `json:"items"`
}

type Item struct {
	ID      int     `json:"id"`
	Text    string  `json:"text"`
	Date    int64   `json:"date"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
	Sizes   []Size  `json:"sizes"`
	OwnerID int     `json:"owner_id"`
}

type Size struct {
	URL string `json:"url"`
}

type RootUsersResponse struct {
	Response []User `json:"response"`
}

type User struct {
	ID    int    `json:"id"`
	Sex   int    `json:"sex"`
	Bdate string `json:"bdate"`
}

func (v *VK) Init(token string) error {
	v.url = "https://api.vk.com/method/"
	v.token = token
	return nil
}

func (v *VK) Download(lat float64, long float64, radius int) ([]models.Point, error) {
	points := make([]models.Point, 0)
	offset := 0

	client := &http.Client{}
	for {
		photos, err := v.getPhotos(lat, long, radius, offset, client)

		if err != nil {
			return nil, err
		}

		if len(photos) == 0 {
			break
		}

		userIds := make([]int, 0, len(photos))
		for k, _ := range photos {
			userIds = append(userIds, k)
		}

		users, err := v.getUsers(userIds, client)

		offset += 1000
		points = append(points, v.mapToPoint(photos, users)...)
	}

	return points, nil
}

func (v *VK) Type() string {
	return "vk"
}

func (v *VK) getPhotos(lat float64, long float64, radius int, offset int, client *http.Client) (map[int][]Item, error) {
	req, err := http.NewRequest("GET", v.url+"photos.search", nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("access_token", v.token)
	query.Add("lat", fmt.Sprintf("%f", lat))
	query.Add("long", fmt.Sprintf("%f", long))
	query.Add("radius", strconv.Itoa(radius))
	query.Add("v", "5.102")
	query.Add("count", "1000")
	query.Add("offset", strconv.Itoa(offset))
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var p RootPhotosResponse
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}

	photos := make(map[int][]Item)
	for _, val := range p.Response.Items {
		if obj, ok := photos[val.OwnerID]; ok {
			photos[val.OwnerID] = append(obj, val)
		} else {
			photos[val.OwnerID] = make([]Item, 0)
			photos[val.OwnerID] = append(obj, val)
		}
	}

	return photos, nil
}

func (v *VK) getUsers(ids []int, client *http.Client) (map[int]User, error) {
	req, err := http.NewRequest("GET", v.url+"users.get", nil)

	if err != nil {
		return nil, err
	}

	strIds := ""
	for i, id := range ids {
		if i == 0 {
			strIds += strconv.Itoa(id)
		} else {
			strIds += "," + strconv.Itoa(id)
		}
	}

	query := req.URL.Query()
	query.Add("access_token", v.token)
	query.Add("user_ids", strIds)
	query.Add("fields", "bdate,sex,city")
	query.Add("v", "5.102")
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var p RootUsersResponse
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}

	users := make(map[int]User)
	for _, val := range p.Response {
		users[val.ID] = val
	}

	return users, nil
}

func (v *VK) mapToPoint(items map[int][]Item, users map[int]User) []models.Point {
	pins := make([]models.Point, 0)

	for k, item := range items {
		gender := new(string)
		age := new(int)
		if user, ok := users[k]; ok {
			if user.Sex == 1 {
				t := "female"
				gender = &t
			} else if user.Sex == 2 {
				t := "male"
				gender = &t
			}

			if user.Bdate != "" {
				t, err := dateToUnix(user.Bdate)
				if err == nil {
					current := time.Now().Unix()
					diff := current - t
					tm := time.Unix(diff, 0)

					temp := tm.Year()
					age = &temp
				}
			}
		}

		for _, val := range item {
			pin := models.Point{
				ID:         val.ID,
				Text:       val.Text,
				Lat:        val.Lat,
				Long:       val.Long,
				SocialType: v.Type(),
				Gender:     gender,
				Age:        age,
				URL:        val.Sizes[len(val.Sizes)-1].URL,
				UserID:     val.OwnerID,
			}

			pins = append(pins, pin)
		}
	}

	return pins
}

func dateToUnix(date string) (int64, error) {
	layout := "21.9.1986"
	t, err := time.Parse(layout, date)

	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
