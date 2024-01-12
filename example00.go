/*  filename: example00.go  */
package main

import "fmt"

func main() {

	for _, c64col := range c64palette {
		c := printColorExample(c64col, "  ")
		fmt.Printf("colour: %02d|%s|:%s\n", c64col.id, c, c64col.name)
	}
}

func printColorExample(c c64color, s string) string {
	p := fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[m", c.r, c.g, c.b, s)
	return p
}

var c64palette = []c64color{
	{0x00, 0x00, 0x00, 0, "black", 0},
	{0xFF, 0xFF, 0xFF, 1, "white", 0},
	{0x68, 0x37, 0x2B, 2, "red", 0},
	{0x70, 0xA4, 0xB2, 3, "cyan", 0},
	{0x6F, 0x3D, 0x86, 4, "purple", 0},
	{0x58, 0x8D, 0x43, 5, "green", 0},
	{0x35, 0x28, 0x79, 6, "navy", 0},
	{0xB8, 0xC7, 0x6F, 7, "yellow", 0},
	{0x6F, 0x4F, 0x25, 8, "orange", 0},
	{0x43, 0x39, 0x00, 9, "brown", 0},
	{0x9A, 0x67, 0x59, 10, "light-red", 0},
	{0x44, 0x44, 0x44, 11, "dark-grey", 0},
	{0x6C, 0x6C, 0x6C, 12, "grey", 0},
	{0x9A, 0xD2, 0x84, 13, "light-green", 0},
	{0x6C, 0x5E, 0xB5, 14, "blue", 0},
	{0x95, 0x95, 0x95, 15, "light-grey", 0},
}

type c64color struct {
	r            uint8
	g            uint8
	b            uint8
	id           uint8
	name         string
	lumaDistance float64
}
