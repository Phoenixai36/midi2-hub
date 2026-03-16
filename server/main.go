package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Phoenixai36/midi2-hub/server/session"
	"github.com/Phoenixai36/midi2-hub/server/timesync"
	"github.com/Phoenixai36/midi2-hub/server/transport"
	"github.com/Phoenixai36/midi2-hub/server/tui"
)

const defaultAddr = ":8080"

func main() {
	addr := os.Getenv("HUB_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	// Core components
	clock := timesync.NewClock(120.0) // default 120 BPM
	manager := session.NewManager(clock)
	server := transport.NewServer(addr, manager)

	// Start TUI in goroutine
	go func() {
		if err := tui.Run(manager, clock); err != nil {
			log.Printf("TUI error: %v", err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	go func() {
		fmt.Printf("midi2-hub server listening on %s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\nShutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(shutCtx)
}
