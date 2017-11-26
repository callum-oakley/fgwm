package grid

import (
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Position struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

func (a Position) Diff(b Position) Size {
	return Size{a.X - b.X, a.Y - b.Y}
}

type Grid struct {
	// size of the screen
	screen wmutils.Size
	// margin at edge of screen
	margin wmutils.Size
	// padding around cells
	pad wmutils.Size
	// border around cells
	border wmutils.Size
	// size of each cell, including pad and border but excluding margin
	cell wmutils.Size
	// the pixel locations of the cell boundaries
	points map[Position]wmutils.Position
}

// The sizes that define the grid layout are made up as follows (bd is border).
// The Y direction is similar.
//
// | <-------------------------------- screen -------------------------------> |
// |        |     |    |      |    |     |     |    |      |    |     |        |
// | margin | pad | bd |      | bd | pad | pad | bd |      | bd | pad | margin |
// |        |     |    |      |    |     |     |    |      |    |     |        |
//          | <--------- cell ---------> | <--------- cell ---------> |
//          |                            |                            |
//          A                            B                            C
//
// The points attribute would list the points A, B, C above -- i.e. the cell
// boundaries.

type Options struct {
	Border    wmutils.Size
	MinMargin wmutils.Size
	Pad       wmutils.Size
	Size      Size
}

func New(opts *Options) (*Grid, error) {
	wid, err := wmutils.Root()
	if err != nil {
		return nil, err
	}
	_, screen, err := wmutils.GetAttributes(wid)
	if err != nil {
		return nil, err
	}
	cell := wmutils.Size{
		W: (screen.W - 2*opts.MinMargin.W) / wmutils.Pixels(opts.Size.W),
		H: (screen.H - 2*opts.MinMargin.H) / wmutils.Pixels(opts.Size.H),
	}
	margin := wmutils.Size{
		W: (screen.W - wmutils.Pixels(opts.Size.W)*cell.W) / 2,
		H: (screen.H - wmutils.Pixels(opts.Size.H)*cell.H) / 2,
	}
	points := make(map[Position]wmutils.Position)
	for x := 0; x <= opts.Size.W; x++ {
		for y := 0; y <= opts.Size.H; y++ {
			points[Position{x, y}] = wmutils.Position{
				X: margin.W + wmutils.Pixels(x)*cell.W,
				Y: margin.H + wmutils.Pixels(y)*cell.H,
			}
		}
	}
	return &Grid{
		screen: screen,
		margin: margin,
		border: opts.Border,
		pad:    opts.Pad,
		cell:   cell,
		points: points,
	}, nil
}

func (g *Grid) closestPoint(p wmutils.Position) Position {
	return Position{
		X: int((p.X - g.margin.W + g.cell.W/2) / g.cell.W),
		Y: int((p.Y - g.margin.H + g.cell.H/2) / g.cell.H),
	}
}

func (g *Grid) Snap(wid wmutils.WindowID) error {
	pPos, pSize, err := wmutils.GetAttributes(wid)
	if err != nil {
		return err
	}
	topRight := g.closestPoint(pPos)
	bottomLeft := g.closestPoint(pPos.Offset(pSize))
	return g.Teleport(wid, topRight, bottomLeft.Diff(topRight))
}

func (g *Grid) Teleport(wid wmutils.WindowID, pos Position, size Size) error {
	// TODO check target pos and size are inside screen
	// More reasearch required to confirm this, but the position seems to be
	// from the top left corner *including the border* but the size *excludes
	// the border*.
	return wmutils.Teleport(
		wid,
		g.points[pos].Offset(g.pad),
		g.pixelSize(size).Add(g.pad.Add(g.border).Scale(-2)),
	)
}

func (g *Grid) pixelSize(size Size) wmutils.Size {
	return wmutils.Size{
		W: wmutils.Pixels(size.W) * g.cell.W,
		H: wmutils.Pixels(size.H) * g.cell.H,
	}
}
