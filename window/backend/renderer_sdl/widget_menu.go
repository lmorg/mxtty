package rendersdl

import (
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

/*var (
	questionColor       = &types.Colour{0x00, 0x77, 0x00}
	questionColorBorder = &types.Colour{0x00, 0xff, 0x00}
)*/

type menuCallbackT func(int)

type menuT struct {
	Items    []string
	Callback menuCallbackT
}

func (sr *sdlRender) DisplayMenu(message string, defaultValue string, callback func(string)) {

}

func (sr *sdlRender) closeMenu() {

}

func (menu *menuT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	// do nothing
}

func (menu *menuT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	// do nothing
}

func (menu *menuT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	// do nothing
}

func (menu *menuT) eventMouseWheel(sr *sdlRender, evt *sdl.MouseWheelEvent) {
	// do nothing
}

func (menu *menuT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	// do nothing
}

func (sr *sdlRender) renderMenu(windowRect *sdl.Rect) {
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, windowRect.W, windowRect.H, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) //TODO: don't panic!
	}
	defer surface.Free()

	sr.setFontStyle(types.SGR_BOLD)

	text, err := sr.font.RenderUTF8BlendedWrapped(sr.inputBox.Message, sdl.Color{R: 200, G: 200, B: 200, A: 255}, int(surface.W-sr.notifyIconSize.X))
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	defer text.Free()

	textShadow, err := sr.font.RenderUTF8BlendedWrapped(sr.inputBox.Message, sdl.Color{R: 0, G: 0, B: 0, A: 150}, int(surface.W-sr.notifyIconSize.X))
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	defer textShadow.Free()

	/*
		FRAME
	*/

	height := text.H + (sr.border * 5) + sr.glyphSize.Y
	offset := (surface.H / 2) - (height / 2)
	padding := sr.border * 2

	// draw border
	sr.renderer.SetDrawColor(questionColorBorder.Red, questionColorBorder.Green, questionColorBorder.Blue, notificationAlpha)
	rect := sdl.Rect{
		X: sr.border - 1,
		Y: offset - 1,
		W: windowRect.W - padding + 2,
		H: height + 2,
	}
	sr.renderer.DrawRect(&rect)
	rect = sdl.Rect{
		X: sr.border,
		Y: offset,
		W: windowRect.W - padding,
		H: height,
	}
	sr.renderer.DrawRect(&rect)

	// fill background
	sr.renderer.SetDrawColor(questionColor.Red, questionColor.Green, questionColor.Blue, notificationAlpha)
	rect = sdl.Rect{
		X: sr.border + 1,
		Y: 1 + offset,
		W: surface.W - padding - 2,
		H: height - 2,
	}
	sr.renderer.FillRect(&rect)

	// render shadow
	rect = sdl.Rect{
		X: padding + sr.notifyIconSize.X + 2,
		Y: sr.border + offset + 2,
		W: surface.W - sr.notifyIconSize.X,
		H: text.H + padding - 2,
	}
	_ = textShadow.Blit(nil, surface, &rect)
	sr._renderNotificationSurface(surface, &rect)

	// render text
	rect = sdl.Rect{
		X: padding + sr.notifyIconSize.X,
		Y: sr.border + offset,
		W: surface.W - sr.notifyIconSize.X,
		H: text.H + padding - 2,
	}
	err = text.Blit(nil, surface, &rect)
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	sr._renderNotificationSurface(surface, &rect)

	/*
		TEXT FIELD
	*/

	height = sr.glyphSize.Y + (sr.border * 2)
	offset += text.H + sr.border + sr.border
	var width int32

	// draw border
	sr.renderer.SetDrawColor(255, 255, 255, 150)
	rect = sdl.Rect{
		X: sr.notifyIconSize.X + padding - 1,
		Y: offset - 1,
		W: windowRect.W - sr.notifyIconSize.X - padding - padding + 2,
		H: height + 2,
	}
	sr.renderer.DrawRect(&rect)
	rect = sdl.Rect{
		X: sr.notifyIconSize.X + padding,
		Y: offset,
		W: windowRect.W - sr.notifyIconSize.X - padding - padding,
		H: height,
	}
	sr.renderer.DrawRect(&rect)

	// fill background
	sr.renderer.SetDrawColor(0, 0, 0, 150)
	rect = sdl.Rect{
		X: sr.notifyIconSize.X + padding + 1,
		Y: 1 + offset,
		W: surface.W - sr.notifyIconSize.X - padding - padding - 2,
		H: height - 2,
	}
	sr.renderer.FillRect(&rect)

	// value
	if len(sr.inputBox.Value) > 0 {
		textValue, err := sr.font.RenderUTF8Blended(sr.inputBox.Value, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer textValue.Free()

		rect = sdl.Rect{
			X: padding + sr.notifyIconSize.X + sr.border,
			Y: sr.border + offset,
			W: surface.W - sr.notifyIconSize.X,
			H: textValue.H + padding - 2,
		}
		err = textValue.Blit(nil, surface, &rect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		sr._renderNotificationSurface(surface, &rect)
		width = textValue.W
	}

	if surface, ok := sr.notifyIcon[types.NOTIFY_QUESTION].Asset().(*sdl.Surface); ok {
		srcRect := &sdl.Rect{
			X: 0,
			Y: 0,
			W: surface.W,
			H: surface.H,
		}

		dstRect := &sdl.Rect{
			X: sr.border + 2,
			Y: offset + text.H - sr.notifyIconSize.Y,
			W: sr.notifyIconSize.X,
			H: sr.notifyIconSize.X,
		}

		texture, err := sr.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err) //TODO: don't panic!
		}
		defer texture.Destroy()

		err = sr.renderer.Copy(texture, srcRect, dstRect)
		if err != nil {
			panic(err) //TODO: don't panic!
		}
	}

	if sr.blinkState {
		sr.renderer.SetDrawColor(255, 255, 255, 255)
		rect = sdl.Rect{
			X: padding + sr.notifyIconSize.X + sr.border + width,
			Y: sr.border + offset,
			W: sr.glyphSize.X,
			H: sr.glyphSize.Y,
		}
		sr.renderer.FillRect(&rect)
	}

	texture, err := sr.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	defer texture.Destroy()

	err = sr.renderer.Copy(texture, windowRect, windowRect)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
}
