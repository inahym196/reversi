package reversi

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
	b[row][col] = cellFromPiece(piece)
	return nil
}

type Game struct {
	board Board
}

func NewGame() *Game {
	return &Game{NewBoard()}
}

func (g Game) Board() Board {
	return g.board
}
