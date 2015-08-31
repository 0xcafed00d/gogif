package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Mono bool
	Once bool
	Help bool
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "Display Help")
	flag.BoolVar(&config.Mono, "m", false, "Play in Monochrome mode")
	flag.BoolVar(&config.Once, "o", false, "Play animation once")
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

func showHelp() {
	fmt.Fprintln(os.Stderr, "Usage: gogif [options] <filename>")
	flag.PrintDefaults()
}

func main() {

	flag.Usage = showHelp
	flag.Parse()

	if len(flag.Args()) == 0 || config.Help {
		showHelp()
		os.Exit(1)
	}

	f, err := openFile(flag.Args()[0])
	exitOnError(err)

	g, err := gif.DecodeAll(f)
	exitOnError(err)

	gc := GameCore{}
	gc.TickTime = time.Millisecond * 50

	gc.OnInit = func(gc *GameCore) error {
		//mode := termbox.SetOutputMode(termbox.OutputGrayscale)
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

	//attribs := mapColours(g, CMapMono)
	attribs := mapColours(g, CMapRGB)

	gc.OnTick = func(gc *GameCore) error {
		err := renderFrameHiRes(g, frameNumber, attribs)
		frameNumber++
		if len(g.Image) == frameNumber {
			frameNumber = 0
			//gc.DoQuit = true
		}
		return err
	}

	gc.Run()
}
