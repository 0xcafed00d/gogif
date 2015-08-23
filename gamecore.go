package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

type OnInitFunc func(gc *GameCore) error
type OnTickFunc func(gc *GameCore) error
type OnEventFunc func(gc *GameCore, ev *termbox.Event) error

type GameCore struct {
	OnInit   OnInitFunc
	OnTick   OnTickFunc
	OnEvent  OnEventFunc
	DoQuit   bool
	TickTime time.Duration
	Ticker   *time.Ticker
}

func (gc *GameCore) Run() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	if gc.OnInit != nil {
		gc.OnInit(gc)
		if err != nil {
			return err
		}
	}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	gc.Ticker = time.NewTicker(gc.TickTime)

	for !gc.DoQuit {
		select {
		case ev := <-eventQueue:
			if gc.OnEvent != nil {
				gc.OnEvent(gc, &ev)
			}
			if ev.Type == termbox.EventResize {
				termbox.Flush()
			}

		case <-gc.Ticker.C:
			if gc.OnTick != nil {
				gc.OnTick(gc)
				termbox.Flush()
			}
		}
	}
	return nil
}
