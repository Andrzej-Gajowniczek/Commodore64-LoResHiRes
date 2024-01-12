/*  example03.go  */
package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	"log"
	"math"
	"os"
)

var maxDistance float64 = 140
var imgFile string = "test23a.png"

var indexC1 int = 6
var indexC2 int = 14

//var imgFile string = "modified.gif"

func main() {

	var color []c64color
	color = append(color, c64palette[indexC1])
	color = append(color, c64palette[indexC2])
	for _, v := range color {
		example := PrintExample(v, "  ")
		fmt.Printf("R:%d G:%d B:%d \t c64 index:%d\tcolor name:|%s|:%s\n", v.r, v.g, v.b, v.id, example, v.name)
	}
	fmt.Printf("the distance between |%s| and |%s| is %f\n", PrintExample(color[0], "  "), PrintExample(color[01], "  "), calculateLumaDistance((color[01]), (color[0])))
	fmt.Println(showShadesOf(color[0], color[1]))
	//os.Exit(0)
	/*newOrderPallette := */
	sortedPal := showSortedLumasInPalette(c64palette[indexC1], c64palette)
	//fmt.Println(sortedPal)
	generateMixedPalette25percent(c64palette)
	fmt.Println("")
	trashedLumaPalette := generateMixedPalette25percent(sortedPal)
	goodLumaPalette := dropBadLuma(trashedLumaPalette)
	fmt.Println("palettte size:", len(goodLumaPalette))
	img := OpenImage(imgFile)
	fmt.Println("25%")
	h, err := os.Create("c64.bin")
	check4errors("opening new file", err)
	defer h.Close()
	frame := renderImgC64(img, goodLumaPalette)
	h.Write(*frame)
	/*colorPairsMap :=  showStatistics(*frame)*/
	//	fmt.Println(colorPairsMap)

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
func getNearestNeighborRGB(r, g, b uint8, cl []c64color) c64color {
	//var i int
	baseR := float64(r)
	baseG := float64(g)
	baseB := float64(b)
	var deltaError float64
	deltaError = 500
	//	fmt.Println(deltaError)
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

		//		fmt.Println(deltaError)
	}
	//	fmt.Println(deltaError)
	return c64color{r, g, b, idFgBg, "", 0}
}
func PrintRGB(rf, gf, bf, rb, gb, bb uint8, s string) {
	p := fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s", rf, gf, bf, rb, gb, bb, s)
	os.Stdout.WriteString(p)
}

type c64color struct {
	r            uint8
	g            uint8
	b            uint8
	id           uint8
	name         string
	lumaDistance float64
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
	{0x9A, 0x67, 0x59, 10, "pink", 0},
	{0x44, 0x44, 0x44, 11, "dark-grey", 0},
	{0x6C, 0x6C, 0x6C, 12, "grey", 0},
	{0x9A, 0xD2, 0x84, 13, "lightgreen", 0},
	{0x6C, 0x5E, 0xB5, 14, "blue", 0},
	{0x95, 0x95, 0x95, 15, "lightgrey", 0},
}

