package game

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/ajkula/golenia/graphics"
)

type Lenia struct {
	width, height int
	grid          [][]float64
	nextGrid      [][]float64
	mu            sync.Mutex
	speed         float64
	renderChan    chan *image.RGBA
	updateTicker  *time.Ticker
}

func NewLenia(width, height int) *Lenia {
	l := &Lenia{
		width:        width,
		height:       height,
		grid:         make([][]float64, height),
		nextGrid:     make([][]float64, height),
		speed:        1.0,
		renderChan:   make(chan *image.RGBA, 10),       // N frames Buffer
		updateTicker: time.NewTicker(time.Second / 60), // 60 updates per sec
	}
	for i := range l.grid {
		l.grid[i] = make([]float64, width)
		l.nextGrid[i] = make([]float64, width)
	}
	l.Reset()
	go l.updateLoop() // backgrnd update loop init
	return l
}

func (l *Lenia) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Clear the grid
	for y := range l.grid {
		for x := range l.grid[y] {
			l.grid[y][x] = 0
		}
	}

	// Create a complex initial pattern
	centerX, centerY := l.width/2, l.height/2
	radius := 40
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			dist := math.Sqrt(float64(x*x + y*y))
			if dist <= float64(radius) {
				angle := math.Atan2(float64(y), float64(x))
				value := (math.Sin(dist/5)+math.Cos(angle*6))*0.25 + 0.5
				l.grid[(centerY+y+l.height)%l.height][(centerX+x+l.width)%l.width] = value
			}
		}
	}
}

func (l *Lenia) updateLoop() {
	for range l.updateTicker.C {
		l.Update()
		l.renderFrame()
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

func (l *Lenia) renderFrame() {
	frame := image.NewRGBA(image.Rect(0, 0, l.width*4, l.height*4))

	l.mu.Lock()
	defer l.mu.Unlock()

	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			value := l.grid[y][x]
			r := uint8(255 * math.Pow(value, 0.5))
			g := uint8(255 * math.Pow(value, 2))
			b := uint8(255 * value)
			c := color.RGBA{R: r, G: g, B: b, A: 255}
			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					frame.Set(x*4+i, y*4+j, c)
				}
			}
		}
	}

	select {
	case l.renderChan <- frame:
		// Frame successfuly sent
	default:
		// chan is full, ignore frame
	}
}

func (l *Lenia) Render(gr *graphics.Graphics) {
	select {
	case frame := <-l.renderChan:
		gr.DrawImage(frame)
	default:
	}
}

func (l *Lenia) kernel(value, neighborhood float64) float64 {
	r := 13.0
	mu := 0.15
	sigma := 0.016
	beta := 1 / (sigma * math.Sqrt(2*math.Pi))

	x := neighborhood * r
	growth := beta * math.Exp(-math.Pow(x-mu, 2)/(2*sigma*sigma))

	growthRate := 0.13 * l.speed
	newValue := value + growthRate*(2*growth-1)
	return math.Max(0, math.Min(1, newValue))
}

func (l *Lenia) getNeighborhood(x, y int) float64 {
	sum := 0.0
	weight := 0.0
	r := 13

	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			dist := math.Sqrt(float64(dx*dx + dy*dy))
			if dist <= float64(r) {
				nx := (x + dx + l.width) % l.width
				ny := (y + dy + l.height) % l.height
				k := math.Exp(-dist * dist / 40)
				sum += l.grid[ny][nx] * k
				weight += k
			}
		}
	}

	return sum / weight
}

func (l *Lenia) IncreaseSpeed() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.speed *= 1.1
	l.updateTicker.Reset(time.Duration(float64(time.Second) / (60 * l.speed)))
}

func (l *Lenia) DecreaseSpeed() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.speed /= 1.1
	l.updateTicker.Reset(time.Duration(float64(time.Second) / (60 * l.speed)))
}

func (l *Lenia) Cleanup() {
	l.updateTicker.Stop()
	close(l.renderChan)
}
