package session

import "sync"

// Client represents a connected producer.
type Client struct {
	ID      string
	Profile string
	Send    chan []byte
}

// Room is a collaboration session between producers.
type Room struct {
	mu      sync.RWMutex
	ID      string
	Name    string
	Members map[string]*Client
}

func (r *Room) Join(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Members[c.ID] = c
}

func (r *Room) Leave(clientID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Members, clientID)
}

func (r *Room) MemberCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Members)
}

// Broadcast sends msg to all members except the one with excludeID.
func (r *Room) Broadcast(msg []byte, excludeID string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for id, c := range r.Members {
		if id == excludeID {
			continue
		}
		select {
		case c.Send <- msg:
		default:
			// drop if buffer full
		}
	}
}
