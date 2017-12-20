package view

import "github.com/hot-leaf-juice/fgwm/wmutils"

type View interface {
	Register(wid wmutils.WindowID)
	Unregister(wid wmutils.WindowID)
	Include(wid wmutils.WindowID, n int) // Include wid in view n
	Set(n int) error                     // Set the view to n
}

type view struct {
	current int
	views   map[int]map[wmutils.WindowID]bool
}

func New(start int) (View, error) {
	v := view{
		current: start,
		views: map[int]map[wmutils.WindowID]bool{
			start: map[wmutils.WindowID]bool{},
		},
	}
	wids, err := wmutils.List()
	if err != nil {
		return nil, err
	}
	for wid := range wids {
		v.Register(wid)
	}
	return &v, nil
}

func (v *view) Register(wid wmutils.WindowID) {
	v.Include(wid, v.current)
}

func (v *view) Unregister(wid wmutils.WindowID) {
	for _, wids := range v.views {
		delete(wids, wid)
	}
}

func (v *view) Include(wid wmutils.WindowID, n int) {
	if _, ok := v.views[n]; !ok {
		v.views[n] = map[wmutils.WindowID]bool{}
	}
	v.views[n][wid] = true
}

func (v *view) Set(n int) error {
	v.current = n
	for _, wids := range v.views {
		if err := forEach(wmutils.Unmap, wids); err != nil {
			return err
		}
	}
	if err := forEach(wmutils.Map, v.views[v.current]); err != nil {
		return err
	}
	return nil
}

func forEach(
	f func(wmutils.WindowID) error,
	wids map[wmutils.WindowID]bool,
) error {
	for wid := range wids {
		if err := f(wid); err != nil {
			return err
		}
	}
	return nil
}
