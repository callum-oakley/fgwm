package view

import "github.com/hot-leaf-juice/fgwm/wmutils"

type View interface {
	Register(wid wmutils.WindowID)
	Unregister(wid wmutils.WindowID) error
	UnregisterAll(wid wmutils.WindowID)
	IsRegistered(wid wmutils.WindowID) bool
	Include(wid wmutils.WindowID, n int) // Include wid in view n
	Set(n int) error                     // Set the view to n
}

type windowState struct {
	position wmutils.Position
	size     wmutils.Size
}

type view struct {
	current int
	views   map[int]map[wmutils.WindowID]*windowState
}

func New(start int) (View, error) {
	v := view{
		current: start,
		views: map[int]map[wmutils.WindowID]*windowState{
			start: map[wmutils.WindowID]*windowState{},
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

func (v *view) Unregister(wid wmutils.WindowID) error {
	delete(v.views[v.current], wid)
	return wmutils.Unmap(wid)
}

func (v *view) UnregisterAll(wid wmutils.WindowID) {
	for _, wids := range v.views {
		delete(wids, wid)
	}
}

func (v *view) IsRegistered(wid wmutils.WindowID) bool {
	for _, wids := range v.views {
		if _, ok := wids[wid]; ok {
			return true
		}
	}
	return false
}

func (v *view) Include(wid wmutils.WindowID, n int) {
	if _, ok := v.views[n]; !ok {
		v.views[n] = map[wmutils.WindowID]*windowState{}
	}
	v.views[n][wid] = nil
}

func (v *view) Set(n int) error {
	for wid := range v.views[v.current] {
		if err := wmutils.Unmap(wid); err != nil {
			return err
		}
		position, size, err := wmutils.GetAttributes(wid)
		if err != nil {
			return err
		}
		v.views[v.current][wid] = &windowState{position, size}
	}
	v.current = n
	for wid, ws := range v.views[v.current] {
		if ws != nil {
			if err := wmutils.Teleport(wid, ws.position, ws.size); err != nil {
				return err
			}
		}
		if err := wmutils.Map(wid); err != nil {
			return err
		}
	}
	return nil
}
