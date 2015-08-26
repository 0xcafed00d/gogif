package main

import (
	"fmt"
	//"image"
	"github.com/nsf/termbox-go"
	"image/gif"
	"os"
	"time"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type AttribVals struct {
	fg termbox.Attribute
	bg termbox.Attribute
}

type AttribTable [256]AttribVals

type ColourMapFunc func(idx uint8) AttribVals

func test(idx uint8) AttribVals {
	return AttribVals{fg: termbox.Attribute(idx), bg: termbox.Attribute(idx)}
}

func mapColours(g *gif.GIF, cmap ColourMapFunc) []AttribTable {
	var attribs []AttribTable

	for f := 0; f < len(g.Image); f++ {

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
			attr := attribs[framenum][x+lineOffset]
			termbox.SetCell(x, y, ' ', attr.fg, attr.bg)
		}
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: gogif <filename>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	exitOnError(err)

	g, err := gif.DecodeAll(f)
	exitOnError(err)

	gc := GameCore{}
	gc.TickTime = time.Millisecond * 40

	gc.OnEvent = func(gc *GameCore, ev *termbox.Event) error {
		if ev.Type == termbox.EventKey {
			if ev.Ch == 'q' {
				gc.DoQuit = true
			}
		}
		return nil
	}

	frameNumber := 0

	attribs := mapColours(g, test)

	gc.OnTick = func(gc *GameCore) error {
		err := renderFrame(g, frameNumber, attribs)
		frameNumber++
		if len(g.Image) == frameNumber {
			frameNumber = 0
		}
		return err
	}

	gc.Run()
}
