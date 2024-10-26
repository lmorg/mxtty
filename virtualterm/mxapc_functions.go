package virtualterm

import (
	"log"

	"github.com/lmorg/murex/utils/json"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) mxapcBegin(element types.ElementID, parameters *types.ApcSlice) {
	term._activeElement = term.renderer.NewElement(element, nil, nil)
	term._activeElement.Begin(parameters)
}

func (term *Term) mxapcEnd() {
	if term._activeElement == nil {
		return
	}
	term._activeElement.End()
	term._activeElement = nil
}

func (term *Term) mxapcInsert(element types.ElementID, parameters *types.ApcSlice) {
	// this is ugly and needs rewriting
	term._activeElement = term.renderer.NewElement(element, nil, nil)

	size := term._activeElement.Insert(parameters)

	if size != nil {
		term.currentCell().Element = term._activeElement

		for i := int32(1); i < size.Y; i++ {
			term.lineFeed()
		}

		start := &types.XY{X: term.curPos().X, Y: term.curPos().Y - size.Y}
		term._elementResizeGrow(term._activeElement, start, size)
	}

	term._activeElement = nil
}

func (term *Term) _elementResizeGrow(el types.Element, start *types.XY, size *types.XY) {
	log.Printf("DEBUG: _elementResizeUpdate(): %s | %s", json.LazyLogging(start), json.LazyLogging(size))
	for y := start.Y; y < start.Y+size.Y && y < int32(len(*term.cells)); y++ {

		for x := start.X; x < start.X+size.X && x < int32(len((*term.cells)[y])); x++ {
			(*term.cells)[y][x].Element = el
		}
	}
}
