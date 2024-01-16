package elementImage

import (
	"log"
	"strconv"

	"github.com/lmorg/mxtty/types"
)

const (
	_KEY_BASE64 = "base64"
	_KEY_WIDTH  = "width"
)

type ElementImage struct {
	renderer types.Renderer
	size     *types.XY
	apc      *types.ApcSlice
	load     func([]byte, *types.XY) (types.Image, error)
	bmp      []byte
	image    types.Image
}

func New(renderer types.Renderer, loadFn func([]byte, *types.XY) (types.Image, error)) *ElementImage {
	return &ElementImage{renderer: renderer, load: loadFn}
}

func (el *ElementImage) Begin(apc *types.ApcSlice) {
	el.apc = apc
	el.size = new(types.XY)

	width := apc.Parameter(_KEY_WIDTH)
	if width != "" {
		i, err := strconv.Atoi(width)
		if err != nil {
			log.Printf("ERROR: cannot convert width: %s", err.Error())
		}
		el.size.X = int32(i)
	}
}

func (el *ElementImage) ReadCell(cell *types.Cell) {
	switch cell {
	case nil:
		el.size.Y++
	}
}

func (el *ElementImage) End() {
	err := el.decode()
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}
}

func (el *ElementImage) Draw(rect *types.Rect) {
	//log.Printf("DEBUG: image.Draw")

	if len(el.bmp) == 0 {
		return
	}

	if el.image == nil {
		// cache image
		var err error
		el.image, err = el.load(el.bmp, el.size)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			return
		}
	}

	el.image.Draw(el.size, rect)
}

func (el *ElementImage) Close() {
	// clear memory (if required)
	el.image.Close()
}

func (el *ElementImage) MouseClick(_ uint8, _ *types.XY) {
	// do nothing
}
