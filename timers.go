package main

type Timers struct {
	delayTimer, soundTimer uint8
}

func NewTimers() *Timers {
	return &Timers{}
}
