package main

import (
	"fmt"
	"log"

	"github.com/hot-leaf-juice/fgwm/wm"
)

func main() {
	for ev := range wm.WatchEvents() {
		switch ev.Type {
		case wm.CreateNotifyEvent:
			fmt.Printf("Window created:   %v\n", ev.WID)
		case wm.DestroyNotifyEvent:
			fmt.Printf("Window destroyed: %v\n", ev.WID)
		case wm.UnmapNotifyEvent:
			fmt.Printf("Window unmapped:  %v\n", ev.WID)
		case wm.MapNotifyEvent:
			fmt.Printf("Window mapped:    %v\n", ev.WID)
		}
	}
	log.Fatal("Something went wrong! Event channel was closed...")
}
