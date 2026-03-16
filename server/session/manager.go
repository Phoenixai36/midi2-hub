package session

import (
	"sync"

	"github.com/google/uuid"
	"github.com/Phoenixai36/midi2-hub/server/timesync"
)

// Manager handles all rooms and their members.
type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	clock *timesync.Clock
}

func NewManager(clock *timesync.Clock) *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
		clock: clock,
	}
}

func (m *Manager) GetOrCreateRoom(name string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.rooms[name]; ok {
		return r
	}
	r := &Room{
		ID:      uuid.NewString(),
		Name:    name,
		Members: make(map[string]*Client),
	}
	m.rooms[name] = r
	return r
}

func (m *Manager) GetRoom(name string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[name]
	return r, ok
}

func (m *Manager) ListRooms() []*Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rooms := make([]*Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

func (m *Manager) Clock() *timesync.Clock {
	return m.clock
}

func (m *Manager) Broadcast(roomName string, msg []byte, exclude string) {
	room, ok := m.GetRoom(roomName)
	if !ok {
		return
	}
	room.Broadcast(msg, exclude)
}
