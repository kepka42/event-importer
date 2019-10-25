package main

import (
	"event-importer/core/importers"
	"fmt"
)

func main() {
	vk := &importers.VK{}
	vk.Init("")
	pins, err := vk.Upload(55.343837, 86.077922, 30000)
	if err != nil {
		panic(err)
	}

	fmt.Println(pins[0].Long)
}


