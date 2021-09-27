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

	var selectedPiece []int = nil
	var tempPiece []int = nil
	var legalMoves MoveSequence = nil

	mousePressed := false
	moveMade := false

	player := 0
	isPlayer := isWhite
	isOpponent := isBlack

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return nil
			case *sdl.KeyboardEvent:
				if t.GetType() == sdl.KEYUP {
					if count >= len(moves) {
						return nil
					}
					makeMove(b, h, moves[count])
					count += 1
					fmt.Println(checkForCheck(b, h, 1))
				}
			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED && !mousePressed {
					tempPiece = []int{int(t.X / SQUARE_WIDTH) + 65, int(t.Y / SQUARE_WIDTH) + 1}
					if selectedPiece != nil {
						for _, move := range legalMoves {
							if (tempPiece[0] == move.DF) && (tempPiece[1] == move.DR) {
								makeMove(b, h, move)
								moveMade = true
								player = 1 - player
								isPlayer, isOpponent = isOpponent, isPlayer
								selectedPiece = nil
								legalMoves = nil
								break
							}
						}
						selectedPiece = nil
						legalMoves = nil
					}
					if !moveMade {
						if b[tempPiece[0]][tempPiece[1]] == EMPTY_SQUARE || isOpponent(b[tempPiece[0]][tempPiece[1]]){
							selectedPiece = nil
							legalMoves = nil
						} else {
							selectedPiece = tempPiece
							legalMoves = generateLegalMoves(b, h, selectedPiece[0], selectedPiece[1], player, false)
						}
						mousePressed = true
					}
					moveMade = false
				}
				if t.State == sdl.RELEASED {
					mousePressed = false
				}
			}

		}


		err = renderBoard(b, selectedPiece, legalMoves, window, renderer)
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