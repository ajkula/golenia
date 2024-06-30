package main

import (
	"time"

	"github.com/ajkula/golenia/game"
	"github.com/ajkula/golenia/graphics"
	"github.com/ajkula/golenia/input"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	gridWidth  = 200
	gridHeight = 150
	scale      = 4
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Lenia",
		Bounds: pixel.R(0, 0, float64(gridWidth*scale), float64(gridHeight*scale)),
		VSync:  true,
	}

	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	graphics := graphics.NewGraphics(gridWidth*scale, gridHeight*scale)
	lenia := game.NewLenia(gridWidth, gridHeight)
	input := input.NewInput()

	timeStep := 1.0 / 60.0 // 30 updates per second
	lastUpdate := time.Now()

	for !window.Closed() {
		input.Update(window)

		now := time.Now()
		elapsed := now.Sub(lastUpdate).Seconds()
		if elapsed >= timeStep {
			lenia.Update()
			lastUpdate = now
		}

		if window.JustPressed(pixelgl.KeySpace) {
			lenia.Reset()
		}

		if window.JustPressed(pixelgl.KeyUp) {
			lenia.IncreaseSpeed()
		}

		if window.JustPressed(pixelgl.KeyDown) {
			lenia.DecreaseSpeed()
		}

		graphics.Clear(pixel.RGB(0, 0, 0))
		lenia.Render(graphics)
		graphics.Render(window)
		window.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
