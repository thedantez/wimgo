package main

type Mode int
const (
	ModeNormal = iota	// 0
	ModeInsert			// 1
	ModeVisible			// 2
	ModeCommand			// 3
)
var currentMode = ModeNormal

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisible:
		return "VISIBLE"
	case ModeCommand:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}
