package main

import (
	"errors"
	"fmt"
)

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

func (c *CPU) EmulationCycle() {
	//fetch
	opcode := c.memory.ReadOpcode(c.PC)

	//decode
	instruction, err :=  c.decode(opcode)

	if err != nil {
		panic(err)	
	}
	//execute
	instruction()

	if c.timers.delayTimer > 0 {
		c.timers.delayTimer--
	}
	
	if c.timers.soundTimer > 0 {
		c.timers.soundTimer --
	} else {
		fmt.Println("BEEP")
	}
}

func (c *CPU) decode(opcode uint16) (func(), error) {

	switch opcode & 0xF000 {
	case 0xA000:
		return c.ANNN(opcode), nil
	default:
		return nil, fmt.Errorf("unknown opcode: %x", opcode)
	}
}

func (c *CPU) ANNN(opcode uint16) func() {
	return func() {
		c.I = opcode & 0x0FFF
		c.PC += 2
	}
}
