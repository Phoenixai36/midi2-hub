package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Phoenixai36/midi2-hub/server/session"
	"github.com/Phoenixai36/midi2-hub/server/timesync"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A855F7"))
	bpmStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22D3EE"))
	roomStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ADE80"))
)

type model struct {
	manager *session.Manager
	clock   *timesync.Clock
	tick    int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tickMsg:
		m.tick++
		return m, tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	bpm, beat, phase := m.clock.Snapshot()
	rooms := m.manager.ListRooms()

	out := headerStyle.Render(" midi2-hub ") + "\n"
	out += fmt.Sprintf(" BPM: %s  Beat: %d  Phase: %.2f\n",
		bpmStyle.Render(fmt.Sprintf("%.1f", bpm)), beat, phase)
	out += "\n"
	out += headerStyle.Render(" Rooms ") + "\n"

	if len(rooms) == 0 {
		out += "  (no active rooms)\n"
	} else {
		for _, r := range rooms {
			out += fmt.Sprintf("  %s — %d members\n",
				roomStyle.Render(r.Name), r.MemberCount())
		}
	}
	out += "\n  [q] quit\n"
	return out
}

func Run(manager *session.Manager, clock *timesync.Clock) error {
	p := tea.NewProgram(model{manager: manager, clock: clock})
	_, err := p.Run()
	return err
}
