package reversi

import (
	"fmt"
)

type Piece bool

const (
	PieceBlack Piece = true
	PieceWhite Piece = false
)

func (p Piece) Opponent() Piece {
	return !p
}

func (p Piece) String() string {
	switch p {
	case PieceBlack:
		return "PieceBlack"
	case PieceWhite:
		return "PieceWhite"
	default:
		panic("invalid piece")
	}
}

type Cell byte

const (
	CellEmpty Cell = iota
	CellBlack
	CellWhite
)

func (c Cell) String() string {
	switch c {
	case CellEmpty:
		return "CellEmpty"
	case CellBlack:
		return "CellBlack"
	case CellWhite:
		return "CellWhite"
	default:
		panic("invalid cell")
	}
}

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

type Position struct {
	Row, Column int
}

func (b *Board) isInBoard(row, column int) bool {
	return 0 <= row && row < BoardWidth && 0 <= column && column < BoardWidth
}

func (b *Board) collectFlippableInDirection(row, column, dy, dx int, piece Piece) (flips []Position) {
	row += dy
	column += dx
	for b.isInBoard(row, column) {
		switch b[row][column] {
		case cellFromPiece(piece.Opponent()):
			flips = append(flips, Position{row, column})
		case cellFromPiece(piece):
			return flips
		case CellEmpty:
			return []Position{}
		}
		row += dy
		column += dx
	}
	return []Position{}
}

func (b *Board) collectFlippable(row, col int, piece Piece) (flips []Position) {
	dirs := []struct{ x, y int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, dir := range dirs {
		dx, dy := dir.x, dir.y
		flipsInDir := b.collectFlippableInDirection(row, col, dy, dx, piece)
		flips = append(flips, flipsInDir...)
	}
	return flips
}

func (b *Board) isPlaced(row, column int) bool {
	return b[row][column] != CellEmpty
}

func (b *Board) PutPiece(row, col int, piece Piece) error {
	if b.isPlaced(row, col) {
		return fmt.Errorf("cell is not empty")
	}
	flips := b.collectFlippable(row, col, piece)
	if len(flips) == 0 {
		return fmt.Errorf("flippable piece not exists")
	}
	for _, flip := range flips {
		b[flip.Row][flip.Column] = cellFromPiece(piece)
	}
	b[row][col] = cellFromPiece(piece)
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

type Winner byte

const (
	WinnerNone Winner = iota
	WinnerWhite
	WinnerBlack
)

func (w Winner) String() string {
	switch w {
	case WinnerNone:
		return "WinnerNone"
	case WinnerWhite:
		return "WinnerWhite"
	case WinnerBlack:
		return "WinnerBlack"
	default:
		panic("invalid winner")
	}
}

type Game struct {
	board     Board
	nextPiece Piece
	nextMoves []Position
	winner    Winner
}

func NewGame() *Game {
	board := NewBoard()
	return &Game{
		board:     board,
		nextPiece: PieceBlack,
		nextMoves: board.GetNextMoves(PieceBlack),
		winner:    Winner(WinnerNone),
	}
}

func (g Game) Board() Board {
	return g.board
}

func (g Game) NextPiece() Piece {
	return g.nextPiece
}

func (g Game) NextMoves() []Position {
	return g.nextMoves
}

func (g Game) Winner() Winner {
	return g.winner
}

func (g *Game) isInNextMoves(row, col int) bool {
	for _, move := range g.nextMoves {
		if row == move.Row && col == move.Column {
			return true
		}
	}
	return false
}

func (g *Game) PutPiece(row, col int, piece Piece) error {
	if g.nextPiece != piece {
		return fmt.Errorf("相手のターンです")
	}
	if !g.isInNextMoves(row, col) {
		return fmt.Errorf("無効な配置場所です")
	}
	g.board.PutPiece(row, col, piece)
	g.nextPiece = g.nextPiece.Opponent()
	g.nextMoves = g.board.GetNextMoves(g.nextPiece)
	return nil
}
