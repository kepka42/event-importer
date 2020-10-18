package importers

import (
	"encoding/json"
	"event-importer/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
	ID        int        `json:"id"`
	Sex       int        `json:"sex"`
	Bdate     string     `json:"bdate"`
	City      City       `json:"city"`
	Relatives []Relative `json:"relatives"`
}

type Relative struct {
	Type string `json:"type"`
}

type City struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (v *VK) Init(token string) error {
	v.url = "https://api.vk.com/method/"
	v.token = token
	return nil
}

func (v *VK) Download(location *models.Location) ([]models.Point, error) {
	points := make([]models.Point, 0)
	offset := 0

	var startFrom *time.Time
	if location.StartFrom.Valid {
		start := time.Unix(location.StartFrom.Int64, 0)
		startFrom = &start
	} else {
		startFrom = nil
	}

	client := &http.Client{}
	for {
		photos, err := v.getPhotos(location.Coordinates.Lat, location.Coordinates.Lng, location.Radius, offset, startFrom, client)

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

func (v *VK) getPhotos(lat float64, long float64, radius int, offset int, startFrom *time.Time, client *http.Client) (map[int][]Item, error) {
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
	if startFrom != nil {
		query.Add("start_time", strconv.FormatInt(startFrom.Unix(), 10))
	}
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
	query.Add("fields", "bdate,sex,city,relatives")
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
	points := make([]models.Point, 0)

	for k, item := range items {
		gender := new(string)
		age := new(int)
		hasChilds := false
		isTourist := new(bool)
		city := new(string)
		if user, ok := users[k]; ok {
			if user.Sex == 1 {
				t := "female"
				gender = &t
			} else if user.Sex == 2 {
				t := "male"
				gender = &t
			}

			if user.Bdate != "" {
				temp := makeAge(user.Bdate)
				age = &temp
			}

			if len(user.Relatives) > 0 {
				for _, relative := range user.Relatives {
					if relative.Type == "child" {
						hasChilds = true
					}
				}
			}

			if user.City.ID > 0 {
				*city = user.City.Title
			}

			isTourist = nil
		}

		if *age == 0 {
			age = nil
		}

		if *gender == "" {
			gender = nil
		}

		for _, val := range item {
			point := models.Point{
				ID:          val.ID,
				Text:        val.Text,
				Coordinates: models.MakePointDB(val.Lat, val.Long),
				Gender:      gender,
				Age:         age,
				URL:         val.Sizes[len(val.Sizes)-1].URL,
				VkUserID:    val.OwnerID,
				HasChildren: hasChilds,
				IsTourist:   isTourist,
				UserCity:	 city,
			}

			points = append(points, point)
		}
	}

	return points
}

func makeAge(date string) int {
	re := regexp.MustCompile(`(?m)^([0-9]{1,2})\.([0-9]{1,2})\.([0-9]{4})$`)
	strs := re.FindStringSubmatch(date)

	if len(strs) != 4 {
		return 0
	}

	days, _ := strconv.Atoi(strs[1])
	months, _ := strconv.Atoi(strs[2])
	years, _ := strconv.Atoi(strs[3])

	current := time.Now()

	age := current.Year() - years
	if (int(current.Month()) < months) || (int(current.Month()) == months && current.Day() < days) {
		age -= 1
	}

	return age
}
