//go:build tinygo

package koebiten

import (
	"fmt"
	"image/color"
	"io/fs"
	"machine"
	"strings"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/image/png"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
)

type Displayer interface {
	drivers.Displayer
	ClearDisplay()
	ClearBuffer()
}

var (
	btn     machine.Pin
	display Displayer
)

var (
	textY int16
)

var (
	white = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	black = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

func init() {
	pngBuffer = map[string]Image{}
}

func Run(d func()) error {
	tick := time.Tick(32 * time.Millisecond)
	for {
		<-tick
		keyUpdate()
		textY = 0
		display.ClearBuffer()
		d()
		display.Display()
	}
	return nil
}

func SetWindowSize(w, h int) {
}

func IsClicked() bool {
	return !btn.Get()
}

func Println(args ...any) {

	str := []string{}
	for _, x := range args {
		s, ok := x.(string)
		if ok {
			str = append(str, s)
			continue
		}

		i, ok := x.(int)
		if ok {
			str = append(str, fmt.Sprintf("%d", i))
			continue
		}
	}

	textY += 8
	tinyfont.WriteLine(display, &tinyfont.Org01, 2, textY, strings.Join(str, " "), white)
}

func DrawRect(x, y, w, h int) {
	tinydraw.Rectangle(display, int16(x), int16(y), int16(w), int16(h), white)
}

func DrawLine(x1, y1, x2, y2 int) {
	tinydraw.Line(display, int16(x1), int16(y1), int16(x2), int16(y2), white)
}

func DrawCircle(x, y, r int) {
	tinydraw.Circle(display, int16(x), int16(y), int16(r), white)
}

type Image struct {
	W   int16
	H   int16
	Buf []bool
}

var (
	buffer    [3 * 8 * 8 * 4]uint16
	pngBuffer map[string]Image
)

func DrawImageFS(fsys fs.FS, path string, x, y int) {
	img, ok := pngBuffer[path]
	if !ok {
		p, err := fsys.Open(path)
		if err != nil {
			return
		}

		png.SetCallback(buffer[:], func(data []uint16, x, y, w, h, width, height int16) {
			img.W = width
			img.H = height

			if img.Buf == nil || len(img.Buf) == 0 {
				img.Buf = make([]bool, width*height)
			}

			for yy := int16(0); yy < h; yy++ {
				for xx := int16(0); xx < w; xx++ {
					c := C565toRGBA(data[yy*w+xx])
					cnt := 0
					if c.R < 0x80 {
						cnt++
					}
					if c.G < 0x80 {
						cnt++
					}
					if c.B < 0x80 {
						cnt++
					}
					if cnt >= 2 {
						img.Buf[y*width+x+yy*w+xx] = true
					}
				}
			}
		})

		_, err = png.Decode(p)
		if err != nil {
			return
		}
		pngBuffer[path] = img
	}

	for yy := int16(0); yy < img.H; yy++ {
		for xx := int16(0); xx < img.W; xx++ {
			if img.Buf[yy*img.W+xx] {
				display.SetPixel(int16(x)+xx, int16(y)+yy, white)
			}
		}
	}
}
