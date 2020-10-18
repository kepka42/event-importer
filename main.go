package main

import (
	"errors"
	"event-importer/core"
	"event-importer/core/importers"
	"flag"
	"fmt"
)

type Params struct {
	VKtoken      string
	DBconnection string
	LocationID   int
}

func main() {
	params := parseParams()

	check_error := checkParams(params)
	if check_error != nil {
		fmt.Println(check_error.Error())
		return
	}

	imps := initImporters(&params)

	manager := &core.Manager{}
	err := manager.Init(imps, params.DBconnection, core.Query{
		LocationID: params.LocationID,
	})

	if err != nil {
		fmt.Println(err)
	}

	err = manager.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func initImporters(params *Params) []core.Importer {
	imps := make([]core.Importer, 0)

	vk := &importers.VK{}
	err := vk.Init(params.VKtoken)
	if err != nil {
		panic(err)
	}

	imps = append(imps, vk)

	return imps
}

func parseParams() Params {
	vk := flag.String("vk", "", "token for vk")
	db := flag.String("db", "", "connection for db")
	loc := flag.Int("location", 0, "location id")

	flag.Parse()

	params := Params{}
	params.VKtoken = *vk
	params.DBconnection = *db
	params.LocationID = *loc

	return params
}

func checkParams(params Params) error {
	if len(params.DBconnection) == 0 {
		return errors.New("DB connection can not be empty")
	}

	if len(params.VKtoken) == 0 {
		return errors.New("VK token can not by empty")
	}

	return nil
}
