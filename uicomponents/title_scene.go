// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/blocks/blocks/titlescene.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var imageBackground *ebiten.Image
var logoMove int
var subtitleMove int
var spacebarCallout int

var audioContext *audio.Context
var audioPlayer *audio.Player

const logoMoveMax int = 64

func init() {
	// img, _, err := ebitenutil.NewImageFromFile("smaug_par_david_demaret.jpg", ebiten.FilterDefault)
	img, _, err := image.Decode(bytes.NewReader(backgroundJpg))
	if err != nil {
		panic(err)
	}
	imageBackground, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	audioContext, err = audio.NewContext(22050)
	if err != nil {
		log.Fatal(err)
	}
}

type TitleScene struct {
	count int
}

func (s *TitleScene) Update(state *GameState) error {
	if audioPlayer == nil {
		// f, err := ioutil.ReadFile("bensound-epic.mp3")
		// if err != nil {
		// 	return err
		// }
		dec, err := mp3.Decode(audioContext, audio.BytesReadSeekCloser(titleMusic))
		if err != nil {
			return err
		}
		audioPlayer, err = audio.NewPlayer(audioContext, dec)
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.SetVolume(0.6)
	}
	if !audioPlayer.IsPlaying() {
		audioPlayer.Play()
	}

	s.count++

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println("Spacebar pressed")
		if audioPlayer.IsPlaying() {
			audioPlayer.Close()
		}
		state.SceneManager.GoTo(NewSettingsScene())
		return nil
	}
	return nil
}

func (s *TitleScene) Draw(r *ebiten.Image) {
	s.drawTitleBackground(r, s.count)
	drawLogo(r, "DUNGEONS & DRAGONS")

	if logoMove >= 64 {
		drawSubtitle(r, "Jogatina de Sao Nunca")
	}
	if subtitleMove >= 112 {
		drawSpacebarCallout(r, "PRESS SPACEBAR TO START")
	}
}

func (s *TitleScene) drawTitleBackground(r *ebiten.Image, c int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	r.DrawImage(imageBackground, op)
}

func drawLogo(r *ebiten.Image, str string) {
	const scale = 4
	x := 0
	y := 64

	if logoMove < logoMoveMax {
		drawTextWithShadowCenter(r, str, x, logoMove, scale, color.NRGBA{0x00, 0x00, 0x00, 0xff}, ScreenWidth)
		logoMove++
		return
	}
	drawTextWithShadowCenter(r, str, x, y, scale, color.NRGBA{0x00, 0x00, 0x00, 0xff}, ScreenWidth)
}

func drawSubtitle(r *ebiten.Image, str string) {
	const scale = 2
	x := 0
	y := 112

	if subtitleMove < y {
		drawTextWithShadowCenter(r, str, x, subtitleMove, scale, color.NRGBA{0x00, 0x00, 0x00, 0xff}, ScreenWidth)
		subtitleMove++
		return
	}
	drawTextWithShadowCenter(r, str, x, y, scale, color.NRGBA{0x00, 0x00, 0x00, 0xff}, ScreenWidth)
}

func drawSpacebarCallout(r *ebiten.Image, str string) {
	const scale = 2
	x := 0
	y := 112

	if spacebarCallout < y {
		drawTextWithShadowCenter(r, str, x, ScreenHeight-spacebarCallout, scale, color.NRGBA{0x99, 0xCC, 0x99, 0xff}, ScreenWidth)
		spacebarCallout++
		return
	}
	drawTextWithShadowCenter(r, str, x, ScreenHeight-y, scale,
		color.NRGBA{0x99, 0xCC, 0x99, 0xff}, ScreenWidth)
}
