// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/titlescene.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
)

var incomingMessage = insertLineBreaksIntoLongString("Await instructions from Dungeon Master", 20)
var mainSceneAudioPlayer *audio.Player

type MainScene struct {
}

func (s *MainScene) Update(state *GameState) error {
	// Display incoming messages:
	select {

	case incoming := <-state.SceneManager.Incoming:
		incomingMessage = insertLineBreaksIntoLongString(fmt.Sprintf("%s", incoming), 20)

	case sounds := <-state.SceneManager.AudioIn:
		dec, err := mp3.Decode(audioContext, audio.BytesReadSeekCloser(sounds.([]byte)))
		if err != nil {
			return err
		}
		mainSceneAudioPlayer, err = audio.NewPlayer(audioContext, dec)
		if err != nil {
			log.Fatal(err)
		}
		mainSceneAudioPlayer.Play()
	default:
		return nil
	}

	return nil
}

func (s *MainScene) Draw(screen *ebiten.Image) {

	drawTextWithShadowCenter(screen, incomingMessage, 20, 100, 2, color.NRGBA{0x99, 0x66, 0xcc, 0xff}, ScreenWidth)
}

func NewMainScene() *MainScene {
	return &MainScene{}
}

func insertLineBreaksIntoLongString(text string, lineBreakAfterNChars int) string {
	rs := []rune(text)
	newrs := make([]string, 0)

	flag := false
	count := 0

	for i, r := range rs {
		if i > 0 && i%lineBreakAfterNChars == 0 {
			flag = true
		}
		if flag && r == []rune(" ")[0] {
			newrs = append(newrs[:i+count], append([]string{"\n"}, newrs[i+count:]...)...)
			flag = false
			count++
		}
		newrs = append(newrs, string(r))
	}

	ss := strings.Join(newrs, "")
	ss = strings.Replace(ss, "\n ", "\n", -1)
	return ss
}
