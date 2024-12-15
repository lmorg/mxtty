package rendersdl

import (
	"strings"
	"time"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/renderer_sdl/layer"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const MENU_SEPARATOR = "-"

type menuT struct {
	options           []string
	title             string
	highlightIndex    int
	highlightCallback types.MenuCallbackT
	selectCallback    types.MenuCallbackT
	cancelCallback    types.MenuCallbackT
	mouseRect         sdl.Rect
	maxLen            int32
	filter            string
	hidden            []bool
	blinkState        bool
}

const (
	_MENU_HIGHLIGHT_HIDDEN = -2
	_MENU_HIGHLIGHT_INIT   = -1
)

func (sr *sdlRender) DisplayMenu(title string, options []string, highlightCallback, selectCallback, cancelCallback types.MenuCallbackT) {
	if highlightCallback == nil {
		highlightCallback = func(_ int) {}
	}
	if selectCallback == nil {
		selectCallback = func(_ int) {}
	}
	if cancelCallback == nil {
		cancelCallback = func(_ int) {}
	}

	sr.footerText = "[Up/Down] Highlight  |  [Return] Choose  |  [Esc] Cancel"
	sr.menu = &menuT{
		title:             title,
		options:           options,
		hidden:            make([]bool, len(options)),
		highlightCallback: highlightCallback,
		selectCallback:    selectCallback,
		cancelCallback:    cancelCallback,
		highlightIndex:    _MENU_HIGHLIGHT_INIT,
	}

	for i := range options {
		if len(options[i]) > int(sr.menu.maxLen) {
			sr.menu.maxLen = int32(len(options[i]))
		}
	}

	sr.term.ShowCursor(false)
	go sr.menu.cursorBlink(sr)
}

func (sr *sdlRender) closeMenu() {
	sr.footerText = ""
	sr.term.ShowCursor(true)
	sr.menu = nil
}

func (menu *menuT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	menu.filter += evt.GetText()

	if menu.highlightIndex < 0 {
		return
	}

	if menu.filter != "" && !strings.Contains(strings.ToLower(menu.options[menu.highlightIndex]), strings.ToLower(menu.filter)) {
		menu.highlightIndex = _MENU_HIGHLIGHT_HIDDEN
	}
}

func (menu *menuT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	mod := keyEventModToCodesModifier(evt.Keysym.Mod)

	var adjust int
	switch evt.Keysym.Sym {
	case sdl.K_RETURN, sdl.K_RETURN2, sdl.K_KP_ENTER:
		if menu.highlightIndex < 0 {
			return
		}
		sr.closeMenu()
		menu.selectCallback(menu.highlightIndex)
		return
	case sdl.K_ESCAPE:
		sr.closeMenu()
		menu.cancelCallback(menu.highlightIndex)
		return

	case sdl.K_BACKSPACE:
		if menu.filter != "" {
			menu.filter = menu.filter[:len(menu.filter)-1]
		} else {
			sr.Bell()
		}
	case sdl.K_u:
		if mod == codes.MOD_CTRL {
			menu.filter = ""
		}

	case sdl.K_UP:
		adjust = -1
	case sdl.K_DOWN:
		adjust = 1
	}

	var attempts int
	for {
		attempts++
		menu.highlightIndex += adjust

		if menu.highlightIndex >= len(menu.options) {
			menu.highlightIndex = 0
		} else if menu.highlightIndex < 0 {
			menu.highlightIndex = len(menu.options) - 1
		}

		if attempts > len(menu.options) {
			menu.highlightIndex = -2
			return
		}

		if menu.hidden[menu.highlightIndex] {
			continue
		}

		if menu.options[menu.highlightIndex] != MENU_SEPARATOR {
			break
		}
	}

	menu.highlightCallback(menu.highlightIndex)
}

func (menu *menuT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	if evt.Button != 1 || evt.State != sdl.RELEASED {
		return
	}
	i := menu._mouseHover(evt.X, evt.Y, sr.glyphSize)
	if i == -1 {
		return
	}

	sr.closeMenu()
	menu.selectCallback(menu.highlightIndex)
}

func (menu *menuT) eventMouseWheel(sr *sdlRender, evt *sdl.MouseWheelEvent) {
	// do nothing
}

func (menu *menuT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	i := menu._mouseHover(evt.X, evt.Y, sr.glyphSize)
	if i == -1 {
		return
	}

	if menu.hidden[i] || menu.options[i] == MENU_SEPARATOR {
		return
	}

	menu.highlightIndex = i
	sr.TriggerRedraw()
	menu.highlightCallback(menu.highlightIndex)
}

func (menu *menuT) _mouseHover(x, y int32, glyphSize *types.XY) int {
	if x < menu.mouseRect.X || x > menu.mouseRect.X+menu.mouseRect.W {
		return -1
	}
	if y < menu.mouseRect.Y || y > menu.mouseRect.Y+menu.mouseRect.H {
		return -1
	}

	rel := y - menu.mouseRect.Y
	i := int(rel / glyphSize.Y)

	if i >= len(menu.options) || menu.options[i] == MENU_SEPARATOR {
		return -1
	}

	return i
}

func (sr *sdlRender) renderMenu(windowRect *sdl.Rect) {
	if sr.menu.highlightIndex == _MENU_HIGHLIGHT_INIT {
		sr.menu.highlightIndex = 0
		sr.menu.highlightCallback(0)
	}

	texture := sr.createRendererTexture()
	if texture == nil {
		return
	}

	padding := int32(10)
	halfPadding := int32(5)

	/*
		FRAME
	*/

	glyphX := sr.glyphSize.X + 1
	iconByGlyphs := (sr.notifyIconSize.X / glyphX) + 1
	maxLen := sr.menu.maxLen
	if int32(len(sr.menu.title))+iconByGlyphs > maxLen {
		maxLen = (int32(len(sr.menu.title)) + iconByGlyphs)
	}
	height := (sr.glyphSize.Y * int32(len(sr.menu.options))) + (padding * 2) + sr.notifyIconSize.Y
	width := (glyphX * maxLen) + (padding * 3)
	menuRect := sdl.Rect{
		X: (windowRect.W - width) / 2,
		Y: (windowRect.H - height) / 2,
		W: width,
		H: height,
	}

	// draw border
	_ = sr.renderer.SetDrawColor(questionColorBorder.Red, questionColorBorder.Green, questionColorBorder.Blue, notificationAlpha)
	rect := sdl.Rect{
		X: menuRect.X - 1,
		Y: menuRect.Y - 1,
		W: menuRect.W + 2,
		H: menuRect.H + 2,
	}
	_ = sr.renderer.DrawRect(&rect)
	_ = sr.renderer.DrawRect(&menuRect)

	// fill background
	_ = sr.renderer.SetDrawColor(questionColor.Red, questionColor.Green, questionColor.Blue, notificationAlpha)
	rect = sdl.Rect{
		X: menuRect.X + 1,
		Y: menuRect.Y + 1,
		W: menuRect.W - 2,
		H: menuRect.H - 2,
	}
	_ = sr.renderer.FillRect(&rect)

	/*
		TITLE
	*/

	surface, err := sdl.CreateRGBSurfaceWithFormat(0, windowRect.W, windowRect.H, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) //TODO: don't panic!
	}
	defer surface.Free()

	sr.font.SetStyle(ttf.STYLE_BOLD)

	text, err := sr.font.RenderUTF8BlendedWrapped(sr.menu.title, sdl.Color{R: 200, G: 200, B: 200, A: 255}, int(glyphX*maxLen))
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	defer text.Free()

	textShadow, err := sr.font.RenderUTF8BlendedWrapped(sr.menu.title, sdl.Color{R: 0, G: 0, B: 0, A: 150}, int(glyphX*maxLen))
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	defer textShadow.Free()

	// render shadow
	rect = sdl.Rect{
		X: menuRect.X + padding + sr.notifyIconSize.X + 2,
		Y: menuRect.Y + padding + 2,
		W: surface.W - (padding * 2),
		H: surface.H - (padding * 2),
	}
	_ = textShadow.Blit(nil, surface, &rect)
	sr._renderNotificationSurface(surface, &rect)

	// render text
	rect = sdl.Rect{
		X: menuRect.X + padding + sr.notifyIconSize.X,
		Y: menuRect.Y + padding,
		W: surface.W - (padding * 2),
		H: surface.H - (padding * 2),
	}
	err = text.Blit(nil, surface, &rect)
	if err != nil {
		panic(err) // TODO: don't panic!
	}
	sr._renderNotificationSurface(surface, &rect)

	// draw border
	offset := sr.notifyIconSize.Y
	width = menuRect.W - padding - padding
	sr.renderer.SetDrawColor(255, 255, 255, 150)
	rect = sdl.Rect{
		X: menuRect.X + padding - 1,
		Y: menuRect.Y + offset - 1,
		W: width + 2, // menuRect.W - padding - padding + 2,
		H: menuRect.H - offset - padding + 2,
	}
	sr.renderer.DrawRect(&rect)

	rect = sdl.Rect{
		X: menuRect.X + padding,
		Y: menuRect.Y + offset,
		W: width, //menuRect.W - padding - padding,
		H: menuRect.H - offset - padding,
	}
	sr.renderer.DrawRect(&rect)

	// fill background
	sr.renderer.SetDrawColor(0, 0, 0, 150)
	rect = sdl.Rect{
		X: menuRect.X + padding + 1,
		Y: menuRect.Y + offset + 1,
		W: width - 2, //menuRect.W - padding - padding - 2,
		H: menuRect.H - offset - padding - 2,
	}
	sr.renderer.FillRect(&rect)

	/*
		MOUSE INTERACTIVE ZONE
	*/

	sr.menu.mouseRect = sdl.Rect{
		X: menuRect.X + padding,
		Y: menuRect.Y + offset + padding,
		W: width,
		H: menuRect.H - offset - padding,
	}

	/*
		OPTIONS
	*/

	offset += halfPadding
	for i := range sr.menu.options {
		if sr.menu.options[i] == MENU_SEPARATOR {
			if sr.menu.filter != "" {
				sr.menu.hidden[i] = true
				continue
			}
			sr.menu.hidden[i] = false

			// draw horizontal separator
			sr.renderer.SetDrawColor(255, 255, 255, 50)
			rect = sdl.Rect{
				X: menuRect.X + padding + halfPadding,
				Y: menuRect.Y + offset + 2 + (sr.glyphSize.Y * int32(i)) + ((sr.glyphSize.Y / 2) - 4),
				W: width - padding - padding,
				H: 4,
			}
			_ = sr.renderer.DrawRect(&rect)
			continue
		}

		if sr.menu.filter != "" && !strings.Contains(strings.ToLower(sr.menu.options[i]), strings.ToLower(sr.menu.filter)) {
			sr.menu.hidden[i] = true
			continue
		}

		sr.menu.hidden[i] = false

		text, err := sr.font.RenderUTF8BlendedWrapped(sr.menu.options[i], sdl.Color{R: 200, G: 200, B: 200, A: 255}, int(surface.W-sr.notifyIconSize.X))
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer text.Free()

		textShadow, err := sr.font.RenderUTF8BlendedWrapped(sr.menu.options[i], sdl.Color{R: 0, G: 0, B: 0, A: 150}, int(surface.W-sr.notifyIconSize.X))
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer textShadow.Free()

		// render shadow
		rect = sdl.Rect{
			X: menuRect.X + padding + halfPadding + 2,
			Y: menuRect.Y + offset + 2 + (sr.glyphSize.Y * int32(i)),
			W: menuRect.X + padding + halfPadding + 2,
			H: surface.H - (padding * 2),
		}
		_ = textShadow.Blit(nil, surface, &rect)
		sr._renderNotificationSurface(surface, &rect)

		// render text
		rect = sdl.Rect{
			X: menuRect.X + padding + halfPadding,
			Y: menuRect.Y + offset + (sr.glyphSize.Y * int32(i)),
			W: surface.W - (padding * 2),
			H: surface.H - (padding * 2),
		}
		err = text.Blit(nil, surface, &rect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		sr._renderNotificationSurface(surface, &rect)
	}

	if surface, ok := sr.notifyIcon[types.NOTIFY_QUESTION].Asset().(*sdl.Surface); ok {
		srcRect := &sdl.Rect{
			X: 0,
			Y: 0,
			W: surface.W,
			H: surface.H,
		}

		dstRect := &sdl.Rect{
			X: menuRect.X,
			Y: menuRect.Y,
			W: sr.notifyIconSize.X,
			H: sr.notifyIconSize.X,
		}

		texture, err := sr.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer texture.Destroy()

		err = sr.renderer.Copy(texture, srcRect, dstRect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}
	}

	sr.AddToOverlayStack(&layer.RenderStackT{texture, windowRect, windowRect, false})
	sr.restoreRendererTexture()

	if sr.menu.filter != "" {
		sr.menu._renderInputBox(sr, surface, windowRect, &sdl.Rect{
			X: sr.menu.mouseRect.X,
			Y: sr.menu.mouseRect.Y + sr.menu.mouseRect.H + padding,
			W: sr.menu.mouseRect.W,
			H: sr.glyphSize.Y + padding,
		})
	}

	if sr.menu.highlightIndex == _MENU_HIGHLIGHT_HIDDEN {
		return
	}

	rect = sdl.Rect{
		X: menuRect.X + padding + halfPadding,
		Y: menuRect.Y + offset + (sr.glyphSize.Y * int32(sr.menu.highlightIndex)),
		W: width - padding,
		H: sr.glyphSize.Y,
	}
	sr._drawHighlightRect(&rect, highlightAlphaBorder, highlightAlphaBorder-20)
}

func (menu *menuT) _renderInputBox(sr *sdlRender, surface *sdl.Surface, windowRect, rect *sdl.Rect) {
	texture := sr.createRendererTexture()
	if texture == nil {
		return
	}

	// draw border
	sr.renderer.SetDrawColor(255, 255, 255, 150)
	borderRect := sdl.Rect{
		X: rect.X - 1,
		Y: rect.Y - 1,
		W: rect.W + 2,
		H: rect.H + 2,
	}
	sr.renderer.DrawRect(&borderRect)
	borderRect = sdl.Rect{
		X: rect.X,
		Y: rect.Y,
		W: rect.W,
		H: rect.H,
	}
	sr.renderer.DrawRect(&borderRect)

	// fill background
	sr.renderer.SetDrawColor(0, 0, 0, 200)
	borderRect = sdl.Rect{
		X: rect.X + 1,
		Y: rect.Y + 1,
		W: rect.W - 2,
		H: rect.H - 2,
	}
	sr.renderer.FillRect(&borderRect)

	// value

	halfPadding := int32(5)
	textRect := sdl.Rect{
		X: rect.X + halfPadding,
		Y: rect.Y + halfPadding,
		W: rect.W,
		H: rect.H,
	}

	var width int32

	if len(sr.menu.filter) > 0 {
		textValue, err := sr.font.RenderUTF8Blended(sr.menu.filter, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer textValue.Free()

		err = textValue.Blit(nil, surface, &textRect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		sr._renderNotificationSurface(surface, rect)
		width = textValue.W
	}

	sr.AddToOverlayStack(&layer.RenderStackT{texture, windowRect, windowRect, false})
	sr.restoreRendererTexture()

	if sr.menu.blinkState {
		cursorRect := sdl.Rect{
			X: textRect.X + width,
			Y: textRect.Y,
			W: sr.glyphSize.X,
			H: sr.glyphSize.Y,
		}
		sr._drawHighlightRect(&cursorRect, 255, 200)
	}
}

func (menu *menuT) cursorBlink(sr *sdlRender) {
	for {
		time.Sleep(500 * time.Millisecond)
		menu.blinkState = !menu.blinkState
		if sr.menu == nil {
			return
		}
	}
}
