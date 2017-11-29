package grid

import (
	"fmt"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func (g *Grid) Snap(wid wmutils.WindowID) error {
	return g.Move(wid, Size{0, 0})
}

func (g *Grid) Center(wid wmutils.WindowID) error {
	center := Position{g.size.W / 2, g.size.H / 2}
	tl, br, err := g.getAttributes(wid)
	if err != nil {
		return err
	}
	size := br.Diff(tl)
	offset := Size{size.W / 2, size.H / 2}
	return g.Teleport(
		wid,
		center.Offset(offset.Scale(-1)),
		center.Offset(offset),
	)
}

func (g *Grid) Move(wid wmutils.WindowID, diff Size) error {
	tl, br, err := g.getAttributes(wid)
	if err != nil {
		return err
	}
	return g.Teleport(wid, tl.Offset(diff), br.Offset(diff))
}

func (g *Grid) Grow(wid wmutils.WindowID, diff Size) error {
	tl, br, err := g.getAttributes(wid)
	if err != nil {
		return err
	}
	if g.inGrid(tl.Offset(diff.Scale(-1))) && g.inGrid(br.Offset(diff)) {
		return g.Teleport(wid, tl.Offset(diff.Scale(-1)), br.Offset(diff))
	}
	if g.inGrid(tl) && g.inGrid(br.Offset(diff.Scale(2))) {
		return g.Teleport(wid, tl, br.Offset(diff.Scale(2)))
	}
	if g.inGrid(tl.Offset(diff.Scale(-2))) && g.inGrid(br) {
		return g.Teleport(wid, tl.Offset(diff.Scale(-2)), br)
	}
	return nil
}

func (g *Grid) Throw(wid wmutils.WindowID, direction Direction) error {
	tl, br, err := g.getAttributes(wid)
	if err != nil {
		return err
	}
	size := br.Diff(tl)
	switch direction {
	case Left:
		return g.Teleport(wid, Position{0, tl.Y}, Position{size.W, br.Y})
	case Right:
		return g.Teleport(
			wid,
			Position{g.size.W - size.W, tl.Y},
			Position{g.size.W, br.Y},
		)
	case Up:
		return g.Teleport(wid, Position{tl.X, 0}, Position{br.X, size.H})
	case Down:
		return g.Teleport(
			wid,
			Position{tl.X, g.size.H - size.H},
			Position{br.X, g.size.H},
		)
	default:
		return fmt.Errorf("Unsupported direction '%v'", direction)
	}
}

func (g *Grid) Spread(wid wmutils.WindowID, direction Direction) error {
	tl, br, err := g.getAttributes(wid)
	if err != nil {
		return err
	}
	switch direction {
	case Left:
		return g.Teleport(wid, Position{0, tl.Y}, br)
	case Right:
		return g.Teleport(wid, tl, Position{g.size.W, br.Y})
	case Up:
		return g.Teleport(wid, Position{tl.X, 0}, br)
	case Down:
		return g.Teleport(wid, tl, Position{br.X, g.size.H})
	default:
		return fmt.Errorf("Unsupported direction '%v'", direction)
	}
}

func (g *Grid) Teleport(wid wmutils.WindowID, tl Position, br Position) error {
	if !g.inGrid(tl) || !g.inGrid(br) || tl.X >= br.X || tl.Y >= br.Y {
		return nil
	}
	return wmutils.Teleport(
		wid,
		g.pixelPosition(tl).Offset(g.pad),
		g.pixelSize(br.Diff(tl)).Add(g.pad.Add(g.border).Scale(-2)),
	)
}
