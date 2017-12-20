package server

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/hot-leaf-juice/fgwm/grid"
)

type Server struct {
	name string
	grid *grid.Grid
}

type Config struct {
	Name string
	Grid *grid.Options
}

func Run(config *Config) error {
	g, err := grid.New(config.Grid)
	if err != nil {
		return err
	}
	go func() {
		log.Fatal(g.WatchWindowEvents())
	}()
	s := &Server{config.Name, g}
	rpc.Register(s)
	rpc.HandleHTTP()
	// TODO replace with named pipes(?) or unix sockets(?) or something
	listener, err := net.Listen("tcp", ":62676")
	if err != nil {
		return err
	}
	log.Printf("%v: listening on localhost:62676\n", s.name)
	http.Serve(listener, nil)
	return nil
}
