package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

const buySoundFile = "buy.mp3"

var (
	buyStreamer beep.StreamSeekCloser
	format      beep.Format
	buyCtrl     *beep.Ctrl
)

func initPlayer(soundFile string) (beep.StreamSeekCloser, *beep.Ctrl, error) {
	f, err := os.Open(soundFile)
	if err != nil {
		return nil, nil, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		streamer.Close()
	})), Paused: true}
	speaker.Play(ctrl)

	return streamer, ctrl, nil
}

func InitPlayers() error {
	var err error

	buyStreamer, buyCtrl, err = initPlayer(buySoundFile)
	if err != nil {
		return err
	}

	return nil
}

func ClosePlayers() {
	buyStreamer.Close()
}

func PausePlayers() {
	buyCtrl.Paused = true
}

func ResumeBuyPlayer() {
	PausePlayers()
	buyCtrl.Paused = false
}
