// Package wm provides wrappers around https://github.com/wmutils
package wm

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

// Focussed returns the WindowID of the currently focussed window. Wraps pfw.
func Focussed() (WindowID, error) {
	var wid WindowID
	cmd := exec.Command("pfw")
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
