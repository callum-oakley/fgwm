package main

import (
	"log"

	"github.com/hot-leaf-juice/fgwm/grid"
	"github.com/hot-leaf-juice/fgwm/wmutils"
)

func main() {
	g, err := grid.New(&grid.Options{
		Border:    wmutils.Size{5, 5},
		MinMargin: wmutils.Size{10, 10},
		Pad:       wmutils.Size{10, 10},
		Size:      grid.Size{24, 24},
	})
	if err != nil {
		log.Fatal(err)
	}
	wid, err := wmutils.Focussed()
	if err != nil {
		log.Fatal(err)
	}
	err = g.Snap(wid)
	if err != nil {
		log.Fatal(err)
	}
	// for ev := range wmutils.WatchEvents() {
	// 	switch ev.Type {
	// 	case wmutils.CreateNotifyEvent:
	// 		fmt.Printf("Window created:   %v\n", ev.WID)
	// 	case wmutils.DestroyNotifyEvent:
	// 		fmt.Printf("Window destroyed: %v\n", ev.WID)
	// 	case wmutils.UnmapNotifyEvent:
	// 		fmt.Printf("Window unmapped:  %v\n", ev.WID)
	// 	case wmutils.MapNotifyEvent:
	// 		fmt.Printf("Window mapped:    %v\n", ev.WID)
	// 	}
	// }
	// log.Fatal("Something went wrong! Event channel was closed...")
}
