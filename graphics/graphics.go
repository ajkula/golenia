package graphics

import (
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Graphics struct {
	screen *image.RGBA
	mu     sync.Mutex
}

func NewGraphics(width, height int) *Graphics {
	return &Graphics{
		screen: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (g *Graphics) Clear(c color.Color) {
	g.mu.Lock()
	defer g.mu.Unlock()
	draw.Draw(g.screen, g.screen.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

func (g *Graphics) DrawCell(x, y int, c color.Color) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.screen.Set(x, y, c)
}

func (g *Graphics) GetScreen() *image.RGBA {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.screen
}

func (g *Graphics) Render(window *pixelgl.Window) {
	g.mu.Lock()
	defer g.mu.Unlock()
	pic := pixel.PictureDataFromImage(g.screen)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))
}
