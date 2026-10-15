package rendererwebkit

import "github.com/lmorg/ttyphoon/types"

func (wr *webkitRender) DrawFrame(tile types.Tile) {
	var flags int32
	if pane, ok := tile.(interface{ InteractiveActionChrome() *types.Colour }); ok {
		switch pane.InteractiveActionChrome() {
		case types.COLOR_OK:
			flags = 1
		case types.COLOR_ERROR:
			flags = 2
		}
	}
	wr.enqueueDrawCommand(DrawCommand{
		Op:     DrawOpFrame,
		Flags:  flags,
		X:      tile.Left(),
		Y:      tile.Top() - 1,
		Width:  tile.Right() - tile.Left() + 2,
		Height: tile.Bottom() - tile.Top() + 2,
	})
}
