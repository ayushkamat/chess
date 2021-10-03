package main

import (
	"errors"
	"fmt"
)

type Piece uint8
type Board map[int]map[int]Piece
type Move struct{
	PL int       // which player made the move              ([PL]ayer)
	SR int       // rank of the piece initially being moved ([S]ource [R]ank)
	SF int       // file of the piece initially being moved ([S]ource [F]ile)
	DR int       // rank of the destination square          ([D]estination [R]ank)
	DF int       // file of the destination square          ([D]estination [F]ile)
	P Piece      // piece to promote to, if applicable      ([P]romotion Piece)
}
type MoveSequence []Move

var FILES = [...]int {65, 66, 67, 68, 69, 70, 71, 72} // ASCII Values for ABCDEFGH 
var RANKS = [...]int {1, 2, 3, 4, 5, 6, 7, 8} 

// FORMAT: 
// [COLOR (1 bit)] 0 0 0 0 [PIECE TYPE (3 bits)]

const (
	EMPTY_SQUARE Piece = 0b00000000

	WHITE_PAWN Piece   = 0b00000001
	WHITE_KNIGHT Piece = 0b00000010
	WHITE_BISHOP Piece = 0b00000011
	WHITE_ROOK Piece   = 0b00000100
	WHITE_QUEEN Piece  = 0b00000101
	WHITE_KING Piece   = 0b00000110

	BLACK_PAWN Piece   = 0b10000001
	BLACK_KNIGHT Piece = 0b10000010
	BLACK_BISHOP Piece = 0b10000011
	BLACK_ROOK Piece   = 0b10000100
	BLACK_QUEEN Piece  = 0b10000101
	BLACK_KING Piece   = 0b10000110
)

var PIECE_NAMES = map[Piece]string{
	EMPTY_SQUARE : "",
	WHITE_PAWN   : "",
	WHITE_KNIGHT : "N",
	WHITE_BISHOP : "B",
	WHITE_ROOK   : "R",
	WHITE_QUEEN  : "Q",
	WHITE_KING   : "K",
	BLACK_PAWN   : "",
	BLACK_KNIGHT : "N",
	BLACK_BISHOP : "B",
	BLACK_ROOK   : "R",
	BLACK_QUEEN  : "Q",
	BLACK_KING   : "K"}

func (m Move) String() string {
	return fmt.Sprintf("%s%c%v", PIECE_NAMES[m.P], rune(m.DF), m.DR)
}

func (b Board) copy() Board {
	newBoard := Board{}
	for _, file := range FILES {
		newBoard[file] = make(map[int]Piece)
		for _, rank := range RANKS {
			newBoard[file][rank] = b[file][rank]
		}
	}
	return newBoard
}

func (h MoveSequence) copy() MoveSequence {
	newH := MoveSequence{}
	for _, move := range h {
		newH = append(newH, move)
	}
	return newH
}

func initializeBoard() (Board, error) {
	board := Board{}
	for _, file := range FILES {
		board[file] = make(map[int]Piece)
		for _, rank := range RANKS {
			if rank == 1 {
				if file == 'A' || file == 'H' {
					board[file][rank] = WHITE_ROOK
				} else if file == 'B' || file == 'G' {
					board[file][rank] = WHITE_KNIGHT
				} else if file == 'C' || file == 'F' {
					board[file][rank] = WHITE_BISHOP
				} else if file == 'D' {
					board[file][rank] = WHITE_QUEEN
				} else if file == 'E' {
					board[file][rank] = WHITE_KING
				} else {
					return nil, errors.New("Unexpected rank/file caused board initialization to break.")
				}
			} else if rank == 2 {
				board[file][rank] = WHITE_PAWN
			} else if rank == 7 {
				board[file][rank] = BLACK_PAWN
			} else if rank == 8 {
				if file == 'A' || file == 'H' {
					board[file][rank] = BLACK_ROOK
				} else if file == 'B' || file == 'G' {
					board[file][rank] = BLACK_KNIGHT
				} else if file == 'C' || file == 'F' {
					board[file][rank] = BLACK_BISHOP
				} else if file == 'D' {
					board[file][rank] = BLACK_QUEEN
				} else if file == 'E' {
					board[file][rank] = BLACK_KING
				} else {
					return nil, errors.New("Unexpected rank/file caused board initialization to break.")
				}
			} else {
				board[file][rank] = EMPTY_SQUARE
			}
		}
	}
	return board, nil
}

