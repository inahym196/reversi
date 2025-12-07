package reversi_test

import (
	"reflect"
	"testing"

	"github.com/inahym196/reversi"
)

const (
	E = reversi.CellEmpty
	B = reversi.CellBlack
	W = reversi.CellWhite
)

func assertBoardState(t *testing.T, board reversi.Board, expected [][]reversi.Cell) {
	t.Helper()

	for row := range expected {
		for col := range expected[row] {
			if board[row][col] != expected[row][col] {
				t.Errorf("board[%d][%d]: expected %v, got %v",
					row, col, expected[row][col], board[row][col])
			}
		}
	}
}

func TestNewBoard(t *testing.T) {
	board := reversi.NewBoard()

	t.Run("初期ボード", func(t *testing.T) {
		expected := [][]reversi.Cell{
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, W, B, E, E, E},
			{E, E, E, B, W, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
		}
		assertBoardState(t, board, expected)
	})

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

	t.Run("初期配置から一つ置く", func(t *testing.T) {
		board := reversi.NewBoard()

		board.PutPiece(2, 3, reversi.PieceBlack)

		expected := [][]reversi.Cell{
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, B, E, E, E, E},
			{E, E, E, B, B, E, E, E},
			{E, E, E, B, W, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
			{E, E, E, E, E, E, E, E},
		}
		assertBoardState(t, board, expected)
	})
	t.Run("NextMoves以外には配置できない", func(t *testing.T) {
		tests := []struct {
			row, column int
			cell        reversi.Cell
		}{
			{1, 1, reversi.CellEmpty},
			{3, 3, reversi.CellWhite},
		}

		board := reversi.NewBoard()

		for _, tt := range tests {
			for _, moves := range board.GetNextMoves(reversi.PieceBlack) {
				p := reversi.Position{tt.row, tt.column}
				if p == moves {
					t.Fatal("unexpected error")
				}
			}
		}

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

func TestNewGame(t *testing.T) {
	game := reversi.NewGame()

	t.Run("最初のnextMovesは4個", func(t *testing.T) {
		expected := []reversi.Position{{2, 3}, {3, 2}, {4, 5}, {5, 4}}
		nextMoves := game.NextMoves()
		if !reflect.DeepEqual(expected, nextMoves) {
			t.Errorf("expected %v, got %v", expected, nextMoves)
		}
	})

	t.Run("最初のnextPieceはBlack", func(t *testing.T) {
		nextPiece := game.NextPiece()
		if nextPiece != reversi.PieceBlack {
			t.Errorf("expected %v, got %v", reversi.PieceBlack, nextPiece)
		}
	})
}

func TestGame_PutPiece(t *testing.T) {
	t.Run("nextPiece以外は置けない", func(t *testing.T) {
		game := reversi.NewGame()
		err := game.PutPiece(2, 3, reversi.PieceWhite)
		if err == nil {
			t.Error("エラーが出るはずなのに出ていない")
		}
	})
	t.Run("ピースを置くと色が変わる", func(t *testing.T) {
		game := reversi.NewGame()
		err := game.PutPiece(2, 3, reversi.PieceBlack)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		if game.NextPiece() != reversi.PieceWhite {
			t.Errorf("expected %v, got %v", reversi.PieceWhite, game.NextPiece())
		}
	})
	t.Run("ピースを置くとnextMovesが変わる", func(t *testing.T) {
		game := reversi.NewGame()
		err := game.PutPiece(2, 3, reversi.PieceBlack)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		expected := []reversi.Position{{2, 2}, {2, 4}, {4, 2}}
		nextMoves := game.NextMoves()
		if !reflect.DeepEqual(expected, nextMoves) {
			t.Errorf("expected %v, got %v", expected, nextMoves)
		}
	})
}
