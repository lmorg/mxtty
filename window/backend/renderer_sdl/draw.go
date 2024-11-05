package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/types"
)

func (sr *sdlRender) DrawTable(pos *types.XY, height int32, boundaries []int32) {
	sr.fnStack = append(sr.fnStack, func() {
		var err error

		tx :=sr.renderer.GetRenderTarget()
		tx.GetTextureUserData()

		sr.renderer.SetDrawColor(255, 255, 255, 64)

		X := (pos.X * sr.glyphSize.X) + sr.border
		Y := (pos.Y * sr.glyphSize.Y) + sr.border
		H := Y + ((height + 1) * sr.glyphSize.Y)

		err = sr.renderer.DrawLine(X, Y, X, H)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}

		for i := range boundaries {
			x := X + (boundaries[i] * sr.glyphSize.X)
			err = sr.renderer.DrawLine(x, Y, x, H)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
				return
			}
		}

		x := X + (boundaries[len(boundaries)-1] * sr.glyphSize.X)
		y := Y + ((height + 1) * sr.glyphSize.Y)
		err = sr.renderer.DrawLine(X, y, x, y)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}

		sr.renderer.SetDrawColor(255, 255, 255, 32)

		for i := int32(0); i <= height; i++ {
			y = Y + (i * sr.glyphSize.Y)
			err = sr.renderer.DrawLine(X, y, x, y)
			if err != nil {
				log.Printf("ERROR: %s", err.Error())
				return
			}
		}
	})
}
