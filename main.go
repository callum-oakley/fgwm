package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hot-leaf-juice/fgwm/client"
	"github.com/hot-leaf-juice/fgwm/config"
	"github.com/hot-leaf-juice/fgwm/server"
)

func main() {
	if len(os.Args) >= 2 {
		client.Run(os.Args)
		return
	}
	options, err := config.Load(fmt.Sprintf(
		"%v/.config/fgwm/fgwm.toml",
		os.Getenv("HOME"),
	))
	if err != nil {
		log.Fatalf("%v: %v\n", os.Args[0], err)
	}
	if err := server.Run(os.Args[0], options); err != nil {
		log.Fatalf("%v: %v\n", os.Args[0], err)
	}
}
