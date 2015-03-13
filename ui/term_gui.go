package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/termbox-go"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var palette = xtermPalette()

var fontPath = "test_data/Hack-Regular.ttf"
var fontSize = 12

type GuiTerm struct {
	w, h         int
	text         [][]char
	textLock     sync.Mutex
	win          wde.Window
	font         *truetype.Font
	charW, charH int // size of characters
	face         font.Face
	ctx          *freetype.Context
	rgba         *image.RGBA
}

type char struct {
	rune
	fg, bg core.Style
}

func NewGuiTerm(h, w int) *GuiTerm {
	win, err := wde.NewWindow(1400, 800)
	if err != nil {
		panic(err)
	}

	t := &GuiTerm{
		win: win,
	}

	t.applyFont(fontPath, fontSize)

	t.text = [][]char{}

	for i := 0; i != t.h; i++ {
		t.text = append(t.text, make([]char, t.w))
	}

	return t
}

func (t *GuiTerm) applyFont(fontPath string, fontSize int) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	t.font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	opts := truetype.Options{}
	opts.Size = float64(fontSize)
	t.face = truetype.NewFace(t.font, &opts)
	bounds, _, _ := t.face.GlyphBounds('░')
	t.charW = int((bounds.Max.X-bounds.Min.X)>>6) + 2
	t.charH = int((bounds.Max.Y-bounds.Min.Y)>>6) + 2
	ww, wh := t.win.Size()
	t.w = ww / t.charW
	t.h = wh / t.charH

	t.rgba = image.NewRGBA(image.Rect(0, 0, ww, wh))

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(t.font)
	c.SetFontSize(float64(fontSize))
	c.SetClip(t.rgba.Bounds())
	c.SetDst(t.rgba)
	c.SetHinting(font.HintingFull)
	t.ctx = c

	fmt.Printf("w:%d h:%d, cw:%d, ch:%d\n", t.w, t.h, t.charW, t.charH)
}

func (t *GuiTerm) Init() error {
	t.win.Show()
	return nil
}

func (t *GuiTerm) Close() {
	t.win.Close()
}

func (t *GuiTerm) Clear(fg, bg uint16) {
	c := image.NewUniform(palette[bg&255])
	x, y := t.win.Size()
	draw.Draw(t.win.Screen(), image.Rect(0, 0, x, y), c, image.ZP, draw.Src)
}

func (t *GuiTerm) Flush() {
	t.paint()
}

func (t *GuiTerm) SetCursor(y, x int) {
	// todo : move cursor
}

func (t *GuiTerm) Char(y, x int, c rune, fg, bg core.Style) {
	t.textLock.Lock()
	defer t.textLock.Unlock()
	if x >= 0 && y >= 0 && y < len(t.text) && x < len(t.text[y]) {
		t.text[y][x] = char{
			rune: c,
			fg:   fg,
			bg:   bg,
		}
	}
}

// size in characters
func (t *GuiTerm) Size() (h, w int) {
	return t.h, t.w
}

// for testing
func (t *GuiTerm) CharAt(y, x int) rune {
	t.textLock.Lock()
	defer t.textLock.Unlock()
	if x < 0 || y < 0 {
		panic("CharAt out of bounds")
	}
	if y >= t.h || x >= t.w {
		panic("CharAt out of bounds")
	}
	return t.text[y][x].rune
}

func (t *GuiTerm) SetMouseMode(m termbox.MouseMode) { // N/A
}

func (t *GuiTerm) SetInputMode(m termbox.InputMode) { // N/A
}

func (t *GuiTerm) SetExtendedColors(b bool) { // N/A
}

func (t *GuiTerm) paint() {
	c := t.ctx
	w := fixed.Int26_6(t.charW << 6)
	h := fixed.Int26_6(t.charH << 6)
	pt := freetype.Pt(1, t.charH-4)
	for y, ln := range t.text {
		for x, r := range ln {
			if r.rune == 0 {
				r.rune = ' '
				r.bg = core.Ed.Theme().Bg
			}
			// TODO: attributes (bold)
			bg := image.NewUniform(palette[r.bg.Uint16()&255])
			fg := image.NewUniform(palette[r.fg.Uint16()&255])
			c.SetSrc(fg)
			//bounds, awidth, _ := t.face.GlyphBounds(r.rune)
			//fmt.Printf("%s %v | %v\n",
			//	string(r.rune),
			//	bounds,
			//	awidth)
			rx := t.charW * x
			ry := t.charH * y
			rect := image.Rect(rx, ry, rx+t.charW, ry+t.charH)
			draw.Draw(t.rgba, rect, bg, image.ZP, draw.Src)
			c.DrawString(string(r.rune), pt)
			pt.X += w
		}
		pt.X = 1
		pt.Y += h
	}
	t.win.Screen().CopyRGBA(t.rgba, t.rgba.Bounds())
	t.win.FlushImage()
}

// Palette based of what's used in gnome-terminal / xterm-256
func xtermPalette() *[256]color.Color {
	a := uint8(255)
	// base colors (from gnome-terminal)
	palette := [256]color.Color{
		color.RGBA{0x2e, 0x34, 0x36, a},
		color.RGBA{0xcc, 0, 0, a},
		color.RGBA{0x4e, 0x9a, 0x06, a},
		color.RGBA{0xc4, 0xa0, 0, a},
		color.RGBA{0x34, 0x65, 0xa4, a},
		color.RGBA{0x75, 0x50, 0x7b, a},
		color.RGBA{0x06, 0x98, 0x9a, a},
		color.RGBA{0xd3, 0xd7, 0xcf, a},
		color.RGBA{0x55, 0x57, 0x53, a},
		color.RGBA{0xef, 0x29, 0x29, a},
		color.RGBA{0x8a, 0xe2, 0x34, a},
		color.RGBA{0xfc, 0xe9, 0x4f, a},
		color.RGBA{0x72, 0x9f, 0xcf, a},
		color.RGBA{0xad, 0x7f, 0xa8, a},
		color.RGBA{0x34, 0xe2, 0xe2, a},
		color.RGBA{0xee, 0xee, 0xec, a},
	}
	// xterm-256 colors
	for i := 16; i != 232; i++ {
		b := ((i - 16) % 6) * 40
		if b != 0 {
			b += 55
		}
		g := (((i - 16) / 6) % 6) * 40
		if g != 0 {
			g += 55
		}
		r := ((i - 16) / 36) * 40
		if r != 0 {
			r += 55
		}
		palette[i] = color.RGBA{uint8(r), uint8(g), uint8(b), a}
	}
	// Shades of grey
	for i := 232; i != 256; i++ {
		h := 8 + (i-232)*10
		palette[i] = color.RGBA{uint8(h), uint8(h), uint8(h), a}
	}

	return &palette
}