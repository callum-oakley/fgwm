package focus

import (
	"time"

	"github.com/hot-leaf-juice/fgwm/wmutils"
)

type Manager interface {
	FocusWID(wid wmutils.WindowID) error
	FocusNext() error
	FocusPrev() error
}

type manager struct {
	// TODO mutex
	// TODO clean this up on window deletion
	wids    []wmutils.WindowID // The window stack, in mru order
	i       int                // The current location in the window stack
	timer   *time.Timer
	timeout time.Duration
}

func NewManager(timeout time.Duration) Manager {
	m := manager{timeout: timeout}
	m.timer = time.AfterFunc(m.timeout, m.update)
	// TODO remove ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// We shouldn't need to do this once we are focussing windows on creation,
	// this is a hack so that we can test the behaviour alongside windowchef...
	wids, _ := wmutils.List()
	for wid := range wids {
		m.wids = append(m.wids, wid)
	}
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	return &m
}

func (m *manager) update() {
	if len(m.wids) == 0 {
		return
	}
	updated := []wmutils.WindowID{m.wids[m.i]}
	for i, wid := range m.wids {
		if i != m.i {
			updated = append(updated, wid)
		}
	}
	m.wids = updated
	m.i = 0
}

func (m *manager) focus() error {
	// TODO border colour
	m.timer.Stop()
	defer m.timer.Reset(m.timeout)
	wid := m.wids[m.i]
	if err := wmutils.Focus(wid); err != nil {
		return err
	}
	return wmutils.Raise(wid)
}

func (m *manager) FocusWID(wid wmutils.WindowID) error {
	i := index(wid, m.wids)
	if i < 0 {
		m.wids = append(m.wids, wid)
		m.i = len(m.wids) - 1
	}
	m.i = i
	return m.focus()
}

func (m *manager) focusFunc(f func(int) int) error {
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
	return m.focus()
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
