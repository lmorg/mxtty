package elementImage

import (
	"strconv"

	"github.com/lmorg/mxtty/types"
)

const (
	_KEY_BASE64   = "base64"
	_KEY_FILENAME = "filename"
	_KEY_HEIGHT   = "height"
	_KEY_WIDTH    = "width"
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

func (el *ElementImage) Begin(_ *types.ApcSlice) {
	// not required for this element
}

func (el *ElementImage) ReadCell(cell *types.Cell) {
	// not required for this element
}

func (el *ElementImage) End() *types.XY {
	// not required for this element
	return nil
}

func (el *ElementImage) Insert(apc *types.ApcSlice) *types.XY {
	el.renderer.DisplayNotification(types.NOTIFY_DEBUG, "Importing image from ANSI escape codes....")

	el.apc = apc
	el.size = new(types.XY)

	width := apc.Parameter(_KEY_WIDTH)
	if width != "" {
		i, err := strconv.Atoi(width)
		if err != nil {
			el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot convert width: "+err.Error())
		}
		el.size.X = int32(i)
	}

	height := apc.Parameter(_KEY_HEIGHT)
	if height != "" {
		i, err := strconv.Atoi(height)
		if err != nil {
			el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot convert height: "+err.Error())
		}
		el.size.Y = int32(i)
	}

	if el.size.X == 0 && el.size.Y == 0 {
		el.size.Y = 15 // default
	}

	err := el.decode()
	if err != nil {
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot decode image: "+err.Error())
		return nil
	}

	return el.size
}

func (el *ElementImage) Draw(rect *types.Rect) *types.XY {
	if len(el.bmp) == 0 {
		return nil
	}

	var updateSize bool

	if el.image == nil {
		// cache image
		var err error

		el.image, err = el.load(el.bmp, el.size)
		if err != nil {
			el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Cannot cache image: "+err.Error())
			rect.End.X = rect.Start.X
			rect.End.Y = rect.Start.Y
			return &types.XY{}
		}

		updateSize = true

	}

	el.renderer.AddRenderFnToStack(func() {
		el.image.Draw(el.size, rect)
	})

	if updateSize {
		return el.size
	}
	return nil
}

func (el *ElementImage) Close() {
	// clear memory (if required)
	el.image.Close()
}

func (el *ElementImage) MouseClick(_ uint8, _ *types.XY) {
	//el.renderer.AddImageToStack(func() {
	err := el.fullscreen()
	if err != nil {
		el.renderer.DisplayNotification(types.NOTIFY_ERROR, "Unable to go fullscreen: "+err.Error())
	}
	//})
}
