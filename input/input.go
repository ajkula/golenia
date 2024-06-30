package input

import (
	"github.com/faiface/pixel/pixelgl"
)

type Input struct {
	keys map[pixelgl.Button]bool
}

func NewInput() *Input {
	return &Input{
		keys: make(map[pixelgl.Button]bool),
	}
}

func (i *Input) Update(window *pixelgl.Window) {
	for key := range i.keys {
		i.keys[key] = window.Pressed(key)
	}
}

func (i *Input) IsKeyPressed(key pixelgl.Button) bool {
	return i.keys[key]
}