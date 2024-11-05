package elementImage

import (
	"errors"
	"fmt"

	"github.com/lmorg/mxtty/types"
)

type ElementImage struct {
	renderer   types.Renderer
	parameters parametersT
	size       *types.XY
	load       func([]byte, *types.XY) (types.Image, error)
	bmp        []byte
	image      types.Image
}

type parametersT struct {
	Base64   string
	Filename string
	Width    int32
	Height   int32
}

func New(renderer types.Renderer, loadFn func([]byte, *types.XY) (types.Image, error)) *ElementImage {
	return &ElementImage{renderer: renderer, load: loadFn}
}

func (el *ElementImage) Generate(apc *types.ApcSlice) error {
	notify := el.renderer.DisplaySticky(types.NOTIFY_DEBUG, "Importing image from ANSI escape codes....")
	defer notify.Close()

	apc.Parameters(&el.parameters)

	el.size = new(types.XY)
	el.size.X, el.size.Y = el.parameters.Width, el.parameters.Height

	if el.size.X == 0 && el.size.Y == 0 {
		el.size.Y = 15 // default
	}

	err := el.decode()
	if err != nil {
		return fmt.Errorf("cannot decode image: %s", err.Error())
	}

	// cache image

	el.image, err = el.load(el.bmp, el.size)
	if err != nil {
		return fmt.Errorf("cannot cache image: %s", err.Error())
	}
	return nil
}

func (el *ElementImage) Write(_ rune) error {
	return errors.New("not supported")
}

func (el *ElementImage) Size() *types.XY {
	return el.size
}

// Draw:
// size: optional. Defaults to element size
// pos:  required. Position to draw element
func (el *ElementImage) Draw(size *types.XY, pos *types.XY) {
	if len(el.bmp) == 0 {
		return
	}

	if size == nil {
		size = el.size // TODO: eh???
	}

	el.renderer.AddRenderFnToStack(func() {
		el.image.Draw(el.size, pos)
	})
}

func (el *ElementImage) Close() {
	// clear memory (if required)
	el.image.Close()
}

func (el *ElementImage) MouseClick(button uint8, _ *types.XY, callback types.MouseClickCallback) {
	if button != 1 {
		callback()
		return
	}

	err := el.fullscreen()
	if err != nil {
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Unable to go fullscreen: "+err.Error())
	}
}

func (el *ElementImage) Rune(_ *types.XY) rune {
	return ' '
}
