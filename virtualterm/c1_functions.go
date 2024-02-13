package virtualterm

import "github.com/lmorg/mxtty/types"

func (term *Term) c1DecalnTestAlignment() {
	term.curPos = types.XY{} // top left
	for i := int32(0); i < term.size.X*term.size.Y; i++ {
		term.writeCell('E')
	}
	term.curPos = types.XY{} // top left
}
