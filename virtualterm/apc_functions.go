package virtualterm

import (
	"github.com/lmorg/mxtty/types"
)

const (
	cellElementXyMask    = (^int32(0)) << 16
	cellElementXyCeiling = int32(^uint16(0) >> 1)
)

func setElementXY(xy *types.XY) rune {
	if xy.X > cellElementXyCeiling || xy.Y > cellElementXyCeiling {
		panic("TODO: proper error handling")
	}
	return (xy.X << 16) | xy.Y
}

func getElementXY(r rune) *types.XY {
	return &types.XY{
		X: r >> 16,
		Y: r &^ cellElementXyMask,
	}
}

func (term *Term) mxapcBegin(element types.ElementID, parameters *types.ApcSlice) {
	term._activeElement = term.renderer.NewElement(element)
}

func (term *Term) mxapcEnd(_ types.ElementID, parameters *types.ApcSlice) {
	if term._activeElement == nil {
		return
	}
	el := term._activeElement           // this needs to be in this order because a
	term._activeElement = nil           // function inside _mxapcGenerate returns
	term._mxapcGenerate(el, parameters) // without processing if _activeElement set
}

func (term *Term) mxapcInsert(element types.ElementID, parameters *types.ApcSlice) {
	term._mxapcGenerate(term.renderer.NewElement(element), parameters)
}

func (term *Term) _mxapcGenerate(el types.Element, parameters *types.ApcSlice) {
	err := el.Generate(parameters)
	if err != nil {
		term.renderer.DisplayNotification(types.NOTIFY_ERROR, err.Error())
		return
	}

	size := el.Size()
	lineWrap := term._noAutoLineWrap
	term._noAutoLineWrap = true

	elPos := new(types.XY)
	for ; elPos.Y < size.Y; elPos.Y++ {
		if term.curPos().X != 0 {
			term.carriageReturn()
			term.lineFeed()
		}
		for elPos.X = 0; elPos.X < size.X && term._curPos.X < term.size.X; elPos.X++ {
			term.writeCell(setElementXY(elPos), el)
		}
	}

	term._noAutoLineWrap = lineWrap
}
