package server

import (
	"github.com/hot-leaf-juice/fgwm/grid"
	"log"
)

func (s *Server) Snap(struct{}, *struct{}) error {
	err := s.grid.Snap()
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Center(struct{}, *struct{}) error {
	err := s.grid.Center()
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Fullscreen(struct{}, *struct{}) error {
	err := s.grid.Fullscreen()
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Kill(struct{}, *struct{}) error {
	err := s.grid.Kill()
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Move(diff grid.Size, _ *struct{}) error {
	err := s.grid.Move(diff)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Grow(diff grid.Size, _ *struct{}) error {
	err := s.grid.Grow(diff)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Throw(direction grid.Direction, _ *struct{}) error {
	err := s.grid.Throw(direction)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Spread(direction grid.Direction, _ *struct{}) error {
	err := s.grid.Spread(direction)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Focus(strategy grid.FocusStrategy, _ *struct{}) error {
	err := s.grid.Focus(strategy)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) Teleport(rectangle grid.Rectangle, _ *struct{}) error {
	err := s.grid.Teleport(rectangle)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) ViewInclude(n int, _ *struct{}) error {
	err := s.grid.ViewInclude(n)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}

func (s *Server) ViewSet(n int, _ *struct{}) error {
	err := s.grid.ViewSet(n)
	if err != nil {
		log.Printf("%v: %v\n", s.name, err)
	}
	return err
}
