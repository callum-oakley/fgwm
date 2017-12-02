package main

import (
	"log"
	"os"

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
	if err := server.Run(&server.Config{
		Name: os.Args[0],
		Grid: &grid.Options{
			Border:    wmutils.Size{5, 5},
			MinMargin: wmutils.Size{10, 10},
			Pad:       wmutils.Size{10, 10},
			Size:      grid.Size{24, 24},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
