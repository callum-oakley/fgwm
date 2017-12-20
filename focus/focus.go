package focus

import (
	"fmt"
	"time"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Focus interface {
	Register(wid wmutils.WindowID) error
	Unregister(wid wmutils.WindowID) error
	Get() (wmutils.WindowID, error)
	Set(wid wmutils.WindowID) error
	Unset(wid wmutils.WindowID) error
	Next() error
	Prev() error
	Top() error
}

type focus struct {
	// TODO mutex?
	// TODO clean this up on window deletion
	wids                             []wmutils.WindowID
	i                                int
	timer                            *time.Timer
	timeout                          time.Duration
	focussedColour, unfocussedColour wmutils.Colour
}

func New(
	timeout time.Duration,
	focussedColour, unfocussedColour wmutils.Colour,
) (Focus, error) {
	f := focus{
		timeout:          timeout,
		focussedColour:   focussedColour,
		unfocussedColour: unfocussedColour,
	}
	f.timer = time.AfterFunc(f.timeout, f.update)
	wids, err := wmutils.List()
	if err != nil {
		return nil, err
	}
	for wid := range wids {
		if err := f.Register(wid); err != nil {
			return nil, err
		}
	}
	return &f, nil
}

func (f *focus) update() {
	wid, err := f.Get()
	if err != nil {
		return
	}
	for j := f.i; j > 0; j-- {
		f.wids[j] = f.wids[j-1]
	}
	f.i = 0
	f.wids[0] = wid
}

func (f *focus) Get() (wmutils.WindowID, error) {
	if f.i >= len(f.wids) {
		return 0, fmt.Errorf(
			"index is %v but we only have %v wids!",
			f.i,
			len(f.wids),
		)
	}
	return f.wids[f.i], nil
}

func (f *focus) Register(wid wmutils.WindowID) error {
	if index(wid, f.wids) < 0 {
		f.wids = append([]wmutils.WindowID{wid}, f.wids...)
		f.i = 0
		return f.Set(wid)
	}
	return nil
}

func (f *focus) Unregister(wid wmutils.WindowID) error {
	if i := index(wid, f.wids); i >= 0 {
		f.wids = append(f.wids[:i], f.wids[i+1:]...)
		return f.Top()
	}
	return nil
}

func (f *focus) Set(wid wmutils.WindowID) error {
	f.timer.Stop()
	for j := 0; j < len(f.wids); j++ {
		if j != f.i {
			err := wmutils.SetBorderColour(f.wids[j], f.unfocussedColour)
			if err != nil {
				return err
			}
		}
	}
	if err := wmutils.Focus(wid); err != nil {
		return err
	}
	if err := wmutils.Raise(wid); err != nil {
		return err
	}
	if err := wmutils.SetBorderColour(wid, f.focussedColour); err != nil {
		return err
	}
	f.timer.Reset(f.timeout)
	return nil
}

func (f *focus) Unset(wid wmutils.WindowID) error {
	w, err := f.Get()
	if err != nil || w == wid {
		return f.Top()
	}
	return nil
}

func (f *focus) Top() error {
	if len(f.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	for f.i = 0; f.i < len(f.wids); f.i++ {
		if visible[f.wids[f.i]] {
			return f.Set(f.wids[f.i])
		}
	}
	return nil
}

func (f *focus) focusFunc(g func(int) int) error {
	if len(f.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	for j := 0; j == 0 || !visible[f.wids[f.i]]; j++ {
		f.i = g(f.i)
		if j == len(f.wids) {
			return nil
		}
	}
	return f.Set(f.wids[f.i])
}

func (f *focus) Next() error {
	return f.focusFunc(func(i int) int {
		return (i + 1) % len(f.wids)
	})
}

func (f *focus) Prev() error {
	return f.focusFunc(func(i int) int {
		return (i + len(f.wids) - 1) % len(f.wids)
	})
}

func index(wid wmutils.WindowID, wids []wmutils.WindowID) int {
	for i, w := range wids {
		if w == wid {
			return i
		}
	}
	return -1
}
