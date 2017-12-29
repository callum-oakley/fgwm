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
	if len(os.Args) >= 2 && os.Args[1] != "-c" && os.Args[1] != "--config" {
		client.Run(os.Args)
		return
	}
	configPath := fmt.Sprintf(
		"%v/.config/fgwm/fgwm.toml",
		os.Getenv("HOME"),
	)
	if len(os.Args) >= 3 && (os.Args[1] == "-c" || os.Args[1] == "--config") {
		configPath = os.Args[2]
	}
	options, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("%v: %v\n", os.Args[0], err)
	}
	if err := server.Run(os.Args[0], options); err != nil {
		log.Fatalf("%v: %v\n", os.Args[0], err)
	}
}
