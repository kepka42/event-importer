package main

import "event-importer/core/importers"

func main() {
	vk := &importers.VK{}
	vk.Init("")
	vk.Upload(55.343837, 86.077922, 30000)
}


