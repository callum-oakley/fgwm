package config

import (
	"github.com/pelletier/go-toml"
	"time"

	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func Load(path string) (*grid.Options, error) {
	c, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	// TODO GetDefault?
	return &grid.Options{
		Border: wmutils.Pixels(c.Get("border").(int64)),
		Margins: grid.Margins{
			Top:    wmutils.Pixels(c.Get("margins.top").(int64)),
			Bottom: wmutils.Pixels(c.Get("margins.bottom").(int64)),
			Left:   wmutils.Pixels(c.Get("margins.left").(int64)),
			Right:  wmutils.Pixels(c.Get("margins.right").(int64)),
		},
		Pad: wmutils.Size{
			wmutils.Pixels(c.Get("pad.width").(int64)),
			wmutils.Pixels(c.Get("pad.height").(int64)),
		},
		Size: grid.Size{
			int(c.Get("grid_size.width").(int64)),
			int(c.Get("grid_size.height").(int64)),
		},
		InitialView: int(c.Get("initial_view").(int64)),
		FocusTimeout: time.Duration(
			c.Get("focus_timeout_ms").(int64),
		) * time.Millisecond,
		FocussedColour:   wmutils.Colour(c.Get("focussed_colour").(int64)),
		UnfocussedColour: wmutils.Colour(c.Get("unfocussed_colour").(int64)),
	}, nil
}
