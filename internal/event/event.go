package event

import "strings"

type HyprEvent struct {
	Event string
	Data  string
}

func NewEvent(line string) HyprEvent {
	split := strings.Split(line, ">>")

	return HyprEvent{Event: split[0], Data: split[1]}
}
