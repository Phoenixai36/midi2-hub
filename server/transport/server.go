package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"github.com/Phoenixai36/midi2-hub/server/session"
)

// Message is the universal wire format.
type Message struct {
	Type    string          `json:"type"`
	Room    string          `json:"room"`
	Client  string          `json:"client"`
	TS      int64           `json:"ts"`
	Payload json.RawMessage `json:"payload"`
}

type Server struct {
	addr    string
	manager *session.Manager
	mux     *http.ServeMux
}

func NewServer(addr string, manager *session.Manager) *Server {
	s := &Server{
		addr:    addr,
		manager: manager,
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("/ws", s.handleWS)
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // TODO: replace with proper origin check
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}
	defer conn.CloseNow()

	clientID := uuid.NewString()
	client := &session.Client{
		ID:   clientID,
		Send: make(chan []byte, 64),
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// writer goroutine
	go func() {
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					return
				}
				conn.Write(ctx, websocket.MessageText, msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	var roomName string
	for {
		var msg Message
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			break
		}
		msg.Client = clientID
		s.handleMessage(ctx, client, &msg, &roomName)
	}

	if roomName != "" {
		if room, ok := s.manager.GetRoom(roomName); ok {
			room.Leave(clientID)
		}
	}
}

func (s *Server) handleMessage(ctx context.Context, client *session.Client, msg *Message, roomName *string) {
	switch msg.Type {
	case "join":
		var p struct {
			Room    string `json:"room"`
			Profile string `json:"profile"`
		}
		json.Unmarshal(msg.Payload, &p)
		room := s.manager.GetOrCreateRoom(p.Room)
		client.Profile = p.Profile
		room.Join(client)
		*roomName = p.Room

	case "clock":
		var p struct {
			BPM   float64 `json:"bpm"`
			Beat  uint64  `json:"beat"`
			Phase float64 `json:"phase"`
		}
		json.Unmarshal(msg.Payload, &p)
		s.manager.Clock().ProposeBPM(p.BPM)
		// broadcast updated clock to all in room
		bpm, beat, phase := s.manager.Clock().Snapshot()
		out, _ := json.Marshal(Message{
			Type: "clock",
			Room: msg.Room,
			Payload: mustMarshal(map[string]interface{}{"bpm": bpm, "beat": beat, "phase": phase}),
		})
		s.manager.Broadcast(msg.Room, out, client.ID)

	case "clip", "control":
		out, _ := json.Marshal(msg)
		s.manager.Broadcast(msg.Room, out, client.ID)
	}
}

func mustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
