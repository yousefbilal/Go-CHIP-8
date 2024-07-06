package main

import (
	"fmt"
	"math/rand"
)

type CPU struct {
	V      [16]byte
	I      uint16 //12-bits
	PC     uint16 //12-bits
	SP     uint16
	memory *Memory
	timers *Timers
	gfx    [64 * 32]byte
	opcode uint16
	keys   [16]bool
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
	c.opcode = c.memory.ReadOpcode(c.PC)
	c.PC += 2
	//decode
	instruction, err := c.decode()

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
		if c.timers.soundTimer == 0 {
			fmt.Println("BEEP")
		}
	}
}

func (c *CPU) decode() (func(), error) {
	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000:
			return c._00E0, nil
		case 0x000E:
			return c._00EE, nil
		}
	case 0x1000:
		return c._1NNN, nil
	case 0x2000:
		return c._2NNN, nil
	case 0x3000:
		return c._3XKK, nil
	case 0x4000:
		return c._4XKK, nil
	case 0x5000:
		return c._5XY0, nil
	case 0x6000:
		return c._6XKK, nil
	case 0x7000:
		return c._7XKK, nil
	case 0x8000:
		switch c.opcode & 0x000F {
		case 0x0000:
			return c._8XY0, nil
		case 0x0001:
			return c._8XY1, nil
		case 0x0002:
			return c._8XY2, nil
		case 0x0003:
			return c._8XY3, nil
		case 0x0004:
			return c._8XY4, nil
		case 0x0005:
			return c._8XY5, nil
		case 0x0006:
			return c._8XY6, nil
		case 0x0007:
			return c._8XY7, nil
		case 0x000E:
			return c._8XYE, nil
		}
	case 0x9000:
		return c._9XY0, nil
	case 0xA000:
		return c.ANNN, nil
	case 0xB000:
		return c.BNNN, nil
	case 0xC000:
		return c.CXKK, nil
	case 0xD000:
		return c.DXYN, nil
	case 0xE000:
		switch c.opcode & 0x000F {
		case 0x0001:
			return c.EXA1, nil
		case 0x000E:
			return c.EX9E, nil
		}
	case 0xF000:
		switch c.opcode & 0x00FF {
		case 0x0007:
			return c.FX07, nil
		case 0x000A:
			return c.FX0A, nil
		case 0x0015:
			return c.FX15, nil
		case 0x0018:
			return c.FX18, nil
		case 0x001E:
			return c.FX1E, nil
		case 0x0029:
			return c.FX29, nil
		case 0x0033:
			return c.FX33, nil
		case 0x0055:
			return c.FX55, nil
		case 0x0065:
			return c.FX65, nil
		}
	}
	return nil, fmt.Errorf("unknown opcode: %x", c.opcode)
}
func (c *CPU) push(val uint16) {
	c.memory.stack[c.SP] = val
	c.SP++
}

func (c *CPU) pop() uint16 {
	c.SP--
	return c.memory.stack[c.SP]
}

func (c *CPU) _00E0() {
	//clear the display
	for i := 0; i < 64*32; i++ {
		c.gfx[i] = 0
	}
}

func (c *CPU) _00EE() {
	//return from subroutine
	c.PC = c.pop()
}

func (c *CPU) _1NNN() {
	//jump to nnn
	c.PC = c.opcode & 0x0FFF
}

func (c *CPU) _2NNN() {
	//call subroutine at nnn
	c.push(c.PC)
	c.PC = c.opcode & 0x0FFF
}

func (c *CPU) _3XKK() {
	//skip next instruction if Vx == kk
	if c.V[SelectNibble(c.opcode, 2)] == byte(c.opcode&0x00FF) {
		c.PC += 2
	}
}

func (c *CPU) _4XKK() {
	//skip next instruction if Vx != kk.
	if c.V[SelectNibble(c.opcode, 2)] != byte(c.opcode&0x00FF) {
		c.PC += 2
	}
}

func (c *CPU) _5XY0() {
	//skip next instruction if Vx == Vy
	if c.V[SelectNibble(c.opcode, 2)] == c.V[SelectNibble(c.opcode, 1)] {
		c.PC += 2
	}
}

func (c *CPU) _6XKK() {
	//load kk into Vx (Vx = kk)
	c.V[SelectNibble(c.opcode, 2)] = byte(c.opcode & 0x00FF)
}

func (c *CPU) _7XKK() {
	//Vx = Vx + kk
	c.V[SelectNibble(c.opcode, 2)] += byte(c.opcode & 0x00FF)
}

func (c *CPU) _8XY0() {
	//set Vx = Vy
	c.V[SelectNibble(c.opcode, 2)] = c.V[SelectNibble(c.opcode, 1)]
}

func (c *CPU) _8XY1() {
	//Vx = Vx OR Vy
	c.V[SelectNibble(c.opcode, 2)] |= c.V[SelectNibble(c.opcode, 1)]
}

func (c *CPU) _8XY2() {
	//Vx = Vx AND Vy
	c.V[SelectNibble(c.opcode, 2)] &= c.V[SelectNibble(c.opcode, 1)]
}

func (c *CPU) _8XY3() {
	//Vx = Vx XOR Vy
	c.V[SelectNibble(c.opcode, 2)] ^= c.V[SelectNibble(c.opcode, 1)]
}

