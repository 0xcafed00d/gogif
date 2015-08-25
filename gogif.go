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

type ColourMapFunc func(idx uint8) (termbox.Attribute, termbox.Attribute)

func test(idx uint8) (termbox.Attribute, termbox.Attribute) {
	return termbox.Attribute(idx), termbox.Attribute(idx)
}

func renderFrame(g *gif.GIF, framenum int, cmap ColourMapFunc) error {

	width, height := termbox.Size()

	for y := 0; y < height; y++ {
		lineOffset := g.Image[framenum].Stride * y
		for x := 0; x < width; x++ {
			fg, bg := cmap(g.Image[framenum].Pix[x+lineOffset])
			termbox.SetCell(x, y, ' ', fg, bg)
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

	gc.OnTick = func(gc *GameCore) error {
		err := renderFrame(g, frameNumber, test)
		frameNumber++
		if len(g.Image) == frameNumber {
			frameNumber = 0
		}
		return err
	}

	gc.Run()
}
