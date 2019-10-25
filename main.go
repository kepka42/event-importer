package main

import (
	"event-importer/core"
	"event-importer/core/importers"
	"fmt"
)

func main() {
	imps := initImporters()

	manager := &core.Manager{}
	err := manager.Init(imps, "")

	if err != nil {
		fmt.Println(err)
	}

	err = manager.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func initImporters() []core.Importer {
	imps := make([]core.Importer, 0)

	vk := &importers.VK{}
	vk.Init("")

	imps = append(imps, vk)

	return imps
}


