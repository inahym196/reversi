package reversi_test

import (
	"testing"

	"github.com/inahym196/reversi"
)

func TestNewBoard(t *testing.T) {
	board := reversi.NewBoard()

	tests := []struct {
		x, y int
		want reversi.Cell
	}{
		{3, 3, reversi.CellWhite},
		{3, 4, reversi.CellBlack},
		{4, 3, reversi.CellBlack},
		{4, 4, reversi.CellWhite},
	}
	for _, tt := range tests {
		if got := board[tt.y][tt.x]; got != tt.want {
			t.Errorf("expected board[%d][%d] = %v, got %v", tt.y, tt.x, tt.want, got)
		}
	}

	for y := range reversi.BoardWidth {
		for x := range reversi.BoardWidth {
			if (y == 3 && x == 3) ||
				(y == 3 && x == 4) ||
				(y == 4 && x == 3) ||
				(y == 4 && x == 4) {
				continue
			}
			if board[y][x] != reversi.CellEmpty {
				t.Errorf("expected board[%d][%d] to be Empty, got %v", y, x, board[y][x])
			}
		}
	}

}
