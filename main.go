package main

import (
	"Awesome/component"
	"log"
)

func main() {
	err := component.DeviceFiner()
	if err != nil {
		log.Fatal(err)
		return
	}
}
