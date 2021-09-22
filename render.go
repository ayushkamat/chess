package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const SQUARE_WIDTH = 100

const screenWidth = 8 * SQUARE_WIDTH
const screenHeight = 8 * SQUARE_WIDTH



func renderBoard(b Board, selectedPiece []int,  w *sdl.Window, r *sdl.Renderer) error {
	for i, file := range FILES {
		for j, rank := range RANKS {
			if (i + j) % 2 == 0 {
				r.SetDrawColor(248, 231, 187, 255)
			} else {
				r.SetDrawColor(0, 68, 116, 255)
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
	return nil
}