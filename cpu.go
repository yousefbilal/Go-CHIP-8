package main

import (
	"fmt"
)

type CPU struct {
	V        [16]byte
	I        uint16 //12-bits
	PC       uint16 //12-bits
	SP       uint16
	memory   *Memory
	timers   *Timers
	gfx      [64 * 32]byte
	DrawFlag bool
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
	c.PC += 2
	//decode
	instruction, err := c.decode(opcode)

	if err != nil {
		panic(err)
	}
	//execute
	instruction()

	if c.timers.delayTimer > 0 {
		c.timers.delayTimer--
	}

	if c.timers.soundTimer > 0 {
		c.timers.soundTimer--
	} else {
		fmt.Println("BEEP")
	}
}

func (c *CPU) decode(opcode uint16) (func(), error) {
	switch opcode & 0xF000 {
	case 0x2000:
		return c._2NNN(opcode), nil
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0004:
			return c._8XY4(opcode), nil
		default:
			return nil, fmt.Errorf("unknown opcode in 0x8000 series: %x", opcode)
		}
	case 0xA000:
		return c.ANNN(opcode), nil
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0033:
			return c.FX33(opcode), nil
		default:
			return nil, fmt.Errorf("unknown opcode in 0xF000 series: %x", opcode)
		}
	default:
		return nil, fmt.Errorf("unknown opcode: %x", opcode)
	}
}
func (c *CPU) _2NNN(opcode uint16) func() {
	return func() {
		c.push(c.PC)
		c.PC = opcode & 0x0FFF
	}
}

func (c *CPU) ANNN(opcode uint16) func() {
	return func() {
		c.I = opcode & 0x0FFF
	}
}

func (c *CPU) push(val uint16) {
	c.memory.stack[c.SP] = val
	c.SP++
}

func (c *CPU) _8XY4(opcode uint16) func() {
	return func() {
		if c.V[SelectNibble(opcode, 2)] > (0xFF - c.V[SelectNibble(opcode, 1)]) {
			c.V[0xF] = 1
		} else {
			c.V[0xF] = 1
		}
		c.V[SelectNibble(opcode, 2)] += c.V[SelectNibble(opcode, 1)]
	}
}

func (c *CPU) FX33(opcode uint16) func() {
	return func() {
		regVal := c.V[SelectNibble(opcode, 2)]
		c.memory.memory[c.I] = regVal / 100
		c.memory.memory[c.I+1] = (regVal / 10) % 10
		c.memory.memory[c.I+2] = regVal % 10
	}
}

func (c *CPU) DXYN(opcode uint16) func() {
	return func() {
		x := SelectNibble(opcode, 2)
		y := SelectNibble(opcode, 1)
		height := SelectNibble(opcode, 0)

		c.V[0xF] = 0
		for _y := uint16(0); _y < height; _y++ {
			pixels := c.memory.memory[c.I+_y]
			for _x := uint16(0); _x < 8; _x++ {
				if (pixels&(0x80>>_x)) != 0 && c.gfx[(y+_y)*64+x+_x] == 1 {
					c.V[0xF] = 1
				}
				c.gfx[(y+_y)*64+x+_x] ^= ((pixels & (0x80 >> _x)) >> (7 - _x))
			}
		}
		c.DrawFlag = true
	}
}
