// Felipe Ryan
// This file is mostly an expansion/customisation of Hajime Hoshi's original to be found here:
// https://github.com/hajimehoshi/ebiten/blob/master/examples/typewriter/main.go
// it was originally copyrighted under the Apache v2.0 license.
// This is an attempt to give credit where credit is due and act in good faith!

package uicomponents

import (
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type InputText struct {
	textBox    *ebiten.Image
	mouseDown  bool
	hasFocus   bool
	text       string
	showText   string
	fontFace   font.Face
	PosX       int
	PosY       int
	Width      int
	Height     int
	onPressed  func(b *InputText)
	onEnter    func(b *InputText)
	counter    int
	editable   bool
	charLength int
}

func createInputText(xPos, yPos, width, height int, sampleText string, numChars int) *InputText {
	itext := InputText{}
	itext.text = sampleText
	itext.showText = sampleText
	itext.PosX = xPos
	itext.PosY = yPos
	itext.Width = width
	itext.Height = height
	itext.counter = 0
	itext.editable = true
	itext.charLength = numChars

	itext.textBox, _ = ebiten.NewImage(width, height, ebiten.FilterDefault)
	itext.textBox.Fill(color.RGBA{0x80, 0x0, 0x80, 0x80})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xPos), float64(yPos))

	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	itext.fontFace = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return &itext
}

func (ipt *InputText) Draw(dst *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()

	if ipt.hasFocus {
		ipt.textBox.Fill(color.RGBA{0x0, 0x80, 0x80, 0x80})
	} else {
		ipt.textBox.Fill(color.RGBA{0x80, 0x80, 0x0, 0x80})
	}

	op.GeoM.Translate(float64(ipt.PosX), float64(ipt.PosY))
	dst.DrawImage(ipt.textBox, op)

	text.Draw(dst, ipt.showText, ipt.fontFace, ipt.PosX+5, ipt.PosY+30, color.White)

}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (b *InputText) Update() {
	if b.editable {

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if b.PosX <= x && x < (b.PosX+b.Width) && b.PosY <= y && y < (b.PosY+b.Height) {
				b.mouseDown = true
			} else {
				b.mouseDown = false
				b.hasFocus = false
				b.showText = b.text
			}
		} else {
			if b.mouseDown {
				b.hasFocus = true
				if b.onPressed != nil { // if we have a func to call then call it
					b.onPressed(b)
				}
			}
			b.mouseDown = false
		}

		if b.hasFocus {
			if len(b.text) < b.charLength {
				b.text += string(ebiten.InputChars())
			}
			if repeatingKeyPressed(ebiten.KeyBackspace) {
				if len(b.text) >= 1 {
					b.text = b.text[:len(b.text)-1]
				}
			}

			// Blink the cursor.
			b.counter++
			b.showText = b.text
			if b.counter%60 < 30 {
				b.showText += "_"
			}
			if b.counter == 601 {
				b.counter = 0
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if b.hasFocus {
				b.onEnter(b)
			}
		}
	} else {
		b.hasFocus = false
	}
}

func (b *InputText) SetOnPressed(f func(b *InputText)) {
	b.onPressed = f
}

func (b *InputText) SetOnEnter(f func(b *InputText)) {
	b.onEnter = f
}
