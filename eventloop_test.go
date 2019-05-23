package falcon

import (
	"fmt"
	"testing"
	"time"
)

func TestEventLoop(t *testing.T) {
	loop := NewEventLoop()
	loop.RunInLoop(func() {
		fmt.Println("close in loop")
		loop.Close()
	})
	loop.Loop()
}

func TestTimer(t *testing.T) {
	loop := NewEventLoop()
	loop.RunAfter(time.Second, func(t time.Time) {
		fmt.Println("on timer")
		timer := loop.RunAfter(time.Second, func(t time.Time) {
			fmt.Println("on timer stop")
		})
		timer.Stop()
		loop.RunAfter(time.Second*2, func(t time.Time) {
			fmt.Println("on timer close")
			loop.Close()
		})
	})
	loop.Loop()
}

func TestTicker(t *testing.T) {
	loop := NewEventLoop()
	times := 0
	var ticker *Ticker
	ticker = loop.RunEvery(time.Second, func(t time.Time) {
		times++
		fmt.Println("on ticker", times)
		if times == 3 {
			ticker.Stop()
			ticker = loop.RunEvery(time.Second, func(t time.Time) {
				times--
				fmt.Println("on ticker", times)
				if times == 0 {
					ticker.Stop()
					loop.Close()
				}
			})
		}
	})
	loop.Loop()
}
