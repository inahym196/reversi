package gamehub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/inahym196/gameHub/backend/pkg/hub"
	"github.com/inahym196/gameHub/backend/pkg/message"
	"github.com/inahym196/reversi"
)

type GamePayload struct {
	Board     reversi.Board      `json:"board"`
	NextPiece reversi.Piece      `json:"nextPiece"`
	NextMoves []reversi.Position `json:"nextMoves"`
	Winner    reversi.Winner     `json:"winner"`
}

func (GamePayload) MessageType() message.MessageType { return message.MessageTypeGame }

type Server struct {
	serveMux http.ServeMux
	game     *reversi.Game
	hub      *hub.Hub
}

func NewServer(hub *hub.Hub) *Server {
	s := &Server{
		hub:  hub,
		game: reversi.NewGame(),
	}
	s.serveMux.HandleFunc("/game", (s.gameHandlerFunc))
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

type PostGameRequest struct {
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Piece  string `json:"piece"`
}

func PieceFromString(s string) (reversi.Piece, error) {
	switch s {
	case "W":
		return reversi.PieceWhite, nil
	case "B":
		return reversi.PieceBlack, nil
	default:
		return false, fmt.Errorf("invalid string")
	}
}

func (s *Server) gameHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		payload := GamePayload{
			Board:     s.game.Board(),
			NextPiece: s.game.NextPiece(),
			NextMoves: s.game.NextMoves(),
			Winner:    s.game.Winner(),
		}
		res, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

	case http.MethodPost:
		body := http.MaxBytesReader(w, r.Body, 8192)
		var req PostGameRequest
		if err := json.NewDecoder(body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Print(req.Column, req.Row, req.Piece)
		piece, err := PieceFromString(req.Piece)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.game.PutPiece(req.Row, req.Column, piece); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		payload := GamePayload{
			Board:     s.game.Board(),
			NextPiece: s.game.NextPiece(),
			NextMoves: s.game.NextMoves(),
			Winner:    s.game.Winner(),
		}
		data, err := message.NewMessage(payload).Encode()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.hub.Enqueue(data)
		w.WriteHeader(http.StatusAccepted)

	case http.MethodDelete:
		s.game = reversi.NewGame()
		payload := GamePayload{
			Board:     s.game.Board(),
			NextPiece: s.game.NextPiece(),
			NextMoves: s.game.NextMoves(),
			Winner:    s.game.Winner(),
		}

		data, err := message.NewMessage(payload).Encode()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.hub.Enqueue(data)
		w.WriteHeader(http.StatusAccepted)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
