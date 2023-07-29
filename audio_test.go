// audio_test.go dosyası
package main

import (
	"testing"
	"time"
)

func TestBuyPlayerAudio(t *testing.T) {
	err := InitPlayers()
	if err != nil {
		t.Fatalf("Error initializing audio players: %v", err)
	}
	defer ClosePlayers()

	// Test buy sound
	ResumeBuyPlayer()
	time.Sleep(5 * time.Second) // Wait for 5 seconds to hear the sound
	PausePlayers()
	time.Sleep(1 * time.Second) // Wait for 1 second
}
