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

func Run(name string, options *grid.Options) error {
	g, err := grid.New(options)
	if err != nil {
		return err
	}
	go func() {
		log.Fatal(g.WatchWindowEvents())
	}()
	s := &Server{name, g}
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
