package config

import (
	"github.com/pelletier/go-toml"
	"os"
	"time"

	"github.com/callum-oakley/fgwm/grid"
	"github.com/callum-oakley/fgwm/wmutils"
)

func Load(path string) (*grid.Options, error) {
	c, err := toml.LoadFile(path)
	if os.IsNotExist(err) {
		c = &toml.Tree{} // Will just revert to defaults for everything.
	} else if err != nil {
		return nil, err
	}
	return &grid.Options{
		Border: wmutils.Pixels(c.GetDefault("border", int64(5)).(int64)),
		Margins: grid.Margins{
			Top: wmutils.Pixels(
				c.GetDefault("margins.top", int64(10)).(int64),
			),
			Bottom: wmutils.Pixels(
				c.GetDefault("margins.bottom", int64(10)).(int64),
			),
			Left: wmutils.Pixels(
				c.GetDefault("margins.left", int64(10)).(int64),
			),
			Right: wmutils.Pixels(
				c.GetDefault("margins.right", int64(10)).(int64),
			),
		},
		Pad: wmutils.Size{
			wmutils.Pixels(c.GetDefault("pad.width", int64(10)).(int64)),
			wmutils.Pixels(c.GetDefault("pad.height", int64(10)).(int64)),
		},
		Size: grid.Size{
			int(c.GetDefault("grid_size.width", int64(24)).(int64)),
			int(c.GetDefault("grid_size.height", int64(24)).(int64)),
		},
		InitialView: int(c.GetDefault("initial_view", int64(1)).(int64)),
		FocusTimeout: time.Duration(
			c.GetDefault("focus_timeout_ms", int64(500)).(int64),
		) * time.Millisecond,
		FocussedColour: wmutils.Colour(
			c.GetDefault("focussed_colour", int64(0xd8dee9)).(int64),
		),
		UnfocussedColour: wmutils.Colour(
			c.GetDefault("unfocussed_colour", int64(0x65737e)).(int64),
		),
	}, nil
}
