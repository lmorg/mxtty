package virtualterm

import "fmt"

func (term *Term) drawElement(cell *cell) error {
	//e :=cell.element
	switch cell.element.(type) {
	case nil:
		return fmt.Errorf("nil pointer to element")

	default:
		return fmt.Errorf("unknown element type")

	}

	//e.Draw(nil) // TODO: this shouldn't be nil
	return nil
}
