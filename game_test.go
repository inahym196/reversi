package reversi_test

import (
	"reflect"
	"strings"
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

var initBoard = [][]reversi.Cell{
	{E, E, E, E, E, E, E, E},
	{E, E, E, E, E, E, E, E},
	{E, E, E, B, E, E, E, E},
	{E, E, E, B, B, E, E, E},
	{E, E, E, B, W, E, E, E},
	{E, E, E, E, E, E, E, E},
	{E, E, E, E, E, E, E, E},
	{E, E, E, E, E, E, E, E},
}

func TestNewBoard(t *testing.T) {
	board := reversi.NewBoard()

	t.Run("初期ボード", func(t *testing.T) {
		expected := initBoard
		assertBoardState(t, board, expected)
	})

	nextMoves := board.GetNextMoves(reversi.PieceBlack)
	expected := []reversi.Position{{2, 3}, {3, 2}, {4, 5}, {5, 4}}
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
			assertBoardState(t, board, initBoard)
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
}

func ParseMoves(s string) [][2]int {
	parts := strings.Fields(s)
	moves := make([][2]int, 0, len(parts))
	for _, p := range parts {
		col := int(p[0] - 'a')
		row := int(p[1] - '1')
		moves = append(moves, [2]int{row, col})
	}
	return moves
}

func ParseMovesCompact(s string) [][2]int {
	if len(s)%2 != 0 {
		return nil
	}
	moves := make([][2]int, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		colChar := s[i]
		rowChar := s[i+1]

		col := int(colChar - 'a')
		row := int(rowChar - '1')
		moves = append(moves, [2]int{row, col})
	}
	return moves
}

func RunScenario(t *testing.T, moves [][2]int, expectedWinner reversi.Winner) {
	t.Helper()

	game := reversi.NewGame()
	currentPiece := reversi.PieceBlack
	for i, move := range moves {
		row, col := move[0], move[1]
		err := game.PutPiece(row, col, currentPiece)
		if err != nil {
			t.Fatalf("turn %d err: expected nil, got %v", i, err)
		}

		got := game.Winner()
		if i < len(moves)-1 && got != reversi.WinnerNone {
			t.Fatalf("turn %d winner: expected %v, got %v", i, reversi.WinnerNone, got)
		} else if i == len(moves)-1 && got != expectedWinner {
			t.Errorf("last turn %d winner: expected %v, got %v", i, expectedWinner, got)
		}
		currentPiece = currentPiece.Opponent()
	}
}

func TestGame_PutPiece_Winner(t *testing.T) {
	tests := []struct {
		name           string
		record         string
		expectedWinner reversi.Winner
	}{
		{"黒勝利", "f5d6c5f4e7f6g5e6e3", reversi.WinnerBlack},
		{"白勝利", "f5f6e6f4e3d2d3d6c4b4", reversi.WinnerWhite},
		{"引き分け", "f5f4c3f6g7f7f3h7f8b2h6e7h8e3d7g3f2f1a1c5", reversi.WinnerDraw},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moves := ParseMovesCompact(tt.record)
			RunScenario(t, moves, tt.expectedWinner)
		})
	}
}
