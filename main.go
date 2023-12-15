package main

import (
	"Awesome/component"
	"time"
)

func main() {
	supervisers := component.Start(1)
	for i := range supervisers {
		go supervisers[i].PrintInfo()
	}
	for {
		time.Sleep(time.Second * 1)
	}
}
