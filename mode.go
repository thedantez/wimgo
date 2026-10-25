package main

type Mode int
const (
	ModeNormal = iota	// 0
	ModeInsert			// 1
	ModeVisible			// 2
	ModeCommand			// 3
	ModeHints			// 4
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
	case ModeHints:
		return "HINTS"
	default:
		return "UNKNOWN"
	}
}
