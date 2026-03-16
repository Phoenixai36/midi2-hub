package timesync

import (
	"sync"
	"time"
)

// Clock maintains the global tempo and beat position.
// Implements Link-style consensus: any peer can propose a tempo change;
// the average of recent proposals converges to a shared tempo.
type Clock struct {
	mu       sync.RWMutex
	bpm      float64
	beat     uint64
	phase    float64
	started  time.Time
	proposals []float64
}

func NewClock(initialBPM float64) *Clock {
	return &Clock{
		bpm:     initialBPM,
		started: time.Now(),
	}
}

func (c *Clock) BPM() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bpm
}

// ProposeBPM allows a client to suggest a tempo.
// The clock converges by averaging recent proposals (simplified Link-style).
func (c *Clock) ProposeBPM(bpm float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.proposals = append(c.proposals, bpm)
	if len(c.proposals) > 8 {
		c.proposals = c.proposals[len(c.proposals)-8:]
	}
	var sum float64
	for _, p := range c.proposals {
		sum += p
	}
	c.bpm = sum / float64(len(c.proposals))
}

func (c *Clock) SetBPM(bpm float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bpm = bpm
	c.proposals = []float64{bpm}
}

// Phase returns the current beat phase [0,1).
func (c *Clock) Phase() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elapsed := time.Since(c.started).Seconds()
	beatsPerSec := c.bpm / 60.0
	totalBeats := elapsed * beatsPerSec
	return totalBeats - float64(uint64(totalBeats))
}

func (c *Clock) Beat() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elapsed := time.Since(c.started).Seconds()
	beatsPerSec := c.bpm / 60.0
	return uint64(elapsed * beatsPerSec)
}

// Snapshot returns a consistent (bpm, beat, phase) tuple.
func (c *Clock) Snapshot() (bpm float64, beat uint64, phase float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	elapsed := time.Since(c.started).Seconds()
	beatsPerSec := c.bpm / 60.0
	totalBeats := elapsed * beatsPerSec
	return c.bpm, uint64(totalBeats), totalBeats - float64(uint64(totalBeats))
}
