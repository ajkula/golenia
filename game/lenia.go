package game

import (
	"image/color"
	"math"
	"math/rand"
	"sync"

	"github.com/ajkula/golenia/graphics"
)

type Lenia struct {
	width, height int
	grid          [][]float64
	nextGrid      [][]float64
	mu            sync.Mutex
}

func NewLenia(width, height int) *Lenia {
	grid := make([][]float64, height)
	nextGrid := make([][]float64, height)
	for i := range grid {
		grid[i] = make([]float64, width)
		nextGrid[i] = make([]float64, width)
	}

	// init random des cellules pour visualisation
	for y := range grid {
		for x := range grid[y] {
			grid[y][x] = rand.Float64()
		}
	}

	return &Lenia{
		width:    width,
		height:   height,
		grid:     grid,
		nextGrid: nextGrid,
	}
}

func (l *Lenia) Update() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			neighborhood := l.getNeighborhood(x, y)
			l.nextGrid[y][x] = l.kernel(l.grid[y][x], neighborhood)
		}
	}

	l.grid, l.nextGrid = l.nextGrid, l.grid
}

func (l *Lenia) Render(g *graphics.Graphics) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			value := l.grid[y][x]
			color := color.RGBA{
				R: uint8(255 * value),
				G: uint8(100 * value),
				B: uint8(255 * (1 - value)),
				A: 255,
			}
			g.DrawCell(x, y, color)
		}
	}
}

func (l *Lenia) kernel(value, neighborhood float64) float64 {
	// transition pour Lenia
	sigma := 0.12
	mu := 0.3
	growthRate := 1.0
	k := neighborhood - mu
	return value + growthRate*sigma*k*math.Exp(-k*k)
}

func (l *Lenia) getNeighborhood(x, y int) float64 {
	// somme des cellules voisines avec une convolution gaussienne
	sum := 0.0
	weight := [][]float64{
		{0.02, 0.05, 0.02},
		{0.05, 0.32, 0.05},
		{0.02, 0.05, 0.02},
	}
	dirs := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 0}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for i, dir := range dirs {
		ny, nx := (y+dir.dy+l.height)%l.height, (x+dir.dx+l.width)%l.width
		sum += l.grid[ny][nx] * weight[i/3][i%3]
	}
	return sum
}
