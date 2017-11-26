package wmutils

import (
	"fmt"
	"os/exec"
)

type EventType uint

type Event struct {
	Type EventType
	WID  WindowID
}

const (
	CreateNotifyEvent EventType = 16 + iota
	DestroyNotifyEvent
	UnmapNotifyEvent
	MapNotifyEvent
)

// this doesn't do any cleanup yet...
func WatchEvents() <-chan Event {
	evChan := make(chan Event)
	// might want to use CommandContext to kill this and clean up when done
	cmd := exec.Command("wew")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		close(evChan)
		return evChan
	}
	if err := cmd.Start(); err != nil {
		close(evChan)
		return evChan
	}
	go func() {
		for {
			var ev Event
			_, err := fmt.Fscanf(stdout, "%v:%v", &ev.Type, &ev.WID)
			if err != nil {
				close(evChan)
				break
			}
			isIgnored, err := IsIgnored(ev.WID)
			if err != nil {
				close(evChan)
				break
			}
			if !isIgnored {
				evChan <- ev
			}
		}
	}()
	return evChan
}
