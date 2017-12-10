package client

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/hot-leaf-juice/fgwm/grid"
)

type client struct {
	name string
}

func Run(args []string) {
	c := client{args[0]}
	conn, err := rpc.DialHTTP("tcp", "localhost:62676")
	if err != nil {
		log.Fatalf("%v: %v", c.name, err)
	}
	switch args[1] {
	case "snap":
		err = conn.Call("Server.Snap", c.noArgs(args[2:]), nil)
	case "center":
		err = conn.Call("Server.Center", c.noArgs(args[2:]), nil)
	case "fullscreen":
		err = conn.Call("Server.Fullscreen", c.noArgs(args[2:]), nil)
	case "kill":
		err = conn.Call("Server.Kill", c.noArgs(args[2:]), nil)
	case "move":
		err = conn.Call("Server.Move", c.sizeArg(args[2:]), nil)
	case "grow":
		err = conn.Call("Server.Grow", c.sizeArg(args[2:]), nil)
	case "throw":
		err = conn.Call("Server.Throw", c.directionArg(args[2:]), nil)
	case "spread":
		err = conn.Call("Server.Spread", c.directionArg(args[2:]), nil)
	case "focus":
		err = conn.Call("Server.Focus", c.nextOrPrevArg(args[2:]), nil)
	case "teleport":
		err = conn.Call("Server.Teleport", c.rectangleArg(args[2:]), nil)
	case "help":
		c.printHelpAndExit(args[2:])
	default:
		c.printHelpAndExit(nil)
	}
	if err != nil {
		log.Fatalf("%v: %v", c.name, err)
	}
}

func (c client) noArgs(args []string) struct{} {
	if len(args) != 0 {
		c.printHelpAndExit(nil)
	}
	return struct{}{}
}

func (c client) sizeArg(args []string) grid.Size {
	if len(args) != 2 {
		c.printHelpAndExit(nil)
	}
	var size grid.Size
	var err error
	if size.W, err = strconv.Atoi(args[0]); err != nil {
		c.printHelpAndExit(nil)
	}
	if size.H, err = strconv.Atoi(args[1]); err != nil {
		c.printHelpAndExit(nil)
	}
	return size
}

func (c client) directionArg(args []string) grid.Direction {
	if len(args) != 1 {
		c.printHelpAndExit(nil)
	}
	var direction grid.Direction
	switch strings.ToLower(args[0]) {
	case "left", "l", "west", "w":
		direction = grid.Left
	case "right", "r", "east", "e":
		direction = grid.Right
	case "up", "u", "north", "n":
		direction = grid.Up
	case "down", "d", "south", "s":
		direction = grid.Down
	default:
		c.printHelpAndExit(nil)
	}
	return direction
}

func (c client) nextOrPrevArg(args []string) grid.NextOrPrev {
	if len(args) != 1 {
		c.printHelpAndExit(nil)
	}
	var nextOrPrev grid.NextOrPrev
	switch strings.ToLower(args[0]) {
	case "next", "n":
		nextOrPrev = grid.Next
	case "prev", "p":
		nextOrPrev = grid.Prev
	default:
		c.printHelpAndExit(nil)
	}
	return nextOrPrev
}

func (c client) rectangleArg(args []string) grid.Rectangle {
	if len(args) != 4 {
		c.printHelpAndExit(nil)
	}
	var r grid.Rectangle
	var err error
	if r.TopLeft.X, err = strconv.Atoi(args[0]); err != nil {
		c.printHelpAndExit(nil)
	}
	if r.TopLeft.Y, err = strconv.Atoi(args[1]); err != nil {
		c.printHelpAndExit(nil)
	}
	if r.BottomRight.X, err = strconv.Atoi(args[2]); err != nil {
		c.printHelpAndExit(nil)
	}
	if r.BottomRight.Y, err = strconv.Atoi(args[3]); err != nil {
		c.printHelpAndExit(nil)
	}
	return r
}

func (c client) printHelpAndExit(args []string) {
	// TODO improve this (including command specific help)
	if len(args) == 0 {
		fmt.Printf("Usage:\n\n\t%v command [arguments]\n\n", c.name)
		fmt.Println("Where command is one of:\n")
		fmt.Println(
			"\tsnap, move, grow, center, throw, spred, teleport, kill, help",
		)
	} else {
		fmt.Printf("help for %v coming soon!\n", args[0])
	}
	os.Exit(0)
}
