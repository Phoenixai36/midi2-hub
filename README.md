# midi2-hub

> **MIDI 2.0 AI Remix Hub** — Multi-producer sync over Network MIDI 2.0 with AI-powered remix engines.
> Built for producers across cities and genres. Open source, modular, sprint-ready.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org)
[![C++17](https://img.shields.io/badge/C++-17-orange.svg)](https://isocpp.org)
[![MIDI 2.0](https://img.shields.io/badge/MIDI-2.0-purple.svg)](https://midi.org)
[![Status](https://img.shields.io/badge/status-Sprint%201%20active-green.svg)](#roadmap)

---

## What is this?

`midi2-hub` is a **modular hub** that connects music producers across the world using **Network MIDI 2.0** as the sync backbone, with **AI engines** (stem separation, audio-to-MIDI, pattern generation, emotional macros) running as isolated microservices.

It is designed as:
- A **Go server** (WebSocket + TUI) handling session management, clock consensus, clip routing
- A **JUCE VST3/CLAP plugin** (C++17) acting as the DAW bridge for any producer
- **Python microservices** (Podman containers) for AI: audio-to-MIDI, pattern generation, emotional engine
- A **Network MIDI 2.0 layer** as the universal sync and control protocol

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    midi2-hub SERVER (Go)                     │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ SessionMgr   │  │ TimeSyncCore │  │  ClipService     │  │
│  │ rooms/routing│  │ consensus BPM│  │  JSON→SMF2       │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐                         │
│  │  WS Transport│  │  TUI (Bubble │                         │
│  │  TLS/WSS     │  │  Tea)        │                         │
│  └──────────────┘  └──────────────┘                         │
└───────────────────────┬─────────────────────────────────────┘
                        │ WebSocket + JSON (→ SMF2)
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ JUCE Plugin  │ │ JUCE Plugin  │ │ JUCE Plugin  │
│ Producer BCN │ │ Producer BER │ │ Producer TKY │
│ VST3/CLAP    │ │ VST3/CLAP    │ │ VST3/CLAP    │
│ Ableton/DAW  │ │ Bitwig/DAW   │ │ Reaper/DAW   │
└──────┬───────┘ └──────────────┘ └──────────────┘
       │
       ▼ HTTP/WS
┌─────────────────────────────────────┐
│         AI Engines (Podman)         │
│  ┌────────────┐ ┌────────────────┐  │
│  │audio-to-   │ │ pattern-gen    │  │
│  │midi        │ │ (text-to-groove│  │
│  │(basic-pitch│ │  by style tag) │  │
│  └────────────┘ └────────────────┘  │
│  ┌─────────────────────────────┐    │
│  │ emotional-engine            │    │
│  │ (tension/density/brightness │    │
│  │  → MIDI 2.0 CC 32-bit)     │    │
│  └─────────────────────────────┘    │
└─────────────────────────────────────┘
```

---

## Stack

| Layer | Technology | Role |
|-------|-----------|------|
| Core MIDI 2.0 | **C++17** + ni-midi2 + libremidi | UMP parsing, MIDI-CI, I/O |
| Plugin DAW | **JUCE** (VST3/CLAP) | DAWBridge for Ableton and others |
| Hub Server | **Go 1.22** | SessionManager, routing, WebSocket |
| TUI | **Go + Bubble Tea** | Interactive CLI for debug and control |
| Transport | **WebSocket + JSON** → SMF2 | Clips, control, messages between nodes |
| AI Engines | **Python 3.12** in Podman | AudioToMidi, PatternGen, EmotionalEngine |
| Sync | **Link-style consensus** | Distributed, no single master |
| Clips | **SMF2** target (JSON in MVP) | Standard format |
| Network | **Internet** (WebSocket over TLS) | Producers in different cities |

---

## Project Structure

```
midi2-hub/
├── server/                    # Go Hub server
│   ├── main.go
│   ├── session/               # SessionManager: rooms, members, routing
│   │   ├── manager.go
│   │   └── room.go
│   ├── transport/             # WebSocket server, TLS
│   │   ├── server.go
│   │   └── client.go
│   ├── timesync/              # Consensus clock (Link-style)
│   │   └── clock.go
│   ├── clips/                 # Clip store JSON→SMF2
│   │   └── store.go
│   ├── bridge/                # Bridge to AI engines
│   │   └── aibridge.go
│   └── tui/                   # Bubble Tea TUI
│       └── app.go
├── plugin/                    # JUCE VST3/CLAP plugin (C++17)
│   ├── CMakeLists.txt
│   ├── PluginProcessor.cpp
│   ├── PluginProcessor.h
│   ├── PluginEditor.cpp
│   ├── PluginEditor.h
│   ├── HubClient.cpp          # WebSocket client → Go server
│   ├── HubClient.h
│   ├── Midi2Adapter.cpp       # ni-midi2 + libremidi I/O
│   ├── Midi2Adapter.h
│   └── LinkBridge.cpp         # Ableton Link peer (optional)
├── engines/                   # AI microservices (Python/Podman)
│   ├── audio-to-midi/
│   │   ├── main.py
│   │   ├── requirements.txt
│   │   └── Containerfile
│   ├── pattern-gen/
│   │   ├── main.py
│   │   ├── requirements.txt
│   │   └── Containerfile
│   └── emotional-engine/
│       ├── main.py
│       ├── requirements.txt
│       └── Containerfile
├── deps/                      # Git submodules
│   ├── ni-midi2/              # NI MIDI 2.0 C++ lib
│   ├── libremidi/             # Cross-platform MIDI I/O
│   ├── link/                  # Ableton Link
│   └── JUCE/                  # JUCE framework
├── CMakeLists.txt             # Root CMake (plugin + core)
├── compose.yaml               # Podman Compose for AI engines
├── .gitmodules                # All deps as submodules
└── README.md
```

---

## Message Protocol (WebSocket + JSON)

All messages between clients and server:

```json
{
  "type": "clock|clip|control|join|leave",
  "room": "room-id",
  "client": "client-id",
  "ts": 1710000000000,
  "payload": {}
}
```

### Examples

```json
// Clock consensus
{ "type": "clock", "payload": { "bpm": 138.0, "beat": 4, "phase": 0.25 } }

// MIDI 2.0 Clip (JSON MVP → SMF2 target)
{ "type": "clip", "payload": { "slot": "drums", "notes": [], "resolution": 32 } }

// Emotional macro (32-bit resolution)
{ "type": "control", "payload": { "macro": "tension", "value": 0.87 } }

// Join room
{ "type": "join", "payload": { "room": "techno-bcn", "profile": "deck" } }
```

---

## Roadmap

### Sprint 1 — Go server + TUI (1-2 weeks)
- [ ] WebSocket server with rooms and members
- [ ] Clock broadcast (host → all, consensus later)
- [ ] Bubble Tea TUI: room list, clients, BPM
- [ ] Test: 2 CLI instances sync over internet

### Sprint 2 — Clock consensus + clips (1-2 weeks)
- [ ] Distributed tempo consensus (Link-style logic)
- [ ] Send/receive JSON clips between clients
- [ ] TUI: show clips by slot, manual BPM change

### Sprint 3 — JUCE plugin (2-3 weeks)
- [ ] JUCE client connecting to Hub via WebSocket
- [ ] Receives clock → syncs DAW transport
- [ ] Sends/receives clip per slot (drums first)
- [ ] Minimal UI: connect, room, slot selector

### Sprint 4 — First AI engine + SMF2 (1-2 weeks)
- [ ] Python container with `basic-pitch` endpoint
- [ ] Hub calls it when client sends audio
- [ ] Returns JSON clip → serialize to SMF2
- [ ] Migrate clip format from JSON to SMF2

---

## Dependencies (Git Submodules)

| Submodule | Repo | Role |
|-----------|------|------|
| `deps/ni-midi2` | [midi2-dev/ni-midi2](https://github.com/midi2-dev/ni-midi2) | UMP parser, MIDI-CI, C++17 |
| `deps/libremidi` | [celtera/libremidi](https://github.com/celtera/libremidi) | Cross-platform MIDI I/O |
| `deps/link` | [Ableton/link](https://github.com/Ableton/link) | Ableton Link bridge |
| `deps/JUCE` | [juce-framework/JUCE](https://github.com/juce-framework/JUCE) | VST3/CLAP plugin framework |

```bash
git clone --recurse-submodules https://github.com/Phoenixai36/midi2-hub
```

---

## Quick Start (Sprint 1 — Server only)

```bash
# Clone with submodules
git clone --recurse-submodules https://github.com/Phoenixai36/midi2-hub
cd midi2-hub/server

# Run server
go mod tidy
go run main.go

# TUI launches automatically
# Connect clients via WebSocket ws://localhost:8080/ws
```

### AI Engines (Podman)

```bash
cd midi2-hub
podman compose -f compose.yaml up
# Engines available at:
# audio-to-midi:    http://localhost:8001
# pattern-gen:      http://localhost:8002
# emotional-engine: http://localhost:8003
```

### Plugin Build (JUCE)

```bash
cd midi2-hub
cmake -B build -DCMAKE_BUILD_TYPE=Release
cmake --build build
# VST3 output: build/plugin/midi2-hub_artefacts/VST3/
```

---

## Go Dependencies

```bash
go get nhooyr.io/websocket
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/google/uuid
go get golang.org/x/crypto
```

---

## Emotional Engine — Concept

Beyond technical sync, `midi2-hub` includes an **Emotional Engine** layer:

- Abstract macros: `tension`, `density`, `brightness`, `complexity`, `release`
- Mapped to **32-bit resolution MIDI 2.0 CC** (per-note controllers)
- AI engines condition pattern generation on emotional state
- Use case: live sets with emotional narrative arcs (R3B0RN project)

```
Tension ──────────────────────── 0.87
Density ────────────────── 0.65
Brightness ──────── 0.32
Complexity ─────────────────────── 0.91
Release ──── 0.18
```

---

## B2B Jam Mode

Two producers in different cities share:
- Clock/transport (Network MIDI 2.0)
- Scene triggers and clip slots
- An AI "third deck" that listens to both, proposes mashups and fills
- Control via emotional macros from each side

---

## Contributing

This is an open project. If you work with MIDI 2.0, AI music tools, or live performance tech — PRs and issues welcome.

---

## License

MIT — see [LICENSE](LICENSE)

---

*Built in Barcelona. Designed for everywhere.*
