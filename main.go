package main

import (
	"image/color"
	"time"

	"github.com/ajkula/golenia/game"
	"github.com/ajkula/golenia/graphics"
	"github.com/ajkula/golenia/input"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Lenia",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	graphics := graphics.NewGraphics(800, 600)
	lenia := game.NewLenia(800, 600)
	input := input.NewInput()

	go func() {
		for {
			lenia.Update()
			time.Sleep(time.Second / 30) // 30 FPS
		}
	}()

	for !window.Closed() {
		input.Update(window)
		graphics.Clear(color.RGBA{0, 0, 0, 255}) // noir
		lenia.Render(graphics)
		graphics.Render(window)
		window.Update()
		time.Sleep(time.Second / 60) // 60 FPS
	}
}

func main() {
	pixelgl.Run(run)
}
