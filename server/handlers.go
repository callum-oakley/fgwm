package server

import (
	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Rectangle struct {
	TopLeft, BottomRight grid.Position
}

func (s *Server) Snap(struct{}, *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Snap(wid)
}

func (s *Server) Center(struct{}, *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Center(wid)
}

func (s *Server) Kill(struct{}, *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return wmutils.Kill(wid)
}

func (s *Server) Move(diff grid.Size, _ *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Move(wid, diff)
}

func (s *Server) Grow(diff grid.Size, _ *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Grow(wid, diff)
}

func (s *Server) Throw(direction grid.Direction, _ *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Throw(wid, direction)
}

func (s *Server) Spread(direction grid.Direction, _ *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Spread(wid, direction)
}

func (s *Server) Teleport(rectangle Rectangle, _ *struct{}) error {
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	return s.grid.Teleport(wid, rectangle.TopLeft, rectangle.BottomRight)
}
