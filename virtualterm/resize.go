package virtualterm

import (
	"github.com/lmorg/mxtty/types"
)

func (term *Term) Resize(size *types.XY) {

}

func (term *Term) resize80() {
	term.renderer.AddRenderFnToStack(func() {
		term.renderer.ResizeWindow(&types.XY{X: 80, Y: 25})
		term.resize80()
	})
}

func (term *Term) resize132() {
	term.renderer.AddRenderFnToStack(func() {
		term.renderer.ResizeWindow(&types.XY{X: 132, Y: 25})
		term.resize132()
	})
}
