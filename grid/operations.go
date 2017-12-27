package grid

import (
	"fmt"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func (g *Grid) Focus(strategy FocusStrategy) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	switch strategy {
	case Next:
		return g.focus.Next()
	case Prev:
		return g.focus.Prev()
	default:
		return fmt.Errorf("unsupported focus strategy '%v'", strategy)
	}
}

func (g *Grid) ViewInclude(n int) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	g.view.Include(wid, n)
	return nil
}

func (g *Grid) ViewSet(n int) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	if err := g.view.Set(n); err != nil {
		return err
	}
	return g.focus.Top()
}

func (g *Grid) Snap() error {
	return g.Move(Size{0, 0})
}

func (g *Grid) Center() error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	return g.centerWID(wid)
}

func (g *Grid) Fullscreen() error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	return g.view.Fullscreen(wid)
}

func (g *Grid) Kill() error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	if err := g.view.Unregister(wid); err != nil {
		return err
	}
	if g.view.IsRegistered(wid) {
		return nil
	}
	return wmutils.Kill(wid)
}

func (g *Grid) Move(diff Size) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	return g.teleportWID(wid, r.Offset(diff))
}

func (g *Grid) Grow(diff Size) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	if rg := r.Grow(diff); g.inGrid(rg) {
		return g.teleportWID(wid, rg)
	}
	if rg := r.Grow(diff).Offset(diff); g.inGrid(rg) {
		return g.teleportWID(wid, rg)
	}
	if rg := r.Grow(diff).Offset(diff.Scale(-1)); g.inGrid(rg) {
		return g.teleportWID(wid, rg)
	}
	return nil
}

func (g *Grid) Throw(direction Direction) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	size := r.Size()
	switch direction {
	case Left:
		return g.teleportWID(wid, Rectangle{
			Position{0, r.TopLeft.Y},
			Position{size.W, r.BottomRight.Y},
		})
	case Right:
		return g.teleportWID(wid, Rectangle{
			Position{g.size.W - size.W, r.TopLeft.Y},
			Position{g.size.W, r.BottomRight.Y},
		})
	case Up:
		return g.teleportWID(wid, Rectangle{
			Position{r.TopLeft.X, 0},
			Position{r.BottomRight.X, size.H},
		})
	case Down:
		return g.teleportWID(wid, Rectangle{
			Position{r.TopLeft.X, g.size.H - size.H},
			Position{r.BottomRight.X, g.size.H},
		})
	default:
		return fmt.Errorf("unsupported direction '%v'", direction)
	}
}

func (g *Grid) Spread(direction Direction) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	r, err := g.getRectangle(wid)
	if err != nil {
		return err
	}
	switch direction {
	case Left:
		return g.teleportWID(wid, Rectangle{
			Position{0, r.TopLeft.Y},
			r.BottomRight,
		})
	case Right:
		return g.teleportWID(wid, Rectangle{
			r.TopLeft,
			Position{g.size.W, r.BottomRight.Y},
		})
	case Up:
		return g.teleportWID(wid, Rectangle{
			Position{r.TopLeft.X, 0},
			r.BottomRight},
		)
	case Down:
		return g.teleportWID(wid, Rectangle{
			r.TopLeft,
			Position{r.BottomRight.X, g.size.H},
		})
	default:
		return fmt.Errorf("Unsupported direction '%v'", direction)
	}
}

func (g *Grid) Teleport(r Rectangle) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	wid, err := g.focus.Get()
	if err != nil {
		return err
	}
	return g.teleportWID(wid, r)
}