func PrintExample(c c64color, s string) string {
	p := fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s\x1b[m", c.r, c.g, c.b, s)
	//os.Stdout.WriteString(p)
	return p
}
func PrintExampleRGB(c, b c64color, s string) string {
	p := fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1b[m", c.r, c.g, c.b, b.r, b.g, b.b, s)
	//os.Stdout.WriteString(p)
	return p
}

func calculateLumaDistance(a, b c64color) float64 {
	y := make(map[c64color]float64)

	y[a] = 0.299*float64(a.r) + 0.587*float64(a.g) + 0.114*float64(a.b)
	y[b] = 0.299*float64(b.r) + 0.587*float64(b.g) + 0.114*float64(b.b)
	return math.Abs(y[a] - y[b])
}

func calculateRgbError(src, dest c64color) float64 {
	sr := float64(src.r)
	sg := float64(src.g)
	sb := float64(src.b)
	dr := float64(dest.r)
	dg := float64(dest.g)
	db := float64(dest.b)
	return math.Sqrt(math.Pow((sr-dr), 2) + math.Pow((sg-dg), 2) + math.Pow((sb-db), 2))

}

func showShadesOf(a, b c64color) string {
	var p string
	for _, ch := range []string{"██", "▓▓", "▒▒", "░░", "  "} {

		p = p + fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s\x1b[m", a.r, a.g, a.b, b.r, b.g, b.b, ch)
	}
	return p
}

func showSortedLumasInPalette(a c64color, b []c64color) []c64color {

	var resultLumaPalette []c64color
	resultLumaPalette = append(resultLumaPalette, a)
	PrintExample(a, "to jest ten kolor a")
	best := 255.01
	var fineColor c64color
	c := make(map[c64color]int)
	for v, k := range b {
		c[k] = v
	}
	fmt.Println("len of palette:", len(c))
	delete(c, a)
	fmt.Println("len of palette:", len(c))
	fmt.Print(PrintExample(a, "  "))
	for {

		for v, _ := range c {
			given := calculateLumaDistance(a, v)
			if given < best {
				best = given
				fineColor = v
			}
		}
		fmt.Print(PrintExample(fineColor, "  "))
		a = fineColor
		delete(c, a)
		resultLumaPalette = append(resultLumaPalette, a)
		best = 300
		if len(c) == 0 {
			break
		}

		//fmt.Println("len of palette:", len(c))
	}
	fmt.Println("\n")
	return resultLumaPalette
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
			//fmt.Print(PrintExample(result, "  "))
			newPalette = append(newPalette, result)
		}
		fmt.Println()
	}
	return newPalette
}
func dropBadLuma(pal []c64color) []c64color {
	var newPalette []c64color
	for i, v := range pal {
		b := v
		cFgRGB := c64palette[v.id>>4]
		cBgRGB := c64palette[b.id&0x0F]
		result := calculateLumaDistance(cFgRGB, cBgRGB)
		v.lumaDistance = result
		if result <= maxDistance {
			newPalette = append(newPalette, v)
			s := PrintExampleRGB(c64palette[v.id>>4], c64palette[b.id&0x0F], "░░")
			fmt.Printf("%s\t%02x\t%02x\t%02x %s %f\n", v.name, v.id>>4, b.id&0x0F, i, s, result)
		}
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
func showStatistics(frame []uint8) map[uint8]int {
	stats := make(map[uint8]int)
	out := make(map[uint8]int)
	var i int
	var v uint8
	for i, v = range frame {
		stats[v] += 1
		out[v] += 1
	}
	fmt.Printf("blocks: %d, total colors number is:%d\n", i, len(stats))
	var cPair uint8
	fmt.Println(len(stats))
	var y int
	for {
		if len(stats) == 0 {
			break
		}
		var value int

		for colPair, count := range stats {
			if count > value {
				value = count
				cPair = colPair
			}
		}
		p := PrintExampleRGB(c64palette[cPair>>4], c64palette[cPair&0x0f], "░░")
		fmt.Printf("%03d colors :%x, %d\t%s\n", y, cPair, value, p)
		delete(stats, cPair)
		y++
	}
	return out
}
func getPairsLumaSorted(m map[uint8]int) []uint8 {
	var slice []uint8
	distance := 254.1
	testRgbColor := c64color{
		r: 0xff,
		g: 0xff,
		b: 0xff,
	}
	var r, k uint8
	var kcolor c64color
	for {
		if len(m) == 0 {
			break
		}
		for k, _ = range m {
			kcolor = c64color{
				r: uint8(float64(c64palette[k>>4].r)*.25 + float64(c64palette[k&0x0f].r)*.75),
				g: uint8(float64(c64palette[k>>4].g)*.25 + float64(c64palette[k&0x0f].g)*.75),
				b: uint8(float64(c64palette[k>>4].b)*.25 + float64(c64palette[k&0x0f].b)*.75),
			}
			y := calculateLumaDistance(testRgbColor, kcolor)
			if y < distance {
				distance = y
				r = k
			}

		}
		fmt.Println(r, len(m), distance)
		slice = append(slice, r)
		testRgbColor = kcolor
		distance = 255
		delete(m, r)
	}
	return slice
}

func makePalettFromPairs(pairs []uint8) []c64color {
	var newPalette []c64color
	for _, v := range pairs {
		newPalette = append(newPalette, pair2rgb(v))
	}
	return newPalette
}

func pair2rgb(v uint8) c64color {
	f := v >> 4
	g := v & 0x0f
	fRgb := c64palette[f]
	bRgb := c64palette[g]
	result := c64color{

		r:    uint8(float64(fRgb.r)*.25 + .75*float64(bRgb.r)),
		g:    uint8(float64(fRgb.g)*.25 + .75*float64(bRgb.g)),
		b:    uint8(float64(fRgb.b)*.25 + .75*float64(bRgb.b)),
		id:   v,
		name: c64palette[f].name + c64palette[g].name,
	}
	return result
}
