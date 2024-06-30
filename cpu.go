package main

type CPU struct {
	V      [16]uint16
	I      uint16 //12-bits
	PC     uint16 //12-bits
	SP     uint16
	memory *Memory
	timers *Timers
}

func NewChip8(fileName string) *CPU {
	return &CPU{
		I:      0,
		PC:     0x200, // Program counter starts at 0x200
		SP:     0,
		memory: NewMemory(fileName),
		timers: NewTimers(),
	}
}
