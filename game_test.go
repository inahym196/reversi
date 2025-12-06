package reversi_test

import (
	"testing"

	"github.com/inahym196/reversi"
)

func TestNewBoard(t *testing.T) {
	board := reversi.NewBoard()

	tests := []struct {
		row, col int
		want     reversi.Cell
	}{
		{3, 3, reversi.CellWhite},
		{3, 4, reversi.CellBlack},
		{4, 3, reversi.CellBlack},
		{4, 4, reversi.CellWhite},
	}
	for _, tt := range tests {
		if got := board[tt.row][tt.col]; got != tt.want {
			t.Errorf("expected board[%d][%d] = %v, got %v", tt.row, tt.col, tt.want, got)
		}
	}

	for row := range reversi.BoardWidth {
		for col := range reversi.BoardWidth {
			if (row == 3 && col == 3) ||
				(row == 3 && col == 4) ||
				(row == 4 && col == 3) ||
				(row == 4 && col == 4) {
				continue
			}
			if board[row][col] != reversi.CellEmpty {
				t.Errorf("expected board[%d][%d] to be Empty, got %v", row, col, board[row][col])
			}
		}
	}
	nextMoves := board.GetNextMoves(reversi.PieceBlack)
	expected := []reversi.Position{
		{Row: 2, Column: 3},
		{Row: 3, Column: 2},
		{Row: 4, Column: 5},
		{Row: 5, Column: 4},
	}
	for i := range nextMoves {
		if nextMoves[i] != expected[i] {
			t.Error("error")
		}
	}
}

func TestBoard_PutPiece(t *testing.T) {

	t.Run("空白セルに配置できる", func(t *testing.T) {
		board := reversi.NewBoard()
		if board[1][1] != reversi.CellEmpty {
			t.Fatalf("invalid board")
		}

		board.PutPiece(1, 1, reversi.PieceBlack)

		if board[1][1] != reversi.CellBlack {
			t.Errorf("expected %d, got %d", reversi.CellBlack, board[1][1])
		}
	})
	t.Run("空白以外のセルに配置できない", func(t *testing.T) {
		tests := []struct {
			row, column int
			cell        reversi.Cell
		}{
			{3, 3, reversi.CellWhite},
			{3, 4, reversi.CellBlack},
			{4, 3, reversi.CellBlack},
			{4, 4, reversi.CellWhite},
		}
		board := reversi.NewBoard()

		for _, tt := range tests {
			err := board.PutPiece(tt.row, tt.column, reversi.PieceBlack)
			if err == nil {
				t.Errorf("expected something error, got %v", err)
			}
			if board[tt.row][tt.column] != tt.cell {
				t.Errorf("expected %d, got %d", tt.cell, board[tt.row][tt.column])
			}

		}

	})
}
