// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/game.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"github.com/hajimehoshi/ebiten"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
)

type Game struct {
	sceneManager *SceneManager
	OutgoingMSG  chan interface{}
	IncomingMSG  chan interface{}
	AudioInMSG   chan interface{}
}

func (g *Game) Update(r *ebiten.Image) error {
	if g.sceneManager == nil {
		g.sceneManager = &SceneManager{}
		g.sceneManager.Outgoing = g.OutgoingMSG
		g.sceneManager.Incoming = g.IncomingMSG
		g.sceneManager.AudioIn = g.AudioInMSG
		g.sceneManager.GoTo(&TitleScene{})
	}

	if err := g.sceneManager.Update(); err != nil {
		return err
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	g.sceneManager.Draw(r)
	return nil
}
