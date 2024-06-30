package game

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/ajkula/golenia/graphics"
)

const (
	epsilon = 1e-5 // 1e-6 less / 1e-5 more diff to survive
)

type Lenia struct {
	width, height int
	grid          []float64
	nextGrid      []float64
	mu            sync.RWMutex
	speed         float64
	renderChan    chan *image.RGBA
	updateTicker  *time.Ticker
	kernel        []float64
}

func NewLenia(width, height int) *Lenia {
	l := &Lenia{
		width:        width,
		height:       height,
		grid:         make([]float64, width*height),
		nextGrid:     make([]float64, width*height),
		speed:        1.0,
		renderChan:   make(chan *image.RGBA, 30),
		updateTicker: time.NewTicker(time.Second / 60),
		kernel:       makeKernel(21), // kernel size
	}
	l.Reset()
	go l.updateLoop()
	return l
}

func makeKernel(size int) []float64 {
	kernel := make([]float64, size*size)
	center := size / 2
	sum := 0.0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx, dy := float64(x-center), float64(y-center)
			dist := math.Sqrt(dx*dx + dy*dy)
			value := math.Exp(-dist * dist / 20) // range
			kernel[y*size+x] = value
			sum += value
		}
	}
	for i := range kernel {
		kernel[i] /= sum
	}
	return kernel
}

func (l *Lenia) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := range l.grid {
		l.grid[i] = 0
	}

	for i := 0; i < 5; i++ {
		centerX, centerY := rand.Intn(l.width), rand.Intn(l.height)
		radius := 10
		for y := -radius; y <= radius; y++ {
			for x := -radius; x <= radius; x++ {
				if x*x+y*y <= radius*radius {
					value := rand.Float64()*0.5 + 0.5
					l.grid[((centerY+y+l.height)%l.height)*l.width+(centerX+x+l.width)%l.width] = value
				}
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

	const numGoroutines = 8
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for y := start; y < l.height; y += numGoroutines {
				for x := 0; x < l.width; x++ {
					neighborhood := l.getNeighborhood(x, y)
					l.nextGrid[y*l.width+x] = l.kernelFunction(l.grid[y*l.width+x], neighborhood)
					if l.nextGrid[y*l.width+x] < epsilon {
						l.nextGrid[y*l.width+x] = 0
					}
				}
			}
		}(i)
	}

	wg.Wait()
	l.grid, l.nextGrid = l.nextGrid, l.grid
}

func (l *Lenia) getNeighborhood(x, y int) float64 {
	sum := 0.0
	kernelSize := int(math.Sqrt(float64(len(l.kernel))))
	halfKernel := kernelSize / 2

	for ky := 0; ky < kernelSize; ky++ {
		for kx := 0; kx < kernelSize; kx++ {
			nx := (x + kx - halfKernel + l.width) % l.width
			ny := (y + ky - halfKernel + l.height) % l.height
			sum += l.grid[ny*l.width+nx] * l.kernel[ky*kernelSize+kx]
		}
	}

	return sum
}

func (l *Lenia) kernelFunction(value, neighborhood float64) float64 {
	mu := 0.05     // "comfort zone" center for survival
	sigma := 0.010 // width of that zone
	beta := 1 / (sigma * math.Sqrt(2*math.Pi))

	growth := beta * math.Exp(-math.Pow(neighborhood-mu, 2)/(2*sigma*sigma))

	growthRate := 0.05 * l.speed // growth rate
	newValue := value + growthRate*(2*growth-1)
	return math.Max(0, math.Min(1, newValue))
}

func (l *Lenia) renderFrame() {
	frame := image.NewRGBA(image.Rect(0, 0, l.width*4, l.height*4))

	l.mu.RLock()
	defer l.mu.RUnlock()

	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			value := l.grid[y*l.width+x]
			if value < epsilon {
				continue
			}
			r := uint8(255 * math.Pow(value, 0.5))
			g := uint8(255 * math.Pow(1-value, 2))
			b := uint8(255 * math.Sin(value*math.Pi))
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
	default:
	}
}

func (l *Lenia) Render(gr *graphics.Graphics) {
	select {
	case frame := <-l.renderChan:
		gr.DrawImage(frame)
	default:
	}
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
