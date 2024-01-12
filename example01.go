/*  example01.go */
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

var indexC1 int = 6
var indexC2 int = 14

func main() {
	args := os.Args
	if len(args[1:]) == 2 {
		indexC1, _ = strconv.Atoi(args[1])
		indexC2, _ = strconv.Atoi(args[2])
	}

	fmt.Println("color 1:", indexC1, "color 2:", indexC2)
	var color []c64color
	color = append(color, c64palette[indexC1])
	color = append(color, c64palette[indexC2])
	for _, v := range color {
		example := PrintExample(v, "  ")
		fmt.Printf("R:%d G:%d B:%d \t c64 index:%d\tcolor name:|%s|:%s\n", v.r, v.g, v.b, v.id, example, v.name)
	}
	fmt.Printf("the distance between |%s| and |%s| is %f\n", PrintExample(color[0], "  "), PrintExample(color[01], "  "), calculateLumaDistance((color[01]), (color[0])))
	fmt.Println(showShadesOf(color[0], color[1]))

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
	{0x9A, 0xD2, 0x84, 13, "lightgreen", 0},
	{0x6C, 0x5E, 0xB5, 14, "blue", 0},
	{0x95, 0x95, 0x95, 15, "lightgrey", 0},
}

type c64color struct {
	r            uint8
	g            uint8
	b            uint8
	id           uint8
	name         string
	lumaDistance float64
}

func PrintExample(c c64color, s string) string {
	p := fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[m", c.r, c.g, c.b, s)
	return p
}

func calculateLumaDistance(a, b c64color) float64 {
	y := make(map[c64color]float64)

	y[a] = 0.299*float64(a.r) + 0.587*float64(a.g) + 0.114*float64(a.b)
	y[b] = 0.299*float64(b.r) + 0.587*float64(b.g) + 0.114*float64(b.b)
	return math.Abs(y[a] - y[b])
}
func showShadesOf(a, b c64color) string {
	var p string
	for _, ch := range []string{"██", "▓▓", "▒▒", "░░", "  "} {

		p = p + fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1b[m", a.r, a.g, a.b, b.r, b.g, b.b, ch)
	}
	return p
}
