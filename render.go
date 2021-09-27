package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const SQUARE_WIDTH = 100
const HIGHLIGHT = SQUARE_WIDTH / 3

const screenWidth = 8 * SQUARE_WIDTH
const screenHeight = 8 * SQUARE_WIDTH



func renderBoard(b Board, selectedPiece []int, highlightedSquares MoveSequence, w *sdl.Window, r *sdl.Renderer) error {
	for i, file := range FILES {
		for j, rank := range RANKS {
			if (i + j) % 2 == 0 {
				r.SetDrawColor(248, 231, 187, 255)
			} else {
				r.SetDrawColor(0, 68, 116, 255)
			}
			if (selectedPiece != nil) && (file == selectedPiece[0]) && (rank == selectedPiece[1]) {
				r.SetDrawColor(19, 196, 163, 255)
			}
			r.FillRect(&sdl.Rect{int32(SQUARE_WIDTH * i), int32(SQUARE_WIDTH * j), int32(SQUARE_WIDTH), int32(SQUARE_WIDTH)})

			path := getPath(b[file][rank])

			if path != "" {	
				img, err := sdl.LoadBMP(path)
				if err != nil {
					return err
				}
				defer img.Free()
				pieceTex, err := r.CreateTextureFromSurface(img)
				if err != nil {
					return err
				}
				defer pieceTex.Destroy()
				r.Copy(pieceTex, &sdl.Rect{0, 0, 141, 141}, &sdl.Rect{int32(SQUARE_WIDTH * i), int32(SQUARE_WIDTH * j), int32(SQUARE_WIDTH), int32(SQUARE_WIDTH)})
			}
		}
	}
	if highlightedSquares != nil {
		r.SetDrawColor(119, 136, 153, 255)
		for _, move := range highlightedSquares {
			i := move.DF - 'A'
			j := move.DR - 1
			r.FillRect(&sdl.Rect{int32(SQUARE_WIDTH * i + HIGHLIGHT), int32(SQUARE_WIDTH * j + HIGHLIGHT), int32(HIGHLIGHT), int32(HIGHLIGHT)})
		}
	}
	return nil
}