func getPath(p Piece) string {
	path := ""
	switch p {
	case WHITE_PAWN:
		path = "assets/Chess_plt45.bmp"
	case WHITE_KNIGHT:
		path = "assets/Chess_nlt45.bmp"
	case WHITE_BISHOP:
		path = "assets/Chess_blt45.bmp"
	case WHITE_ROOK:
		path = "assets/Chess_rlt45.bmp"
	case WHITE_QUEEN:
		path = "assets/Chess_qlt45.bmp"
	case WHITE_KING:
		path = "assets/Chess_klt45.bmp"
	case BLACK_PAWN:
		path = "assets/Chess_pdt45.bmp"
	case BLACK_KNIGHT:
		path = "assets/Chess_ndt45.bmp"
	case BLACK_BISHOP:
		path = "assets/Chess_bdt45.bmp"
	case BLACK_ROOK:
		path = "assets/Chess_rdt45.bmp"
	case BLACK_QUEEN:
		path = "assets/Chess_qdt45.bmp"
	case BLACK_KING:
		path = "assets/Chess_kdt45.bmp"
	}

	return path
}

func isWhite(p Piece) bool {
	return (p < 128) && (p > 0)
}

func isBlack(p Piece) bool {
	return p > 128
}

func threatens(b Board, h MoveSequence, targetFile int, targetRank int, p int, fromCheck bool) bool {
	// For a given position, determine if a particular square is being covered by one of p's pieces
	allied := isBlack
	if p == 0 {
		allied = isWhite
	}
	for _, file := range FILES {
		for _, rank := range RANKS {
			if allied(b[file][rank]) {
				moves := generateLegalMoves(b, h, file, rank, p, fromCheck)
				for _, move := range moves {
					if (move.DR == targetRank) && (move.DF == targetFile) {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkForCheck(b Board, h MoveSequence, p int) bool {
	// check if a player is in check by another player
	king := BLACK_KING
	if p == 0 {
		king = WHITE_KING
	}
	kingRank := 0
	kingFile := 0
	for _, file := range FILES {
		for _, rank := range RANKS {
			if b[file][rank] == king {
				kingRank = rank
				kingFile = file
			}
		}
	}

	return threatens(b, h, kingFile, kingRank, 1 - p, true)
}

func generateLegalMoves(b Board, h MoveSequence, file int, rank int, p int, fromCheck bool) MoveSequence {
	moves := make(MoveSequence, 0)
	if p == 0 {
		switch b[file][rank] {
		case WHITE_PAWN:
			// Cases if not in check: 
			// (1) FORWARD 1:  Legal iff the square in front of the pawn is empty and the pawn isn't pinned.
			// (2) FORWARD 2:  Legal iff the pawn is on the 2nd rank, the first two squares in front of the pawn are empty, and the pawn isn't pinned.
			// (3) CAPTURE:    Legal iff there is an opposing piece either up and to the left or up and to the right of the pawn, and the pawn isn't pinned.
			// (4) EN PASSANT: Legal iff an opposing pawn moved two squares forward to a square directly adjacent to the source pawn in the previous 
			//                 move, the square behind the opposing pawn is empty (vacuous), and the pawn isn't pinned.
			// (5a) PROMOTION: Legal iff the pawn is on the 7th rank and FORWARD 1 is legal
			// (5b) PROMOTION: Legal iff the pawn is on the 7th rank and CAPTURE is legal.

			// (1) FORWARD 1 / (5a) PROMOTION:
			if (rank + 1 <= 8) && (b[file][rank + 1] == EMPTY_SQUARE) {
				if rank == 7 {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 1, P: WHITE_KNIGHT})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 1, P: WHITE_BISHOP})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 1, P: WHITE_ROOK})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 1, P: WHITE_QUEEN})
				} else {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 1, P: WHITE_PAWN})
				}
			}

			// (2) FORWARD 2:
			if (rank == 2) && (b[file][rank + 1] == EMPTY_SQUARE) && (b[file][rank + 2] == EMPTY_SQUARE) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file, DR: rank + 2, P: WHITE_PAWN})
			}

			// (3) CAPTURE / (5b) PROMOTION:
			if (rank + 1 <= 8) && (file + 1 <= 'H') && (b[file + 1][rank + 1] > 128) {
				if rank == 7 {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_KNIGHT})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_BISHOP})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_ROOK})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_QUEEN})
				} else {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_PAWN})
				}
			}
			if (rank + 1 <= 8) && (file - 1 >= 'A') && (b[file - 1][rank + 1] > 128) {
				if rank == 7 {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_KNIGHT})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_BISHOP})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_ROOK})
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_QUEEN})
				} else {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_PAWN})
				}
			}

			// (4) EN PASSANT:
			if (rank == 5) && (file + 1 <= 'H') && (len(h) > 0) && (h[len(h) - 1] == Move{PL: 1, SF: file + 1, SR: 7, DF: file + 1, DR: 5, P: BLACK_PAWN}) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 1, P: WHITE_PAWN})
			}
			if (rank == 5) && (file - 1 >= 'A') && (len(h) > 0) && (h[len(h) - 1] == Move{PL: 1, SF: file - 1, SR: 7, DF: file - 1, DR: 5, P: BLACK_PAWN}) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 1, P: WHITE_PAWN})
			}
		case WHITE_KNIGHT:
			// Cases if not in check:
			// (1) JUMP:       Legal iff the target square is a knight's move away, is empty or has an opposing piece in it, and the knight isn't pinned.
			if (file + 1 <= 'H') && (rank + 2 <= 8) && ((b[file + 1][rank + 2] == EMPTY_SQUARE) || (b[file + 1][rank + 2] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank + 2, P: WHITE_KNIGHT})
			}
			if (file + 2 <= 'H') && (rank + 1 <= 8) && ((b[file + 2][rank + 1] == EMPTY_SQUARE) || (b[file + 2][rank + 1] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 2, DR: rank + 1, P: WHITE_KNIGHT})
			}
			if (file + 2 <= 'H') && (rank - 1 >= 1) && ((b[file + 2][rank - 1] == EMPTY_SQUARE) || (b[file + 2][rank - 1] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 2, DR: rank - 1, P: WHITE_KNIGHT})
			}
			if (file + 1 <= 'H') && (rank - 2 >= 1) && ((b[file + 1][rank - 2] == EMPTY_SQUARE) || (b[file + 1][rank - 2] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file + 1, DR: rank - 2, P: WHITE_KNIGHT})
			}
			if (file - 1 >= 'A') && (rank - 2 >= 1) && ((b[file - 1][rank - 2] == EMPTY_SQUARE) || (b[file - 1][rank - 2] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank - 2, P: WHITE_KNIGHT})
			}
			if (file - 2 >= 'A') && (rank - 1 >= 1) && ((b[file - 2][rank - 1] == EMPTY_SQUARE) || (b[file - 2][rank - 1] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 2, DR: rank - 1, P: WHITE_KNIGHT})
			}
			if (file - 2 >= 'A') && (rank + 1 <= 8) && ((b[file - 2][rank + 1] == EMPTY_SQUARE) || (b[file - 2][rank + 1] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 2, DR: rank + 1, P: WHITE_KNIGHT})
			}
			if (file - 1 >= 'A') && (rank + 2 <= 8) && ((b[file - 1][rank + 2] == EMPTY_SQUARE) || (b[file - 1][rank + 2] > 128)) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: file - 1, DR: rank + 2, P: WHITE_KNIGHT})
			}
		case WHITE_BISHOP:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on the same diagonal as the bishop, every square diagonally between the bishop and the target square 
			//                 is empty, the target square is empty or has an opposing piece, and the bishop isn't pinned.
			targetRank := rank + 1
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file + 1
			for (targetRank >= 0) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile -= 1
			}
			targetRank = rank + 1
			targetFile = file - 1
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile -= 1
			}
		case WHITE_ROOK:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on either the same rank or the same file as the rook, every square between the rook and the target 
			//                 square is empty, the target square is empty or has an opposing piece, and the rook isn't pinned.
			targetRank := rank
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile += 1
			}
			targetRank = rank
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile -= 1
			}
			targetRank = rank - 1
			targetFile = file
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
			}
			targetRank = rank + 1
			targetFile = file
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
			}
		case WHITE_QUEEN:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on either the same rank, the same file, or the same diagonal as the queen, every square between the
			//                 the queen and the target square is empty, the target square is empty or has an opposing piece, and the queen isn't pinned.
			targetRank := rank + 1
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file + 1
			for (targetRank >= 0) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile -= 1
			}
			targetRank = rank + 1
			targetFile = file - 1
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile -= 1
			}
			targetRank = rank
			targetFile = file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile += 1
			}
			targetRank = rank
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile -= 1
			}
			targetRank = rank - 1
			targetFile = file
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
			}
			targetRank = rank + 1
			targetFile = file
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isBlack(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
			}
		case WHITE_KING:
			targetRank := rank + 1
			targetFile := file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank + 1
			targetFile = file
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank + 1 
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank
			targetFile = file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank - 1
			targetFile = file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank - 1
			targetFile = file
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}
			targetRank = rank - 1 
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isBlack(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: WHITE_KING})
			}

			// CASTLING	
			if !fromCheck {
				kingside := !checkForCheck(b, h, p) && (b['H'][1] == WHITE_ROOK)
				queenside := !checkForCheck(b, h, p) && (b['A'][1] == WHITE_ROOK) && (b['B'][1] == EMPTY_SQUARE)
				if len(h) > 0 {
					for _, move := range h {
						kingside = kingside && (move.P != WHITE_KING) && !((move.P == WHITE_ROOK) && (move.SF == 'H') && (move.SR == 1))
						queenside = queenside && (move.P != WHITE_KING) && !((move.P == WHITE_ROOK) && (move.SF == 'A') && (move.SR == 1))
					}
				}
				for _, f := range []int{'F', 'G'} {
					kingside = kingside && !threatens(b, h, f, 1, 1 - p, true) && (b[f][1] == EMPTY_SQUARE)
				}
				for _, f := range []int{'C', 'D'} {
					queenside = queenside && !threatens(b, h, f, 1, 1 - p, true) && (b[f][1] == EMPTY_SQUARE)
				}
				if kingside {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: 'G', DR: 1, P: WHITE_KING})
				}
				if queenside {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: 'C', DR: 1, P: WHITE_KING})
				}
			}
		}
	} else {
		switch b[file][rank] {
		case BLACK_PAWN:
			// Cases if not in check: 
			// (1) FORWARD 1:  Legal iff the square in front of the pawn is empty and the pawn isn't pinned.
			// (2) FORWARD 2:  Legal iff the pawn is on the 2nd rank, the first two squares in front of the pawn are empty, and the pawn isn't pinned.
			// (3) CAPTURE:    Legal iff there is an opposing piece either up and to the left or up and to the right of the pawn, and the pawn isn't pinned.
			// (4) EN PASSANT: Legal iff an opposing pawn moved two squares forward to a square directly adjacent to the source pawn in the previous 
			//                 move, the square behind the opposing pawn is empty (vacuous), and the pawn isn't pinned.
			// (5a) PROMOTION: Legal iff the pawn is on the 7th rank and FORWARD 1 is legal
			// (5b) PROMOTION: Legal iff the pawn is on the 7th rank and CAPTURE is legal.

			// (1) FORWARD 1 / (5a) PROMOTION:
			if (rank - 1 >= 1) && (b[file][rank - 1] == EMPTY_SQUARE) {
				if rank == 2 {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 1, P: BLACK_KNIGHT})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 1, P: BLACK_BISHOP})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 1, P: BLACK_ROOK})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 1, P: BLACK_QUEEN})
				} else {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 1, P: BLACK_PAWN})
				}
			}

			// (2) FORWARD 2:
			if (rank == 7) && (b[file][rank - 1] == EMPTY_SQUARE) && (b[file][rank - 2] == EMPTY_SQUARE) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file, DR: rank - 2, P: BLACK_PAWN})
			}

			// (3) CAPTURE / (5b) PROMOTION:
			if (rank - 1 >= 1) && (file + 1 <= 'H') && (isWhite(b[file + 1][rank - 1])) {
				if rank == 2 {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_KNIGHT})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_BISHOP})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_ROOK})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_QUEEN})
				} else {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_PAWN})
				}
			}
			if (rank - 1 >= 1) && (file - 1 >= 'A') && (isWhite(b[file - 1][rank - 1])) {
				if rank == 2 {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_KNIGHT})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_BISHOP})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_ROOK})
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_QUEEN})
				} else {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_PAWN})
				}
			}

			// (4) EN PASSANT:
			if (rank == 4) && (file + 1 <= 'H') && (len(h) > 0) && (h[len(h) - 1] == Move{PL: 0, SF: file + 1, SR: 2, DF: file + 1, DR: 4, P: WHITE_PAWN}) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 1, P: BLACK_PAWN})
			}
			if (rank == 4) && (file - 1 >= 'A') && (len(h) > 0) && (h[len(h) - 1] == Move{PL: 0, SF: file - 1, SR: 2, DF: file - 1, DR: 4, P: WHITE_PAWN}) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 1, P: BLACK_PAWN})
			}
		case BLACK_KNIGHT:
			// Cases if not in check:
			// (1) JUMP:       Legal iff the target square is a knight's move away, is empty or has an opposing piece in it, and the knight isn't pinned.
			if (file + 1 <= 'H') && (rank + 2 <= 8) && ((b[file + 1][rank + 2] == EMPTY_SQUARE) || (b[file + 1][rank + 2] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank + 2, P: BLACK_KNIGHT})
			}
			if (file + 2 <= 'H') && (rank + 1 <= 8) && ((b[file + 2][rank + 1] == EMPTY_SQUARE) || (b[file + 2][rank + 1] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 2, DR: rank + 1, P: BLACK_KNIGHT})
			}
			if (file + 2 <= 'H') && (rank - 1 >= 1) && ((b[file + 2][rank - 1] == EMPTY_SQUARE) || (b[file + 2][rank - 1] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 2, DR: rank - 1, P: BLACK_KNIGHT})
			}
			if (file + 1 <= 'H') && (rank - 2 >= 1) && ((b[file + 1][rank - 2] == EMPTY_SQUARE) || (b[file + 1][rank - 2] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file + 1, DR: rank - 2, P: BLACK_KNIGHT})
			}
			if (file - 1 >= 'A') && (rank - 2 >= 1) && ((b[file - 1][rank - 2] == EMPTY_SQUARE) || (b[file - 1][rank - 2] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank - 2, P: BLACK_KNIGHT})
			}
			if (file - 2 >= 'A') && (rank - 1 >= 1) && ((b[file - 2][rank - 1] == EMPTY_SQUARE) || (b[file - 2][rank - 1] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 2, DR: rank - 1, P: BLACK_KNIGHT})
			}
			if (file - 2 >= 'A') && (rank + 1 <= 8) && ((b[file - 2][rank + 1] == EMPTY_SQUARE) || (b[file - 2][rank + 1] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 2, DR: rank + 1, P: BLACK_KNIGHT})
			}
			if (file - 1 >= 'A') && (rank + 2 <= 8) && ((b[file - 1][rank + 2] == EMPTY_SQUARE) || (b[file - 1][rank + 2] < 128)) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: file - 1, DR: rank + 2, P: BLACK_KNIGHT})
			}
		case BLACK_BISHOP:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on the same diagonal as the bishop, every square diagonally between the bishop and the target square 
			//                 is empty, the target square is empty or has an opposing piece, and the bishop isn't pinned.
			targetRank := rank + 1
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file + 1
			for (targetRank >= 0) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile -= 1
			}
			targetRank = rank + 1
			targetFile = file - 1
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_BISHOP})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile -= 1
			}
		case BLACK_ROOK:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on either the same rank or the same file as the rook, every square between the rook and the target 
			//                 square is empty, the target square is empty or has an opposing piece, and the rook isn't pinned.
			targetRank := rank
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile += 1
			}
			targetRank = rank
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile -= 1
			}
			targetRank = rank - 1
			targetFile = file
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
			}
			targetRank = rank + 1
			targetFile = file
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_ROOK})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
			}
		case BLACK_QUEEN:
			// Cases if not in check:
			// (1) MOVE:       Legal iff the target square is on either the same rank, the same file, or the same diagonal as the queen, every square between the
			//                 the queen and the target square is empty, the target square is empty or has an opposing piece, and the queen isn't pinned.
			targetRank := rank + 1
			targetFile := file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file + 1
			for (targetRank >= 0) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile += 1
			}
			targetRank = rank - 1
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
				targetFile -= 1
			}
			targetRank = rank + 1
			targetFile = file - 1
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
				targetFile -= 1
			}
			targetRank = rank
			targetFile = file + 1
			for (targetRank <= 8) && (targetFile <= 'H') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile += 1
			}
			targetRank = rank
			targetFile = file - 1
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetFile -= 1
			}
			targetRank = rank - 1
			targetFile = file
			for (targetRank >= 0) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank -= 1
			}
			targetRank = rank + 1
			targetFile = file
			for (targetRank <= 8) && (targetFile >= 'A') {
				if (b[targetFile][targetRank] == EMPTY_SQUARE) || (isWhite(b[targetFile][targetRank])) {
					moves = append(moves, Move{PL: 0, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_QUEEN})
				}
				if b[targetFile][targetRank] != EMPTY_SQUARE {
					break
				}
				targetRank += 1
			}
		case BLACK_KING:
			targetRank := rank + 1
			targetFile := file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank + 1
			targetFile = file
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank + 1 
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank
			targetFile = file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank - 1
			targetFile = file + 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank - 1
			targetFile = file
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}
			targetRank = rank - 1 
			targetFile = file - 1
			if (targetRank >= 1) && (targetRank <= 8) && (targetFile >= 'A') && (targetFile <= 'H') && ((b[targetFile][targetRank] == EMPTY_SQUARE) || isWhite(b[targetFile][targetRank])) {
				moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: targetFile, DR: targetRank, P: BLACK_KING})
			}

			// CASTLING
			if !fromCheck {
				kingside := !checkForCheck(b, h, p) && (b['H'][8] == BLACK_ROOK)
				queenside := !checkForCheck(b, h, p) && (b['A'][8] == BLACK_ROOK) && (b['B'][8] == EMPTY_SQUARE)
				if len(h) > 0 {
					for _, move := range h {
						kingside = kingside && (move.P != BLACK_KING) && !((move.P == BLACK_ROOK) && (move.SF == 'H') && (move.SR == 8))
						queenside = queenside && (move.P != BLACK_KING) && !((move.P == BLACK_ROOK) && (move.SF == 'A') && (move.SR == 8))
					}
				}
				for _, f := range []int{'F', 'G'} {
					kingside = kingside && !threatens(b, h, f, 8, 1 - p, true) && (b[f][8] == EMPTY_SQUARE)
				}
				for _, f := range []int{'C', 'D'} {
					queenside = queenside && !threatens(b, h, f, 8, 1 - p, true) && (b[f][8] == EMPTY_SQUARE)
				}
				if kingside {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: 'G', DR: 8, P: BLACK_KING})
				}
				if queenside {
					moves = append(moves, Move{PL: 1, SF: file, SR: rank, DF: 'C', DR: 8, P: BLACK_KING})
				}
			}
		}
	}
	if !fromCheck {
		newMoves := make(MoveSequence, 0)
		tempB := b.copy()
		for i := 0; i < len(moves); i++ {
			tempH := h.copy()
			makeMove(tempB, tempH, moves[i])
			if !checkForCheck(tempB, tempH, p) {
				newMoves = append(newMoves, moves[i])
			}
			undoMove(b, tempB, moves[i])
		}
		return newMoves
	}
	return moves
}

