package server

import (
	"fmt"
	"io/ioutil"
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
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	err = ioutil.WriteFile("/tmp/fgwm-port", []byte(fmt.Sprint(port)), 0666)
	if err != nil {
		return err
	}
	log.Printf("%v: listening on localhost:%v\n", s.name, port)
	http.Serve(listener, nil)
	return nil
}
