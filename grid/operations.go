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
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	size := r.Size()
	offset := Size{size.W / 2, size.H / 2}
	return g.Teleport(wid, Rectangle{
		center.Offset(offset.Scale(-1)),
		center.Offset(offset),
	})
}

func (g *Grid) Move(wid wmutils.WindowID, diff Size) error {
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	return g.Teleport(wid, r.Offset(diff))
}

func (g *Grid) Grow(wid wmutils.WindowID, diff Size) error {
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	if rg := r.Grow(diff); g.inGrid(rg) {
		return g.Teleport(wid, rg)
	}
	if rg := r.Grow(diff).Offset(diff); g.inGrid(rg) {
		return g.Teleport(wid, rg)
	}
	if rg := r.Grow(diff).Offset(diff.Scale(-1)); g.inGrid(rg) {
		return g.Teleport(wid, rg)
	}
	return nil
}

func (g *Grid) Throw(wid wmutils.WindowID, direction Direction) error {
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	size := r.Size()
	switch direction {
	case Left:
		return g.Teleport(wid, Rectangle{
			Position{0, r.TopLeft.Y},
			Position{size.W, r.BottomRight.Y},
		})
	case Right:
		return g.Teleport(wid, Rectangle{
			Position{g.size.W - size.W, r.TopLeft.Y},
			Position{g.size.W, r.BottomRight.Y},
		})
	case Up:
		return g.Teleport(wid, Rectangle{
			Position{r.TopLeft.X, 0},
			Position{r.BottomRight.X, size.H},
		})
	case Down:
		return g.Teleport(wid, Rectangle{
			Position{r.TopLeft.X, g.size.H - size.H},
			Position{r.BottomRight.X, g.size.H},
		})
	default:
		return fmt.Errorf("Unsupported direction '%v'", direction)
	}
}

func (g *Grid) Spread(wid wmutils.WindowID, direction Direction) error {
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	switch direction {
	case Left:
		return g.Teleport(wid, Rectangle{
			Position{0, r.TopLeft.Y},
			r.BottomRight,
		})
	case Right:
		return g.Teleport(wid, Rectangle{
			r.TopLeft,
			Position{g.size.W, r.BottomRight.Y},
		})
	case Up:
		return g.Teleport(wid, Rectangle{
			Position{r.TopLeft.X, 0},
			r.BottomRight},
		)
	case Down:
		return g.Teleport(wid, Rectangle{
			r.TopLeft,
			Position{r.BottomRight.X, g.size.H},
		})
	default:
		return fmt.Errorf("Unsupported direction '%v'", direction)
	}
}

func (g *Grid) Teleport(wid wmutils.WindowID, r Rectangle) error {
	if !g.inGrid(r) || !r.Valid() {
		return nil
	}
	return wmutils.Teleport(
		wid,
		g.pixelPosition(r.TopLeft).Offset(g.pad),
		g.pixelSize(r.Size()).Add(g.pad.Add(g.border).Scale(-2)),
	)
}