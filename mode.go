package main

type Mode int
const (
	ModeNormal = iota
	ModeInsert
	ModeVisible
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
	default:
		return "UNKNOWN"
	}
}
