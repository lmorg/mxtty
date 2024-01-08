package rendererimgui

import (
	imgui "github.com/AllenDang/cimgui-go"
	"github.com/lmorg/mxtty/types"
)

func (renderer *imguiRender) PrintRuneColour(r rune, posX, posY int32, fg *types.Colour, bg *types.Colour, style types.SgrFlag) error {
	imgui.Begin("bob")
	imgui.TextColored(colour(fg), string(r))
	imgui.End()
	return nil
}
