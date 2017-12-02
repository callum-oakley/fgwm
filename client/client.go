package client

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/server"
)

type client struct {
	name string
}

func Run(args []string) {
	conn, err := rpc.DialHTTP("tcp", "localhost:62676")
	if err != nil {
		log.Fatal(err)
	}
	c := client{args[0]}
	switch args[1] {
	case "snap":
		conn.Call("Server.Snap", c.noArgs(args[2:]), nil)
	case "center":
		conn.Call("Server.Center", c.noArgs(args[2:]), nil)
	case "kill":
		conn.Call("Server.Kill", c.noArgs(args[2:]), nil)
	case "move":
		conn.Call("Server.Move", c.sizeArg(args[2:]), nil)
	case "grow":
		conn.Call("Server.Grow", c.sizeArg(args[2:]), nil)
	case "throw":
		conn.Call("Server.Throw", c.directionArg(args[2:]), nil)
	case "spread":
		conn.Call("Server.Spread", c.directionArg(args[2:]), nil)
	case "teleport":
		conn.Call("Server.Teleport", c.rectangleArg(args[2:]), nil)
	case "help":
		c.printHelpAndExit(args[2:])
	default:
		c.printHelpAndExit(nil)
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

func (c client) rectangleArg(args []string) server.Rectangle {
	if len(args) != 4 {
		c.printHelpAndExit(nil)
	}
	var r server.Rectangle
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
