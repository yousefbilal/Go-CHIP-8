package main

func SelectNibble(val uint16, index uint16) uint16 {
	return val & (0x000F << (4 * index)) >> (4 * index)
}
