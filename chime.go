package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

var player *oto.Context

//go:embed chime.mp3
var chimeData []byte

func setupChime() error {
	d, _ := mp3.NewDecoder(bytes.NewReader(chimeData))
	ctx, _, err := oto.NewContext(d.SampleRate(), 2, 2)
	if err != nil {
		return err
	}
	player = ctx
	//	<-ready

	return nil
}

func chime() {
	d, _ := mp3.NewDecoder(bytes.NewReader(chimeData))
	ding := player.NewPlayer(d)

	ding.Play()

	go func() {
		for ding.IsPlaying() {
		}

		ding.Close()
	}()
}
