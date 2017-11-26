package grid

import (
	"fmt"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Position struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

func (p Position) Offset(s Size) Position {
	return Position{p.X + s.W, p.Y + s.H}
}

func (a Position) Diff(b Position) Size {
	return Size{a.X - b.X, a.Y - b.Y}
}

func (a Size) Add(b Size) Size {
	return Size{a.W + b.W, a.H + b.H}
}

func (a Size) Scale(k int) Size {
	return Size{k * a.W, k * a.H}
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
	// the size of the grid in cells
	size Size
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
		size:   opts.Size,
	}, nil
}

func (g *Grid) closestPoint(p wmutils.Position) Position {
	return Position{
		X: int((p.X - g.margin.W + g.cell.W/2) / g.cell.W),
		Y: int((p.Y - g.margin.H + g.cell.H/2) / g.cell.H),
	}
}

func (g *Grid) GetAttributes(wid wmutils.WindowID) (Position, Position, error) {
	pPos, pSize, err := wmutils.GetAttributes(wid)
	if err != nil {
		return Position{}, Position{}, err
	}
	return g.closestPoint(pPos), g.closestPoint(pPos.Offset(pSize)), nil
}

func (g *Grid) Move(wid wmutils.WindowID, diff Size) error {
	tl, br, err := g.GetAttributes(wid)
	if err != nil {
		return err
	}
	return g.Teleport(wid, tl.Offset(diff), br.Offset(diff))
}

func (g *Grid) Grow(wid wmutils.WindowID, diff Size) error {
	tl, br, err := g.GetAttributes(wid)
	if err != nil {
		return err
	}
	return g.Teleport(wid, tl.Offset(diff.Scale(-1)), br.Offset(diff))
}

func (g *Grid) Snap(wid wmutils.WindowID) error {
	return g.Move(wid, Size{0, 0})
}

func restriction(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func (g *Grid) pixelSize(size Size) wmutils.Size {
	return wmutils.Size{
		W: wmutils.Pixels(size.W) * g.cell.W,
		H: wmutils.Pixels(size.H) * g.cell.H,
	}
}

func (g *Grid) Center(wid wmutils.WindowID) error {
	center := Position{g.size.W / 2, g.size.H / 2}
	tl, br, err := g.GetAttributes(wid)
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

func (g *Grid) Throw(wid wmutils.WindowID, direction Direction) error {
	tl, br, err := g.GetAttributes(wid)
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
		return fmt.Errorf("Unknown direction '%v'", direction)
	}
}

func (g *Grid) inGrid(p Position) bool {
	return 0 <= p.X && p.X <= g.size.W && 0 <= p.Y && p.Y <= g.size.H
}

func (g *Grid) Teleport(wid wmutils.WindowID, tl Position, br Position) error {
	if !g.inGrid(tl) || !g.inGrid(br) {
		return nil
	}
	return wmutils.Teleport(
		wid,
		g.points[tl].Offset(g.pad),
		g.pixelSize(br.Diff(tl)).Add(g.pad.Add(g.border).Scale(-2)),
	)
}
