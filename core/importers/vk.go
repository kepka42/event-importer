package importers

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"event-importer/core"
)

type VK struct {
	url   string
	token string
}

type RootResponse struct {
	Response VKResponse `json:"response"`
}

type VKResponse struct {
	Count uint64 `json:"count"`
	Items []Item `json:"items"`
}

type Item struct {
	Id int64 `json:"id"`
	Album_id int64 `json:"album_id"`
	Owner_id int64 `json:"owner_id"`
	Text string `json:"text"`
	Date int64 `json:"date"`
	Lat float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (v *VK) Init(token string) error {
	v.url = "https://api.vk.com/method/photos.search"
	v.token = token
	return nil
}

func (v *VK) Upload(lat float64, long float64, radius int) ([]core.Pin, error) {
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

	return v.mapToPin(p.Response.Items), nil
}

func (v *VK) mapToPin(items []Item) []core.Pin {
	pins := make([]core.Pin, 0)

	for _, item := range items {
		pin := core.Pin{
			Lat: item.Lat,
			Long: item.Long,
		}

		pins = append(pins, pin)
	}

	return pins
}
