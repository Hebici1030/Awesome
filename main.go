package main

import (
	"Awesome/component"
	"os/exec"
)

func main() {
	supervisers := component.Start(1)
	for i := range supervisers {
		go supervisers[i].PrintInfo()
	}
	command := exec.Command("echo", "Start...")
	command.Start()
}
