package virtualterm

import (
	"fmt"

	"github.com/lmorg/mxtty/virtualterm/cell"
)

func (term *Term) drawElement(cell *cell.Cell) error {
	//e :=cell.element
	switch cell.Element.(type) {
	case nil:
		return fmt.Errorf("nil pointer to element")

	default:
		return fmt.Errorf("unknown element type")

	}

	//e.Draw(nil) // TODO: this shouldn't be nil
	return nil
}