func (c *CPU) _8XY4() {
	//Vx = Vx + Vy, set Vy = carry
	x := SelectNibble(c.opcode, 2)
	y := SelectNibble(c.opcode, 1)
	if c.V[x] > (0xFF - c.V[y]) {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] += c.V[y]
}

func (c *CPU) _8XY5() {
	//Vx = Vx - Vy, set VF = NOT borrow
	x := SelectNibble(c.opcode, 2)
	y := SelectNibble(c.opcode, 1)
	if c.V[x] > c.V[y] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] -= c.V[y]
}

func (c *CPU) _8XY6() {
	//Vx = Vx SHR 1
	x := SelectNibble(c.opcode, 2)
	c.V[0xF] = c.V[x] & 0x1
	c.V[x] >>= 1
}

func (c *CPU) _8XY7() {
	//Vx = Vy - Vx, set VF = NOT borrow.
	x := SelectNibble(c.opcode, 2)
	y := SelectNibble(c.opcode, 1)
	if c.V[y] > c.V[x] {
		c.V[0xF] = 1
	} else {
		c.V[0xF] = 0
	}
	c.V[x] = c.V[y] - c.V[x]
}

func (c *CPU) _8XYE() {
	//Vx = Vx SHL 1
	x := SelectNibble(c.opcode, 2)
	c.V[0xF] = (c.V[x] & 0x80) >> 7
	c.V[x] <<= 1
}

func (c *CPU) _9XY0() {
	//Skip next instruction if Vx != Vy.
	if c.V[SelectNibble(c.opcode, 2)] != c.V[SelectNibble(c.opcode, 1)] {
		c.PC += 2
	}
}

func (c *CPU) ANNN() {
	//I = nnn
	c.I = c.opcode & 0x0FFF
}

func (c *CPU) BNNN() {
	//Jump to location nnn + V0.
	c.PC = uint16(c.V[0]) + (c.opcode & 0x0FFF)
}

func (c *CPU) CXKK() {
	//Vx = random byte AND kk.
	c.V[SelectNibble(c.opcode, 2)] = byte(rand.Intn(256)) & byte(c.opcode&0x00FF)
}

func (c *CPU) DXYN() {
	//Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision
	x := uint16(c.V[SelectNibble(c.opcode, 2)])
	y := uint16(c.V[SelectNibble(c.opcode, 1)])
	height := SelectNibble(c.opcode, 0)

	c.V[0xF] = 0
	for _y := uint16(0); _y < height; _y++ {
		pixels := c.memory.memory[c.I+_y]
		for _x := uint16(0); _x < 8; _x++ {
			pixel := (pixels & (0x80 >> _x)) >> (7 - _x)
			x_pos := (x + _x) % 64
			y_pos := (y + _y) % 32
			if pixel == 1 && c.gfx[(y_pos)*64+x_pos] == 1 {
				c.V[0xF] = 1
			}
			c.gfx[(y_pos)*64+x_pos] ^= pixel
		}
	}
}

func (c *CPU) EX9E() {
	//Skip next instruction if key with the value of Vx is pressed
	if c.keys[c.V[SelectNibble(c.opcode, 2)]&0xF] {
		c.PC += 2
	}
}

func (c *CPU) EXA1() {
	//Skip next instruction if key with the value of Vx is not pressed
	if !c.keys[c.V[SelectNibble(c.opcode, 2)]&0xF] {
		c.PC += 2
	}
}

func (c *CPU) FX07() {
	//Vx = delay timer value
	c.V[SelectNibble(c.opcode, 2)] = c.timers.delayTimer
}

func (c *CPU) FX0A() {
	//Wait for a key press, store the value of the key in Vx
	for i, v := range c.keys {
		if v {
			c.V[SelectNibble(c.opcode, 2)] = byte(i)
			return
		}
	}
	//return to the same instruction
	c.opcode -= 2
}

func (c *CPU) FX15() {
	//delay timer = Vx
	c.timers.delayTimer = c.V[SelectNibble(c.opcode, 2)]
}

func (c *CPU) FX18() {
	//sound timer = Vx
	c.timers.soundTimer = c.V[SelectNibble(c.opcode, 2)]
}

func (c *CPU) FX1E() {
	//I = I + Vx
	c.I += uint16(c.V[SelectNibble(c.opcode, 2)])
}

func (c *CPU) FX29() {
	//I = location of sprite for digit Vx
	c.I = 5 * uint16(c.V[SelectNibble(c.opcode, 2)]&0xF)
}

func (c *CPU) FX33() {
	//Store BCD representation of Vx in memory locations I, I+1, and I+2
	regVal := c.V[SelectNibble(c.opcode, 2)]
	c.memory.memory[c.I] = regVal / 100
	c.memory.memory[c.I+1] = (regVal / 10) % 10
	c.memory.memory[c.I+2] = regVal % 10

}

func (c *CPU) FX55() {
	//Store registers V0 through Vx in memory starting at location I
	for i := uint16(0); i < SelectNibble(c.opcode, 2); i++ {
		c.memory.memory[c.I+i] = c.V[i]
	}
}

func (c *CPU) FX65() {
	//Read registers V0 through Vx from memory starting at location I
	for i := uint16(0); i < SelectNibble(c.opcode, 2); i++ {
		c.V[i] = c.memory.memory[c.I+i]
	}
}
