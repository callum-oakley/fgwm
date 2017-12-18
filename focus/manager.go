package focus

import (
	"fmt"
	"time"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Manager interface {
	Register(wid wmutils.WindowID) error
	Unregister(wid wmutils.WindowID)
	Focus(wid wmutils.WindowID) error
	Unfocus(wid wmutils.WindowID) error
	FocusNext() error
	FocusPrev() error
}

type manager struct {
	// TODO mutex?
	// TODO clean this up on window deletion
	wids                             []wmutils.WindowID // The window stack, in mru order
	timer                            *time.Timer
	timeout                          time.Duration
	focussedColour, unfocussedColour wmutils.Colour
}

func NewManager(
	timeout time.Duration,
	focussedColour, unfocussedColour wmutils.Colour,
) Manager {
	m := manager{
		timeout:          timeout,
		focussedColour:   focussedColour,
		unfocussedColour: unfocussedColour,
	}
	m.timer = time.AfterFunc(m.timeout, m.update)
	return &m
}

func (m *manager) update() {
	if len(m.wids) == 0 {
		return
	}
	wid, err := wmutils.Focussed()
	if err != nil {
		return
	}
	i := index(wid, m.wids)
	if i < 0 {
		return
	}
	for j := i; j > 0; j-- {
		m.wids[j] = m.wids[j-1]
	}
	m.wids[0] = wid
}

func (m *manager) Register(wid wmutils.WindowID) error {
	if index(wid, m.wids) < 0 {
		m.wids = append([]wmutils.WindowID{wid}, m.wids...)
		return m.Focus(wid)
	}
	return nil
}

func (m *manager) Unregister(wid wmutils.WindowID) {
	if i := index(wid, m.wids); i >= 0 {
		m.wids = append(m.wids[:i], m.wids[i+1:]...)
	}
}

func (m *manager) Focus(wid wmutils.WindowID) error {
	m.timer.Stop()
	if w, err := wmutils.Focussed(); err == nil {
		if err := wmutils.SetBorderColour(w, m.unfocussedColour); err != nil {
			return err
		}
	}
	if err := wmutils.Focus(wid); err != nil {
		return err
	}
	if err := wmutils.Raise(wid); err != nil {
		return err
	}
	if err := wmutils.SetBorderColour(wid, m.focussedColour); err != nil {
		return err
	}
	m.timer.Reset(m.timeout)
	return nil
}

func (m *manager) Unfocus(wid wmutils.WindowID) error {
	w, err := wmutils.Focussed()
	if err != nil || w == wid {
		return m.focusTop()
	}
	return nil
}

func (m *manager) focusTop() error {
	if len(m.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	for i := 0; i < len(m.wids); i++ {
		if visible[m.wids[i]] {
			return m.Focus(m.wids[i])
		}
	}
	return nil
}

func (m *manager) focusFunc(f func(int) int) error {
	if len(m.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	wid, err := wmutils.Focussed()
	if err != nil {
		return err
	}
	i := index(wid, m.wids)
	if i < 0 {
		return fmt.Errorf("can't find window with id %v", wid)
	}
	for j := 0; j == 0 || !visible[m.wids[i]]; j++ {
		i = f(i)
		if j == len(m.wids) {
			return nil
		}
	}
	return m.Focus(m.wids[i])
}

func (m *manager) FocusNext() error {
	return m.focusFunc(func(i int) int {
		return (i + 1) % len(m.wids)
	})
}

func (m *manager) FocusPrev() error {
	return m.focusFunc(func(i int) int {
		return (i + len(m.wids) - 1) % len(m.wids)
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
