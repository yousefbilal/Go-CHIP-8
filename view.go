package main

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/jroimartin/gocui"
)

const (
	red   = "\x1b[31m"
	reset = "\x1b[0m|"
)

func colorize(str string, color string) string {
	return fmt.Sprintf("%v%v%v", color, str, reset)
}

func generateKeypadLayout(keys [16]bool) string {
	keySymbols := [16]string{"1", "2", "3", "C", "4", "5", "6", "D", "7", "8", "9", "E", "A", "0", "B", "F"}
	layout := ""
	for i, key := range keySymbols {
		if i%4 == 0 {
			layout += "\n+-+-+-+-+\n|"
		}
		keyIndex, _ := strconv.ParseInt(key, 16, 0)
		if keys[keyIndex] {
			layout += colorize(keySymbols[i], red)
		} else {
			layout += fmt.Sprintf("%s|", keySymbols[i])
		}
	}
	return layout
}

func layout(g *gocui.Gui, chip8 *CPU) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("keypad", 1, 1, 11, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keypad"
		keypadLayout := generateKeypadLayout(chip8.keys)
		fmt.Fprintln(v, keypadLayout)
	}

	if v, err := g.SetView("registers", 12, 1, 28, 18); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Registers-V"
		for i, val := range chip8.V {
			fmt.Fprintf(v, "0x%x : %v 0x%x\n", i, val, val)
		}
	}

	if v, err := g.SetView("memory", 29, 1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Memory"
		v.Wrap = true
		fmt.Fprint(v, hex.EncodeToString(chip8.memory.memory[:]))
	}

	if v, err := g.SetView("misc", 1, 11, 11, 18); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Misc"
		fmt.Fprintf(v, "PC: %x\nI: %x\nop: %x\nSP: %x\nDT: %x\nST: %x",
			chip8.PC, chip8.I, chip8.opcode, chip8.SP, chip8.timers.delayTimer, chip8.timers.soundTimer)
	}
	return nil
}

func updateLayout(g *gocui.Gui, chip8 *CPU) error {
	keypadView, err := g.View("keypad")
	if err != nil {
		return err
	}
	// Clear the view and re-print the updated values
	keypadView.Clear()
	keypadLayout := generateKeypadLayout(chip8.keys)
	fmt.Fprintln(keypadView, keypadLayout)

	registersView, err := g.View("registers")
	if err != nil {
		return err
	}
	registersView.Clear()
	for i, val := range chip8.V {
		fmt.Fprintf(registersView, "0x%x : %v 0x%x\n", i, val, val)
	}

	memoryView, err := g.View("memory")
	if err != nil {
		return err
	}
	memoryView.Clear()
	fmt.Fprint(memoryView, hex.EncodeToString(chip8.memory.memory[:]))

	miscView, err := g.View("misc")
	if err != nil {
		return err
	}
	miscView.Clear()
	fmt.Fprintf(miscView, "PC: %x\nI: %x\nop: %x\nSP: %x\nDT: %x\nST: %x",
		chip8.PC, chip8.I, chip8.opcode, chip8.SP, chip8.timers.delayTimer, chip8.timers.soundTimer)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
