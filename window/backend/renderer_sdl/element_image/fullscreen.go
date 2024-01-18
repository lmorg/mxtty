package elementImage

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func (el *ElementImage) fullscreen() error {
	var (
		fullscreen uint32
		err        error
		dispIndex  int
		dispMode   sdl.DisplayMode
	)

	win, ok := el.renderer.GetWindowMeta().(*sdl.Window)
	if !ok {
		log.Println("DEBUG: (el *ElementImage) fullscreen(): WindowMeta is not a window")
		fullscreen = sdl.WINDOW_FULLSCREEN_DESKTOP
		goto createWindow
	}
	dispIndex, err = win.GetDisplayIndex()
	if err != nil {
		log.Printf("DEBUG: (el *ElementImage) fullscreen(): %s", err)
		fullscreen = sdl.WINDOW_FULLSCREEN_DESKTOP
		goto createWindow
	}

	dispMode, err = sdl.GetDesktopDisplayMode(dispIndex)
	if err != nil {
		log.Printf("DEBUG: (el *ElementImage) fullscreen(): %s", err)
		fullscreen = sdl.WINDOW_FULLSCREEN_DESKTOP
		goto createWindow
	}

	fullscreen = sdl.WINDOW_FULLSCREEN
createWindow:

	window, err := sdl.CreateWindow(
		"mxtty",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, dispMode.W, dispMode.W,
		sdl.WINDOW_SHOWN|fullscreen|sdl.WINDOW_ALWAYS_ON_TOP,
	)
	if err != nil {
		return err
	}
	defer el.renderer.FocusWindow()
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}
	defer renderer.Destroy()

	imgSurface, ok := el.image.Asset().(*sdl.Surface)
	if !ok {
		return fmt.Errorf("image asset is not a surface")
	}

	texture, err := renderer.CreateTextureFromSurface(imgSurface)
	if err != nil {
		return err
	}

	imgRatio := float64(imgSurface.W) / float64(imgSurface.H)
	x, y = window.GetSize()
	winRatio := float64(x) / float64(y)

	x, y := imgSurface.W, imgSurface.H
	srcRect := &sdl.Rect{W: x, H: y}

	dstRect := &sdl.Rect{W: x, H: y}

	err = renderer.Copy(texture, srcRect, dstRect)
	if err != nil {
		return err
	}

	renderer.Present()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {

			case *sdl.QuitEvent:
				return nil

			case *sdl.KeyboardEvent:
				return nil

			case *sdl.MouseButtonEvent:
				if evt.State == sdl.PRESSED {
					continue
				}
				return nil
			}
		}

		sdl.Delay(15)
	}
}
