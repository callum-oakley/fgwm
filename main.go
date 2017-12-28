package main

import (
	"log"
	"os"
	"time"

	"github.com/hot-leaf-juice/fgwm/client"
	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/server"
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func main() {
	if len(os.Args) >= 2 {
		client.Run(os.Args)
		return
	}
	if err := server.Run(os.Args[0], &grid.Options{
		Border:           5,
		MinMargin:        wmutils.Size{10, 10},
		Pad:              wmutils.Size{10, 10},
		Size:             grid.Size{24, 24},
		InitialView:      1,
		FocusTimeout:     time.Second,
		FocussedColour:   0xd8dee9,
		UnfocussedColour: 0x65737e,
	}); err != nil {
		log.Fatal(err)
	}
}
