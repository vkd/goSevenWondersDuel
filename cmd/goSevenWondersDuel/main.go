package main

import (
	"log"

	"github.com/vkd/goSevenWondersDuel/gfx"
)

func main() {
	log.Printf("7 Wonders!")

	err := gfx.Run()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
