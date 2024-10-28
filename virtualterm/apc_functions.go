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
	/*el := term.renderer.NewElement(element)
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

	term._noAutoLineWrap = lineWrap*/
}

func (term *Term) mxapcEnd(element types.ElementID, parameters *types.ApcSlice) {
	/*el := term.renderer.NewElement(element)
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

	term._noAutoLineWrap = lineWrap*/
}

func (term *Term) mxapcInsert(element types.ElementID, parameters *types.ApcSlice) {
	el := term.renderer.NewElement(element)
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
