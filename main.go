package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH int32 = 1000
const HEIGHT int32 = 1000
const CELLSIZE int32 = 20

type cell struct {
	point sdl.Point
	state bool
	rect  sdl.Rect
}

var cells []cell

var mousePos sdl.Point
var pause bool = false
var grid bool
var cellColor sdl.Color = sdl.Color{R: 210, G: 100, B: 30, A: 255}

func nextGeneration(cells []cell) []cell {
	nextGen := []cell{}
	for _, cell := range cells {
		ns := checkNeighbours(cell, cells)
		//fmt.Println(ns)
		if cell.state == false && ns == 3 {
			cell.state = true
		} else if cell.state == true && (ns < 2 || ns > 3) {
			cell.state = false
		}
		nextGen = append(nextGen, cell)
	}
	return nextGen
}

func checkNeighbours(c cell, cells []cell) int32 {
	var ns int32 = 0
	for _, cell := range cells {
		if cell.state == true {
			if (cell.point.X+1 == c.point.X && cell.point.Y+1 == c.point.Y) ||
				(cell.point.X+1 == c.point.X && cell.point.Y == c.point.Y) ||
				(cell.point.X+1 == c.point.X && cell.point.Y-1 == c.point.Y) ||
				(cell.point.X == c.point.X && cell.point.Y+1 == c.point.Y) ||
				(cell.point.X == c.point.X && cell.point.Y-1 == c.point.Y) ||
				(cell.point.X-1 == c.point.X && cell.point.Y+1 == c.point.Y) ||
				(cell.point.X-1 == c.point.X && cell.point.Y == c.point.Y) ||
				(cell.point.X-1 == c.point.X && cell.point.Y-1 == c.point.Y) {
				ns += 1

			}
		}

	}
	return ns
}

func reset() {
	cells = []cell{}
	for i := 0; i < (int(WIDTH / CELLSIZE)); i++ {
		for j := 0; j < (int(HEIGHT / CELLSIZE)); j++ {
			cells = append(cells, cell{point: sdl.Point{X: int32(i), Y: int32(j)}, state: false, rect: sdl.Rect{X: int32(i) * CELLSIZE, Y: int32(j) * CELLSIZE, H: CELLSIZE, W: CELLSIZE}})
		}
	}
}

func random() {
	rand.Seed(time.Now().UnixNano())
	cells = []cell{}
	for i := 0; i < (int(WIDTH / CELLSIZE)); i++ {
		for j := 0; j < (int(HEIGHT / CELLSIZE)); j++ {
			var state bool = false
			if rand.Intn(2) == 1 {
				state = true
			}
			cells = append(cells, cell{point: sdl.Point{X: int32(i), Y: int32(j)}, state: state, rect: sdl.Rect{X: int32(i) * CELLSIZE, Y: int32(j) * CELLSIZE, H: CELLSIZE, W: CELLSIZE}})
		}
	}
}

func run() (err error) {
	var window *sdl.Window

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow("Interactive Game of life -____-", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, 0, sdl.GL_ACCELERATED_VISUAL)
	if err != nil {
		return err
	}
	defer renderer.Destroy()
	running := true
	for running {
		//events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				{
					running = false
				}
			case *sdl.MouseButtonEvent:
				{
					if t.Button == sdl.BUTTON_LEFT && t.State == sdl.PRESSED {
						mousePos = sdl.Point{X: t.X, Y: t.Y}
						for i, cello := range cells {
							if mousePos.InRect(&cello.rect) {
								cells[i] = cell{point: cello.point, state: !cello.state, rect: cello.rect}
								break
							}
						}

					}
				}
			case *sdl.KeyboardEvent:
				{
					//fmt.Println(t)
					if t.State == sdl.PRESSED {
						switch t.Keysym.Scancode {
						case 10: //G - turns grid on or off
							{
								grid = !grid
								break
							}
						case 21: //R - sets all cells to false
							{
								reset()
								break
							}
						case 23: //T - sets all cells to a random state
							{
								random()
								break
							}
						case 44: //SPACEBAR - pauses the game, the pause also turns the grid on
							{
								pause = !pause
								break
							}
						}
						break
					}

				}
			}
		}
		//drawing
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		for _, cell := range cells {
			if cell.state {
				renderer.SetDrawColor(cellColor.R, cellColor.G, cellColor.B, cellColor.A)
				renderer.FillRect(&cell.rect)
			}
		}
		if pause || grid { //*WHEN UNCOMMENTED THE GRID WILL ONLY SHOW WHEN THE GAME IS PAUSED
			renderer.SetDrawColor(0, 0, 0, 255)
			for i := 0; i < int(WIDTH/CELLSIZE); i++ {
				renderer.DrawLine(int32(i)*CELLSIZE, 0, int32(i)*CELLSIZE, HEIGHT)
			}
			for i := 0; i < int(HEIGHT/CELLSIZE); i++ {
				renderer.DrawLine(0, int32(i)*CELLSIZE, WIDTH, int32(i)*CELLSIZE)
			}
		} //*
		renderer.Present()

		//calculationg next generation
		if !pause {
			cells = nextGeneration(cells)
		}
		//sdl.Delay(1600)
	}
	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reset()
	if err := run(); err != nil {
		os.Exit(1)
	}
}
