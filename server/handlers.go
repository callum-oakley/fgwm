package server

import "github.com/hot-leaf-juice/fgwm/grid"

func (s *Server) Snap(struct{}, *struct{}) error {
	return s.grid.Snap()
}

func (s *Server) Center(struct{}, *struct{}) error {
	return s.grid.Center()
}

func (s *Server) Fullscreen(struct{}, *struct{}) error {
	return s.grid.Fullscreen()
}

func (s *Server) Kill(struct{}, *struct{}) error {
	return s.grid.Kill()
}

func (s *Server) Move(diff grid.Size, _ *struct{}) error {
	return s.grid.Move(diff)
}

func (s *Server) Grow(diff grid.Size, _ *struct{}) error {
	return s.grid.Grow(diff)
}

func (s *Server) Throw(direction grid.Direction, _ *struct{}) error {
	return s.grid.Throw(direction)
}

func (s *Server) Spread(direction grid.Direction, _ *struct{}) error {
	return s.grid.Spread(direction)
}

func (s *Server) Focus(strategy grid.FocusStrategy, _ *struct{}) error {
	return s.grid.Focus(strategy)
}

func (s *Server) Teleport(rectangle grid.Rectangle, _ *struct{}) error {
	return s.grid.Teleport(rectangle)
}
