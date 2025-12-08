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
	t.Run("NextMovesへ置くとNextPieceとNextMovesが更新され、Winnerは更新されない", func(t *testing.T) {
		game := reversi.NewGame()
		err := game.PutPiece(2, 3, reversi.PieceBlack)

		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
		board := game.Board()
		if board[2][3] != reversi.CellBlack {
			t.Errorf("expected %v, got %v", reversi.CellBlack, board[2][3])
		}
		if game.NextPiece() != reversi.PieceWhite {
			t.Errorf("expected %v, got %v", reversi.PieceWhite, game.NextPiece())
		}
		expected := []reversi.Position{{2, 2}, {2, 4}, {4, 2}}
		nextMoves := game.NextMoves()
		if !reflect.DeepEqual(expected, nextMoves) {
			t.Errorf("expected %v, got %v", expected, nextMoves)
		}
		if game.Winner() != reversi.WinnerNone {
			t.Errorf("expected %v, got %v", reversi.WinnerNone, game.Winner())
		}
	})

	t.Run("黒勝利の最短決着", func(t *testing.T) {
		game := reversi.NewGame()
		tests := []struct {
			row    int
			col    int
			piece  reversi.Piece
			winner reversi.Winner
		}{
			{4, 5, reversi.PieceBlack, reversi.WinnerNone},
			{5, 3, reversi.PieceWhite, reversi.WinnerNone},
			{4, 2, reversi.PieceBlack, reversi.WinnerNone},
			{3, 5, reversi.PieceWhite, reversi.WinnerNone},
			{6, 4, reversi.PieceBlack, reversi.WinnerNone},
			{5, 5, reversi.PieceWhite, reversi.WinnerNone},
			{4, 6, reversi.PieceBlack, reversi.WinnerNone},
			{5, 4, reversi.PieceWhite, reversi.WinnerNone},
			{2, 4, reversi.PieceBlack, reversi.WinnerBlack},
		}
		for i, tt := range tests {
			err := game.PutPiece(tt.row, tt.col, tt.piece)
			if err != nil {
				t.Fatalf("turn %d err: expected nil, got %v", i, err)
			}
			if game.Winner() != tt.winner {
				t.Errorf("turn %d winner: expected %v, got %v", i, tt.winner, game.Winner())
			}
		}
	})
	t.Run("白勝利の最短決着", func(t *testing.T) {
		game := reversi.NewGame()
		tests := []struct {
			row    int
			col    int
			piece  reversi.Piece
			winner reversi.Winner
		}{
			{4, 5, reversi.PieceBlack, reversi.WinnerNone},
			{5, 5, reversi.PieceWhite, reversi.WinnerNone},
			{5, 4, reversi.PieceBlack, reversi.WinnerNone},
			{3, 5, reversi.PieceWhite, reversi.WinnerNone},
			{2, 4, reversi.PieceBlack, reversi.WinnerNone},
			{1, 3, reversi.PieceWhite, reversi.WinnerNone},
			{2, 3, reversi.PieceBlack, reversi.WinnerNone},
			{5, 3, reversi.PieceWhite, reversi.WinnerNone},
			{3, 2, reversi.PieceBlack, reversi.WinnerNone},
			{3, 1, reversi.PieceWhite, reversi.WinnerWhite},
		}
		for i, tt := range tests {
			err := game.PutPiece(tt.row, tt.col, tt.piece)
			if err != nil {
				t.Fatalf("turn %d err: expected nil, got %v", i, err)
			}
			if game.Winner() != tt.winner {
				t.Errorf("turn %d winner: expected %v, got %v", i, tt.winner, game.Winner())
			}
		}
	})
	t.Run("引き分け決着", func(t *testing.T) {
		game := reversi.NewGame()
		tests := []struct {
			row    int
			col    int
			piece  reversi.Piece
			winner reversi.Winner
		}{
			{5, 4, reversi.PieceBlack, reversi.WinnerNone},
			{5, 3, reversi.PieceWhite, reversi.WinnerNone},
			{2, 2, reversi.PieceBlack, reversi.WinnerNone},
			{5, 5, reversi.PieceWhite, reversi.WinnerNone},
			{6, 6, reversi.PieceBlack, reversi.WinnerNone},
			{5, 6, reversi.PieceWhite, reversi.WinnerNone},
			{5, 2, reversi.PieceBlack, reversi.WinnerNone},
			{7, 6, reversi.PieceWhite, reversi.WinnerNone},
			{5, 7, reversi.PieceBlack, reversi.WinnerNone},
			{1, 1, reversi.PieceWhite, reversi.WinnerNone},
			{7, 5, reversi.PieceBlack, reversi.WinnerNone},
			{4, 6, reversi.PieceWhite, reversi.WinnerNone},
			{7, 7, reversi.PieceBlack, reversi.WinnerNone},
			{4, 2, reversi.PieceWhite, reversi.WinnerNone},
			{3, 6, reversi.PieceBlack, reversi.WinnerNone},
			{6, 2, reversi.PieceWhite, reversi.WinnerNone},
			{5, 1, reversi.PieceBlack, reversi.WinnerNone},
			{5, 0, reversi.PieceWhite, reversi.WinnerNone},
			{0, 0, reversi.PieceBlack, reversi.WinnerNone},
			{2, 4, reversi.PieceWhite, reversi.WinnerDraw},
		}
		for i, tt := range tests {
			err := game.PutPiece(tt.row, tt.col, tt.piece)
			if err != nil {
				t.Fatalf("turn %d err: expected nil, got %v", i, err)
			}
			if game.Winner() != tt.winner {
				t.Errorf("turn %d winner: expected %v, got %v", i, tt.winner, game.Winner())
			}
		}
	})
}
