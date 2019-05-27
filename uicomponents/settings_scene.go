// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/titlescene.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type SettingsScene struct {
	serverAddressBox *InputText
	nameBox          *InputText
	incomingText     string
}

func setCallBacksForInputText(s *SettingsScene, state *GameState) {
	// Set textfield function for pressing enter, only need to set it once
	if s.serverAddressBox.onEnter == nil {
		s.serverAddressBox.SetOnEnter(func(b *InputText) {
			s.serverAddressBox.showText = s.serverAddressBox.text
			fmt.Println("server address enter")
			s.serverAddressBox.editable = false
			state.SceneManager.Outgoing <- fmt.Sprintf("%s", s.serverAddressBox.text)
			s.nameBox.editable = true
		})
	}

	if s.nameBox.onEnter == nil {
		s.nameBox.SetOnEnter(func(b *InputText) {
			s.nameBox.showText = s.nameBox.text
			fmt.Println("name box Enter")
			s.nameBox.editable = false
			state.SceneManager.Outgoing <- fmt.Sprintf("/setname %s", s.nameBox.text)

			transitionToMain(state) // Go to main scene
		})
	}

}

func transitionToMain(state *GameState) {
	state.SceneManager.GoTo(NewMainScene())
}

func (s *SettingsScene) Update(state *GameState) error {
	s.serverAddressBox.Update()
	s.nameBox.Update()
	setCallBacksForInputText(s, state)

	// Display incoming messages:
	select {

	case incoming := <-state.SceneManager.Incoming:
		incomingMessage = insertLineBreaksIntoLongString(fmt.Sprintf("%s", incoming), 20)
	default:
		return nil
	}

	return nil
}

func (s *SettingsScene) Draw(screen *ebiten.Image) {
	s.serverAddressBox.Draw(screen)
	s.nameBox.Draw(screen)

	drawTextWithShadowCenter(screen, "Enter server address + Enter", 20, 20, 2, color.NRGBA{0x99, 0x66, 0xcc, 0xff}, ScreenWidth)
	drawTextWithShadowCenter(screen, "Enter your name + Enter", 20, 180, 2, color.NRGBA{0x99, 0x66, 0xcc, 0xff}, ScreenWidth)
	drawTextWithShadowCenter(screen, incomingMessage, 20, 350, 2, color.NRGBA{0x99, 0x66, 0xcc, 0xff}, ScreenWidth)

}

func NewSettingsScene() *SettingsScene {
	setScene := &SettingsScene{}

	boxWidth := 300
	horizontalPos := (ScreenWidth - boxWidth) / 2
	setScene.serverAddressBox = createInputText(horizontalPos, 60, boxWidth, 48, "go-dand-server.herokuapp.com", 36)
	setScene.serverAddressBox.hasFocus = true

	setScene.nameBox = createInputText(horizontalPos, 220, boxWidth, 48, "", 36)
	setScene.nameBox.editable = false

	return setScene
}
