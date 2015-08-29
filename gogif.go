package main

import (
	"errors"
	"fmt"
	"github.com/nsf/termbox-go"
	"image/color"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type AttribVals struct {
	fg    termbox.Attribute
	bg    termbox.Attribute
	trans bool
}

type AttribTable [256]AttribVals

type ColourMapFunc func(color.Color) AttribVals

func CMapMono(c color.Color) AttribVals {

	g := color.GrayModel.Convert(c).(color.Gray)

	return AttribVals{
		fg: termbox.Attribute(g.Y/11 + 1),
		bg: termbox.Attribute(g.Y/11 + 1),
	}
}

func CMapRGB(c color.Color) AttribVals {

	rgb := color.RGBAModel.Convert(c).(color.RGBA)

	r, g, b := int(rgb.R), int(rgb.G), int(rgb.B)

	i := uint8((r*6/256)*36 + (g*6/256)*6 + (b * 6 / 256))

	return AttribVals{
		fg:    termbox.Attribute(i + 17),
		bg:    termbox.Attribute(i + 17),
		trans: (rgb.A == 0),
	}
}

func mapColours(g *gif.GIF, cmap ColourMapFunc) []AttribTable {
	var attribs []AttribTable

	for f := 0; f < len(g.Image); f++ {
		var at AttribTable
		for i := 0; i < len(g.Image[f].Palette); i++ {
			at[i] = cmap(g.Image[f].Palette[i])
		}
		attribs = append(attribs, at)
	}

	return attribs
}

func renderFrame(g *gif.GIF, framenum int, attribs []AttribTable) error {

	width, height := termbox.Size()

	if width > g.Image[framenum].Rect.Dx() {
		width = g.Image[framenum].Rect.Dx()
	}

	if height > g.Image[framenum].Rect.Dy() {
		height = g.Image[framenum].Rect.Dy()
	}

	for y := 0; y < height; y++ {
		lineOffset := g.Image[framenum].Stride * y
		for x := 0; x < width; x++ {
			i := g.Image[framenum].Pix[x+lineOffset]
			attr := attribs[framenum][i]
			if !attr.trans {
				termbox.SetCell(x, y, ' ', attr.fg, attr.bg)
			}
		}
	}

	return nil
}

func openFile(name string) (io.ReadCloser, error) {
	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		resp, err := http.Get(name)
		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}

	return os.Open(name)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: gogif <filename>")
		os.Exit(1)
	}

	f, err := openFile(os.Args[1])
	exitOnError(err)

	g, err := gif.DecodeAll(f)
	exitOnError(err)

	_ = "breakpoint"

	gc := GameCore{}
	gc.TickTime = time.Millisecond * 50

	gc.OnInit = func(gc *GameCore) error {
		//		mode := termbox.SetOutputMode(termbox.OutputGrayscale)
		mode := termbox.SetOutputMode(termbox.Output256)

		if mode != termbox.OutputGrayscale {
			return errors.New("Failed to set output mode")
		}

		return nil
	}

	gc.OnEvent = func(gc *GameCore, ev *termbox.Event) error {
		if ev.Type == termbox.EventKey {
			if ev.Ch == 'q' {
				gc.DoQuit = true
			}
		}
		return nil
	}

	frameNumber := 0

	//	attribs := mapColours(g, CMapMono)
	attribs := mapColours(g, CMapRGB)

	gc.OnTick = func(gc *GameCore) error {
		err := renderFrame(g, frameNumber, attribs)
		frameNumber++
		if len(g.Image) == frameNumber {
			frameNumber = 0
			//gc.DoQuit = true
		}
		return err
	}

	gc.Run()
}
