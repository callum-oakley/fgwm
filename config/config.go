package config

import (
	"github.com/BurntSushi/toml"
	"time"

	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type sizeConfig struct {
	Width  int `toml:"width"`
	Height int `toml:"height"`
}

type config struct {
	Border           int        `toml:"border"`
	MinMargin        int        `toml:"min_margin"`
	Pad              int        `toml:"pad"`
	Size             sizeConfig `toml:"grid_size"`
	InitialView      int        `toml:"initial_view"`
	FocusTimeoutMs   int        `toml:"focus_timeout_ms"`
	FocussedColour   int        `toml:"focussed_colour"`
	UnfocussedColour int        `toml:"unfocussed_colour"`
}

func Load(path string) (*grid.Options, error) {
	var c config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		return nil, err
	}
	// TODO defaults?
	return &grid.Options{
		Border: wmutils.Pixels(c.Border),
		MinMargin: wmutils.Size{ // TODO allow different margins for different sides
			wmutils.Pixels(c.MinMargin),
			wmutils.Pixels(c.MinMargin),
		},
		Pad: wmutils.Size{
			wmutils.Pixels(c.Pad),
			wmutils.Pixels(c.Pad),
		},
		Size:             grid.Size{c.Size.Width, c.Size.Height},
		InitialView:      c.InitialView,
		FocusTimeout:     time.Duration(c.FocusTimeoutMs) * time.Millisecond,
		FocussedColour:   wmutils.Colour(c.FocussedColour),
		UnfocussedColour: wmutils.Colour(c.UnfocussedColour),
	}, nil
}
