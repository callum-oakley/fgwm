// Package wmutils provides wrappers around https://github.com/wmutils
package wmutils

import (
	"fmt"
	"os/exec"
)

type WindowID uint

func (wid WindowID) String() string {
	return fmt.Sprintf("0x%08x", uint(wid))
}

type Pixels int

type Position struct {
	X Pixels
	Y Pixels
}

type Size struct {
	W Pixels
	H Pixels
}

func (p Position) Offset(s Size) Position {
	return Position{p.X + s.W, p.Y + s.H}
}

func (a Size) Add(b Size) Size {
	return Size{a.W + b.W, a.H + b.H}
}

func (a Size) Scale(k Pixels) Size {
	return Size{k * a.W, k * a.H}
}

func fetchWID(cmd *exec.Cmd) (WindowID, error) {
	var wid WindowID
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	if _, err := fmt.Fscanf(stdout, "%v", &wid); err != nil {
		return 0, err
	}
	if err := cmd.Wait(); err != nil {
		return 0, err
	}
	return wid, nil
}

// Focussed returns the WindowID of the currently focussed window. Wraps pfw.
func Focussed() (WindowID, error) {
	return fetchWID(exec.Command("pfw"))
}

// Root gets the window ID of the root window. Wraps lsw -r.
func Root() (WindowID, error) {
	return fetchWID(exec.Command("lsw", "-r"))
}

// Focus sets the keyboard input focus to the window with the given ID if it
// exists and is viewable. Wraps wtf.
func Focus(wid WindowID) error {
	return exec.Command("wtf", wid.String()).Run()
}

// Kills the window with the given ID. Wraps killw.
func Kill(wid WindowID) error {
	return exec.Command("killw", wid.String()).Run()
}

// Teleports the window with given ID to the given position, and resizes it to
// the given size. Wraps wtp.
func Teleport(wid WindowID, pos Position, size Size) error {
	return exec.Command(
		"wtp",
		fmt.Sprint(pos.X),
		fmt.Sprint(pos.Y),
		fmt.Sprint(size.W),
		fmt.Sprint(size.H),
		wid.String(),
	).Run()
}

// Map (show) the window with the given ID. Wraps mapw -m.
func Map(wid WindowID) error {
	return exec.Command("mapw", "-m", wid.String()).Run()
}

// Unmap (hide) the window with the given ID. Wraps mapw -u.
func Unmap(wid WindowID) error {
	return exec.Command("mapw", "-u", wid.String()).Run()
}

// Toggle the visibility of the window with the given ID. Wraps mapw -t.
func Toggle(wid WindowID) error {
	return exec.Command("mapw", "-t", wid.String()).Run()
}

func exitStatusOk(cmd *exec.Cmd) (bool, error) {
	err := cmd.Run()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// IsIgnored returns true if and only if the window with the given ID has the
// override_redirect attribute set. Wraps wattr o.
func IsIgnored(wid WindowID) (bool, error) {
	return exitStatusOk(exec.Command("wattr", "o", wid.String()))
}

// Exists returns true if there is a window with the given ID, false otherwise.
// Wraps wattr.
func Exists(wid WindowID) (bool, error) {
	return exitStatusOk(exec.Command("wattr", wid.String()))
}

// GetAttributes returns the size and position of the window with the given ID.
// Wraps wattr xywh.
func GetAttributes(wid WindowID) (Position, Size, error) {
	var pos Position
	var size Size
	cmd := exec.Command("wattr", "xywh", wid.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return pos, size, err
	}
	if err := cmd.Start(); err != nil {
		return pos, size, err
	}
	if _, err := fmt.Fscanf(
		stdout,
		"%v %v %v %v",
		&pos.X,
		&pos.Y,
		&size.W,
		&size.H,
	); err != nil {
		return pos, size, err
	}
	if err := cmd.Wait(); err != nil {
		return pos, size, err
	}
	return pos, size, nil
}
