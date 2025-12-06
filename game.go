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

type Game struct {
	board Board
}

func NewGame() *Game {
	return &Game{NewBoard()}
}

func (g Game) Board() Board {
	return g.board
}
