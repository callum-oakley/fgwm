package grid

import (
	"fmt"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func (g *Grid) getRectangle(wid wmutils.WindowID) (Rectangle, error) {
	pPos, pSize, err := wmutils.GetAttributes(wid)
	if err != nil {
		return Rectangle{}, fmt.Errorf(
			"error getting rectangle for wid %v: %v",
			wid,
			err,
		)
	}
	return Rectangle{
		g.closestPoint(pPos.Offset(g.pad.Scale(-1))),
		g.closestPoint(pPos.Offset(pSize.Add(
			g.pad.Add(wmutils.Size{g.border, g.border}),
		))),
	}, nil
}

func (g *Grid) closestPoint(p wmutils.Position) Position {
	return Position{
		X: round((float64(p.X - g.margins.Left)) / float64(g.cell.W)),
		Y: round((float64(p.Y - g.margins.Top)) / float64(g.cell.H)),
	}
}

func round(x float64) int {
	return int(x + 0.5)
}

func (g *Grid) pInGrid(p Position) bool {
	return 0 <= p.X && p.X <= g.size.W && 0 <= p.Y && p.Y <= g.size.H
}

func (g *Grid) inGrid(r Rectangle) bool {
	return g.pInGrid(r.TopLeft) && g.pInGrid(r.BottomRight)
}

func (g *Grid) pixelSize(size Size) wmutils.Size {
	return wmutils.Size{
		W: wmutils.Pixels(size.W) * g.cell.W,
		H: wmutils.Pixels(size.H) * g.cell.H,
	}
}

func (g *Grid) pixelPosition(pos Position) wmutils.Position {
	return wmutils.Position{
		X: g.margins.Left + wmutils.Pixels(pos.X)*g.cell.W,
		Y: g.margins.Top + wmutils.Pixels(pos.Y)*g.cell.H,
	}
}

func index(wids []wmutils.WindowID, wid wmutils.WindowID) (int, error) {
	for i := 0; i < len(wids); i++ {
		if wids[i] == wid {
			return i, nil
		}
	}
	return 0, fmt.Errorf("can't find %v in %v", wid, wids)
}

func (g *Grid) centerWID(wid wmutils.WindowID) error {
	center := Position{g.size.W / 2, g.size.H / 2}
	r, err := g.getRectangle(wid)
	if err != nil {
		return nil
	}
	size := r.Size()
	offset := Size{size.W / 2, size.H / 2}
	return g.teleportWID(wid, Rectangle{
		center.Offset(offset.Scale(-1)),
		center.Offset(offset),
	})
}

func (g *Grid) teleportWID(wid wmutils.WindowID, r Rectangle) error {
	g.view.Unfullscreen(wid)
	if !g.inGrid(r) || !r.Valid() {
		return nil
	}
	return wmutils.Teleport(
		wid,
		g.pixelPosition(r.TopLeft).Offset(g.pad),
		g.pixelSize(r.Size()).Add(
			g.pad.Add(wmutils.Size{g.border, g.border}).Scale(-2),
		),
	)
}
