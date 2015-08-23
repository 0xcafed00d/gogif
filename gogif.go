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

func renderFrame(g *gif.GIF, framenum int) error {

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
		err := renderFrame(g, frameNumber)
		frameNumber++
		return err
	}

	gc.Run()
}
