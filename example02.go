/*  example02.go */
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
)

var indexC1 int = 6
var indexC2 int = 14

func main() {
	pal50 := generateMixedPalette50percent(c64palette)
	separator()
	pal25 := generateMixedPalette25percent(c64palette)
	separator()
	img := OpenImage("./test07.png")
	renderImgC64(img, c64palette)
	separator()
	renderImgC64(img, pal50)
	separator()
	renderImgC64(img, pal25)
	separator()

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
func generateMixedPalette25percent(pal []c64color) []c64color {
	var newPalette []c64color
	for _, bG := range pal {
		for _, fG := range pal {
			result := c64color{
				r:    uint8((float64(bG.r)*75 + 25*float64(fG.r)) / 100),
				g:    uint8((float64(bG.g)*75 + 25*float64(fG.g)) / 100),
				b:    uint8((float64(bG.b)*75 + 25*float64(fG.b)) / 100),
				id:   fG.id<<4 | bG.id,
				name: alignString(fG.name+"+"+bG.name, 16),
			}
			color := PrintExampleRGB(c64palette[fG.id], c64palette[bG.id], "░░")
			fmt.Print(color)
			newPalette = append(newPalette, result)
		}
		fmt.Println()
	}
	return newPalette
}
func generateMixedPalette50percent(pal []c64color) []c64color {
	var newPalette []c64color
	for _, bG := range pal {
		for _, fG := range pal {
			result := c64color{

				r:    uint8((float64(bG.r)*50 + 50*float64(fG.r)) / 100),
				g:    uint8((float64(bG.g)*50 + 50*float64(fG.g)) / 100),
				b:    uint8((float64(bG.b)*50 + 50*float64(fG.b)) / 100),
				id:   fG.id<<4 | bG.id,
				name: alignString(fG.name+"+"+bG.name, 16),
			}
			color := PrintExampleRGB(c64palette[fG.id], c64palette[bG.id], "▒▒")
			fmt.Print(color)
			newPalette = append(newPalette, result)
		}
		fmt.Println()
	}
	return newPalette
}
func alignString(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}
func check4errors(s string, err error) {
	if err != nil {
		log.Println(s, err)
	}
}
func OpenImage(s string) image.Image {
	freader, err := os.Open(s)
	check4errors("opening image file", err)
	defer freader.Close()
	img, _, err := image.Decode(freader)
	check4errors("decoding img", err)
	return img
}
func renderImgC64(i image.Image, p []c64color) *[]uint8 {
	xSize := i.Bounds().Dx()
	ySize := i.Bounds().Dy()
	var frame []uint8
	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			rr, gg, bb, _ := i.At(x, y).RGBA()
			r := uint8(rr >> 8)
			g := uint8(gg >> 8)
			b := uint8(bb >> 8)

			ou := getNearestNeighborRGB(r, g, b, p)
			PrintRGB(0, 0, 0, ou.r, ou.g, ou.b, "  ")
			frame = append(frame, ou.id)
		}
		os.Stderr.WriteString("\x1b[m\n")
	}
	return &frame
}

func PrintExampleRGB(c, b c64color, s string) string {
	p := fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1b[m", c.r, c.g, c.b, b.r, b.g, b.b, s)
	return p
}
func getNearestNeighborRGB(r, g, b uint8, cl []c64color) c64color {
	baseR := float64(r)
	baseG := float64(g)
	baseB := float64(b)
	var deltaError float64
	deltaError = 500
	idFgBg := uint8(0)
	for _, v := range cl {
		commR := float64(v.r)
		commG := float64(v.g)
		commB := float64(v.b)

		curretError := math.Sqrt(math.Pow((baseR-commR), 2) + math.Pow((baseG-commG), 2) + math.Pow((baseB-commB), 2))
		if curretError < deltaError {
			deltaError = curretError
			r, g, b = uint8(commR), uint8(commG), uint8(commB)
			idFgBg = v.id
		}
	}
	return c64color{r, g, b, idFgBg, "", 0}
}
func PrintRGB(rf, gf, bf, rb, gb, bb uint8, s string) {
	p := fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s", rf, gf, bf, rb, gb, bb, s)
	os.Stdout.WriteString(p)
}
func separator() {
	fmt.Println()
}