func checkNoLegalMoves(b Board, h MoveSequence, p int) bool {
	stalemate := true
	for _, file := range FILES {
		for _, rank := range RANKS {
			temp := generateLegalMoves(b, h, file, rank, p, false)
			stalemate = stalemate && (len(temp) == 0)
		}
	}
	return stalemate
}

func gameOver(b Board, h MoveSequence, p int) bool {
	return checkForCheck(b, h, p) && checkNoLegalMoves(b, h, p)
}

func makeMove(b Board, h MoveSequence, m Move) error {
	// if legal, err := CheckLegalMove(b, h, m); !legal || (err != nil) {
	// 	return errors.New("Illegal Move.")
	// }
	if (m.P == WHITE_KING) && (m.SF == 'E') && (m.SR == 1) && (m.DF == 'G') {
		// white kingside castling
		b['E'][1] = EMPTY_SQUARE
		b['H'][1] = EMPTY_SQUARE
		b['F'][1] = WHITE_ROOK
		b['G'][1] = WHITE_KING
	} else if (m.P == BLACK_KING) && (m.SF == 'E') && (m.SR == 8) && (m.DF == 'G') {
		// black kingside castling
		b['E'][8] = EMPTY_SQUARE
		b['H'][8] = EMPTY_SQUARE
		b['F'][8] = BLACK_ROOK
		b['G'][8] = BLACK_KING
	} else if (m.P == WHITE_KING) && (m.SF == 'E') && (m.SR == 1) && (m.DF == 'C') {
		// white queenside castling
		b['E'][1] = EMPTY_SQUARE
		b['A'][1] = EMPTY_SQUARE
		b['D'][1] = WHITE_ROOK
		b['C'][1] = WHITE_KING
	} else if (m.P == BLACK_KING) && (m.SF == 'E') && (m.SR == 8) && (m.DF == 'C') {
		// black queenside castling
		b['E'][8] = EMPTY_SQUARE
		b['A'][8] = EMPTY_SQUARE
		b['D'][8] = BLACK_ROOK
		b['C'][8] = BLACK_KING
	} else {
		b[m.DF][m.DR] = m.P
		b[m.SF][m.SR] = EMPTY_SQUARE
	}
	h = append(h, m)
	return nil
}

func undoMove(b Board, tempB Board, m Move) error {
	tempB[m.DF][m.DR] = b[m.DF][m.DR]
	tempB[m.SF][m.SR] = b[m.SF][m.SR]
	return nil
}