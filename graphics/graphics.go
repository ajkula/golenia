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

func (g *Graphics) Clear(c pixel.RGBA) {
	g.mu.Lock()
	defer g.mu.Unlock()
	draw.Draw(g.screen, g.screen.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func (g *Graphics) DrawCell(x, y int, c color.Color) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.screen.Set(x, y, c)
}

func (g *Graphics) Render(window *pixelgl.Window) {
	g.mu.Lock()
	defer g.mu.Unlock()
	pic := pixel.PictureDataFromImage(g.screen)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))
}

func (g *Graphics) DrawImage(img *image.RGBA) {
	g.mu.Lock()
	defer g.mu.Unlock()
	draw.Draw(g.screen, g.screen.Bounds(), img, image.Point{}, draw.Src)
}

func (g *Graphics) DrawImageScaled(img *image.RGBA, scale int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					g.screen.Set(x*scale+dx, y*scale+dy, c)
				}
			}
		}
	}
}
