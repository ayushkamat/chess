package main

import (
	"os"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Error initializing SDL:", err)
		return err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Chess",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		screenWidth,
		screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("Error creating window:", err)
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Error initializing renderer:", err)
		return err
	}
	defer renderer.Destroy()

	b, err := initializeBoard()
	if err != nil {
		fmt.Println("Board is broken:", err)
		return err
	}

	moves := SCHOLAR_MATE
	h := make([]Move, 0)
	count := 0

	selectedPiece := []int{0, 0}

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return nil
			case *sdl.KeyboardEvent:
				if event.GetType() == sdl.KEYUP {
					if count >= len(moves) {
						return nil
					}
					makeMove(b, h, moves[count])
					count += 1
					fmt.Println(checkForCheck(b, h, 1))
				}
			case *sdl.MouseButtonEvent:
				if event.State == sdl.PRESSED {
					if selectedPiece == []int{0, 0} {
						selectedPiece = []int{event.(type).Y / SQUARE_WIDTH + 1, event.(type).X / SQUARE_WIDTH + 65}
					}
					fmt.Println(selectedPiece)
				}
			}

		}
		err = renderBoard(b, selectedPiece, window, renderer)
		if err != nil {
			fmt.Println("Board is broken:", err)
			return err
		}
		renderer.Present()
	}
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}