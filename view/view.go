package view

import "github.com/hot-leaf-juice/fgwm/wmutils"

type View interface {
	Register(wid wmutils.WindowID)           // Register wid with the current view
	Unregister(wid wmutils.WindowID) error   // Unregister wid from the current view
	UnregisterAll(wid wmutils.WindowID)      // Unregister wid from all views
	IsRegistered(wid wmutils.WindowID) bool  // true if wid is registered in any view
	Include(wid wmutils.WindowID, n int)     // Include wid in view n
	Set(n int) error                         // Set the view to n
	Fullscreen(wid wmutils.WindowID) error   // Toggle fullsceen for wid
	Unfullscreen(wid wmutils.WindowID) error // Mark wid as not fullscreen
}

type windowState struct {
	position   wmutils.Position
	size       wmutils.Size
	fullscreen bool
}

type view struct {
	screen  wmutils.Size
	border  wmutils.Pixels
	current int
	views   map[int]map[wmutils.WindowID]*windowState
}

func New(screen wmutils.Size, border wmutils.Pixels, start int) (View, error) {
	v := view{
		screen:  screen,
		border:  border,
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
	if _, ok := v.views[n][wid]; !ok {
		v.views[n][wid] = nil
	}
}

func (v *view) Set(n int) error {
	for wid := range v.views[v.current] {
		if err := wmutils.Unmap(wid); err != nil {
			return err
		}
		if err := v.save(wid); err != nil {
			return err
		}
	}
	v.current = n
	for wid := range v.views[v.current] {
		if err := v.restore(wid); err != nil {
			return err
		}
		if err := wmutils.Map(wid); err != nil {
			return err
		}
	}
	return nil
}

func (v *view) Fullscreen(wid wmutils.WindowID) error {
	if ws := v.views[v.current][wid]; ws != nil && ws.fullscreen {
		ws.fullscreen = false
		return v.restore(wid)
	}
	if err := v.save(wid); err != nil {
		return err
	}
	v.views[v.current][wid].fullscreen = true
	if err := wmutils.SetBorderWidth(wid, 0); err != nil {
		return err
	}
	return wmutils.Teleport(wid, wmutils.Position{}, v.screen)
}

func (v *view) Unfullscreen(wid wmutils.WindowID) error {
	if ws := v.views[v.current][wid]; ws != nil && ws.fullscreen {
		ws.fullscreen = false
		if err := wmutils.SetBorderWidth(wid, v.border); err != nil {
			return err
		}
	}
	return nil
}

func (v *view) save(wid wmutils.WindowID) error {
	if ws := v.views[v.current][wid]; ws != nil && ws.fullscreen {
		return nil
	}
	position, size, err := wmutils.GetAttributes(wid)
	if err != nil {
		return err
	}
	v.views[v.current][wid] = &windowState{position, size, false}
	return nil
}

func (v *view) restore(wid wmutils.WindowID) error {
	ws := v.views[v.current][wid]
	if ws == nil {
		return nil
	}
	if ws.fullscreen {
		if err := wmutils.SetBorderWidth(wid, 0); err != nil {
			return err
		}
		return wmutils.Teleport(wid, wmutils.Position{}, v.screen)
	}
	if err := wmutils.SetBorderWidth(wid, v.border); err != nil {
		return err
	}
	return wmutils.Teleport(wid, ws.position, ws.size)
}
