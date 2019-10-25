package importers

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"event-importer/models"
)

type VK struct {
	url   string
	token string
}

type RootResponse struct {
	Response VKResponse `json:"response"`
}

type VKResponse struct {
	Count int `json:"count"`
	Items []Item `json:"items"`
}

type Item struct {
	ID int `json:"id"`
	Text string `json:"text"`
	Date int64 `json:"date"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
	Sizes []Size `json:"sizes"`
}

type Size struct {
	URL string `json:"url"`
}

func (v *VK) Init(token string) error {
	v.url = "https://api.vk.com/method/photos.search"
	v.token = token
	return nil
}

func (v *VK) Upload(lat float64, long float64, radius int) ([]models.Point, error) {
	items := make([]Item, 0)
	offset := 0

	for {
		req, err := http.NewRequest("GET", v.url, nil)

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

		client := &http.Client{}

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var p RootResponse
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}

		if len(p.Response.Items) == 0 {
			break
		}

		offset += 1000
		items = append(items, p.Response.Items...)
	}

	return v.mapToPin(items), nil
}

func (v *VK) Type() string {
	return "vk"
}

func (v *VK) mapToPin(items []Item) []models.Point {
	pins := make([]models.Point, 0)

	for _, item := range items {
		pin := models.Point{
			ID: item.ID,
			Text: item.Text,
			Lat: item.Lat,
			Long: item.Long,
			SocialType: v.Type(),
			Gender: "female",
			Age: 13,
			URL: item.Sizes[len(item.Sizes) - 1].URL,
		}

		pins = append(pins, pin)
	}

	return pins
}
