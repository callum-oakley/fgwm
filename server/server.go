package server

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/hot-leaf-juice/fgwm/grid"
)

type Server struct {
	grid *grid.Grid
}

type Config struct {
	Grid *grid.Options
}

func Run(config *Config) error {
	g, err := grid.New(config.Grid)
	if err != nil {
		return err
	}
	s := &Server{g}
	rpc.Register(s)
	rpc.HandleHTTP()
	// TODO replace with named pipes
	listener, err := net.Listen("tcp", ":62676")
	if err != nil {
		return err
	}
	fmt.Println("Listening on port localhost:62676")
	http.Serve(listener, nil)
	return nil
}

// Just testing

type Args struct {
	A, B int
}

func (s *Server) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
