package main

import (
	"fmt"
	"log"

	"github.com/hot-leaf-juice/fgwm/wm"
)

func main() {
	wid, err := wm.Focussed()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pfw: %v\n", wid)
	wm.Teleport(
		wid,
		wm.Position{wm.Pixels(800), wm.Pixels(100)},
		wm.Size{wm.Pixels(700), wm.Pixels(500)},
	)
}
