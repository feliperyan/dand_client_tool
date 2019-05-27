// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/scenemanager.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"github.com/hajimehoshi/ebiten"
)

var (
	transitionFrom *ebiten.Image
	transitionTo   *ebiten.Image
)

func init() {
	transitionFrom, _ = ebiten.NewImage(ScreenWidth, ScreenHeight, ebiten.FilterDefault)
	transitionTo, _ = ebiten.NewImage(ScreenWidth, ScreenHeight, ebiten.FilterDefault)
}

type Scene interface {
	Update(state *GameState) error
	Draw(screen *ebiten.Image)
}

const transitionMaxCount = 20 // to do with next scene fading in

type SceneManager struct {
	current         Scene
	next            Scene
	transitionCount int
	Outgoing        chan interface{}
	Incoming        chan interface{}
	AudioIn         chan interface{}
}

type GameState struct {
	SceneManager *SceneManager
}

func (s *SceneManager) Update() error {
	if s.transitionCount == 0 {
		return s.current.Update(&GameState{
			SceneManager: s,
		})
	}

	s.transitionCount--
	if s.transitionCount > 0 {
		return nil
	}

	s.current = s.next
	s.next = nil
	return nil
}

func (s *SceneManager) Draw(r *ebiten.Image) {
	if s.transitionCount == 0 {
		s.current.Draw(r)
		return
	}

	transitionFrom.Clear()
	s.current.Draw(transitionFrom)

	transitionTo.Clear()
	s.next.Draw(transitionTo)

	r.DrawImage(transitionFrom, nil)

	alpha := 1 - float64(s.transitionCount)/float64(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, alpha)
	r.DrawImage(transitionTo, op)
}

func (s *SceneManager) GoTo(scene Scene) {
	if s.current == nil {
		s.current = scene
	} else {
		s.next = scene
		s.transitionCount = transitionMaxCount
	}
}
