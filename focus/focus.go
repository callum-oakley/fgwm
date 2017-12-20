package focus

import (
	"errors"
	"time"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Focus interface {
	Register(wid wmutils.WindowID) error
	Unregister(wid wmutils.WindowID)
	Get() (wmutils.WindowID, error)
	Set(wid wmutils.WindowID) error
	Unset(wid wmutils.WindowID) error
	Next() error
	Prev() error
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
	m := focus{
		timeout:          timeout,
		focussedColour:   focussedColour,
		unfocussedColour: unfocussedColour,
	}
	m.timer = time.AfterFunc(m.timeout, m.update)
	wids, err := wmutils.List()
	if err != nil {
		return nil, err
	}
	for wid := range wids {
		if err := m.Register(wid); err != nil {
			return nil, err
		}
	}
	return &m, nil
}

func (m *focus) update() {
	wid, err := m.Get()
	if err != nil {
		return
	}
	for j := m.i; j > 0; j-- {
		m.wids[j] = m.wids[j-1]
	}
	m.i = 0
	m.wids[0] = wid
}

func (m *focus) Get() (wmutils.WindowID, error) {
	if len(m.wids) == 0 {
		return 0, errors.New("No windows")
	}
	return m.wids[m.i], nil
}

func (m *focus) Register(wid wmutils.WindowID) error {
	if index(wid, m.wids) < 0 {
		m.wids = append([]wmutils.WindowID{wid}, m.wids...)
		return m.Set(wid)
	}
	return nil
}

func (m *focus) Unregister(wid wmutils.WindowID) {
	if i := index(wid, m.wids); i >= 0 {
		m.wids = append(m.wids[:i], m.wids[i+1:]...)
		if i <= m.i {
			m.i--
		}
	}
}

func (m *focus) Set(wid wmutils.WindowID) error {
	m.timer.Stop()
	for j := 0; j < len(m.wids); j++ {
		if j != m.i {
			err := wmutils.SetBorderColour(m.wids[j], m.unfocussedColour)
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
	if err := wmutils.SetBorderColour(wid, m.focussedColour); err != nil {
		return err
	}
	m.timer.Reset(m.timeout)
	return nil
}

func (m *focus) Unset(wid wmutils.WindowID) error {
	w, err := m.Get()
	if err != nil || w == wid {
		return m.focusTop()
	}
	return nil
}

func (m *focus) focusTop() error {
	if len(m.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	for m.i = 0; m.i < len(m.wids); m.i++ {
		if visible[m.wids[m.i]] {
			return m.Set(m.wids[m.i])
		}
	}
	return nil
}

func (m *focus) focusFunc(f func(int) int) error {
	if len(m.wids) == 0 {
		return nil
	}
	visible, err := wmutils.List()
	if err != nil {
		return err
	}
	for j := 0; j == 0 || !visible[m.wids[m.i]]; j++ {
		m.i = f(m.i)
		if j == len(m.wids) {
			return nil
		}
	}
	return m.Set(m.wids[m.i])
}

func (m *focus) Next() error {
	return m.focusFunc(func(i int) int {
		return (i + 1) % len(m.wids)
	})
}

func (m *focus) Prev() error {
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
