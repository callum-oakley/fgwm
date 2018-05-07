package main

import (
	"fmt"
	"log"
	"os"

	"github.com/callum-oakley/fgwm/client"
	"github.com/callum-oakley/fgwm/config"
	"github.com/callum-oakley/fgwm/server"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] != "-c" && os.Args[1] != "--config" {
		client.Run(os.Args)
		return
	}
	configPath := fmt.Sprintf(
		"%v/.config/fgwm/config.toml",
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
