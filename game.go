package reversi

import "fmt"

type Piece bool

const (
	PieceBlack Piece = true
	PieceWhite Piece = false
)

func (p Piece) Opponent() Piece {
	return !p
}

type Cell byte

const (
	CellEmpty Cell = iota
	CellBlack
	CellWhite
)

func cellFromPiece(p Piece) Cell {
	switch p {
	case PieceBlack:
		return CellBlack
	case PieceWhite:
		return CellWhite
	default:
		panic("invalid piece")
	}
}

const (
	BoardWidth = 8
)

type Board [BoardWidth][BoardWidth]Cell

func NewBoard() Board {
	b := Board{}
	b[3][3] = CellWhite
	b[3][4] = CellBlack
	b[4][3] = CellBlack
	b[4][4] = CellWhite
	return b
}

func (b *Board) PutPiece(row, col int, piece Piece) error {
	if b[row][col] != CellEmpty {
		return fmt.Errorf("cell is not empty")
	}
	b[row][col] = cellFromPiece(piece)
	return nil
}

type Position struct {
	Row, Column int
}

func (b *Board) isPlaced(row, column int) bool {
	return b[row][column] != CellEmpty
}

func (b *Board) isInBoard(row, column int) bool {
	return 0 <= row && row < BoardWidth && 0 <= column && column < BoardWidth
}

func (b *Board) collectFlippableInDirection(row, column, dy, dx int, piece Piece) (flips []Position) {
	row += dy
	column += dx
	for b.isInBoard(row, column) {
		switch b[row][column] {
		case CellEmpty:
			return nil
		case cellFromPiece(piece):
			return flips
		default:
			flips = append(flips, Position{row, column})
		}
		row += dy
		column += dx
	}
	return nil
}

func (b *Board) canFlipPieces(row, column int, piece Piece) bool {
	dirs := []struct{ x, y int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, dir := range dirs {
		flips := b.collectFlippableInDirection(row, column, dir.y, dir.x, piece)
		if len(flips) > 0 {
			return true
		}
	}
	return false
}

func (b *Board) GetNextMoves(piece Piece) (nextMoves []Position) {
	for row := range BoardWidth {
		for col := range BoardWidth {
			if !b.isPlaced(row, col) && b.canFlipPieces(row, col, piece) {
				nextMoves = append(nextMoves, Position{row, col})
			}
		}
	}
	return nextMoves
}

type Game struct {
	board     Board
	nextMoves []Position
}

func NewGame() *Game {
	b := NewBoard()
	moves := b.GetNextMoves(PieceBlack)
	return &Game{b, moves}
}

func (g Game) Board() Board {
	return g.board
}
