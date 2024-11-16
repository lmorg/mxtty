package rendersdl

import (
	"fmt"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/tmux"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/utils/octal"
	"github.com/veandco/go-sdl2/sdl"
)

type termWidgetT struct{}

func (tw *termWidgetT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	sr.footerText = ""
	b := []byte(evt.GetText())

	if config.Config.Tmux.Enabled {
		b = octal.Escape(b)
	}

	sr.term.Reply(b)
}

func (tw *termWidgetT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	sr.footerText = ""
	sr.keyModifier = evt.Keysym.Mod

	switch evt.Keysym.Sym {
	case sdl.K_LSHIFT, sdl.K_RSHIFT, sdl.K_LALT, sdl.K_RALT,
		sdl.K_LCTRL, sdl.K_RCTRL, sdl.K_LGUI, sdl.K_RGUI,
		sdl.K_CAPSLOCK, sdl.K_NUMLOCKCLEAR, sdl.K_SCROLLLOCK, sdl.K_SPACE:
		// modifier keys pressed on their own shouldn't trigger anything
		return
	}

	if evt.Keysym.Sym < 256 && evt.Keysym.Sym > sdl.K_SPACE &&
		(evt.Keysym.Mod == sdl.KMOD_NONE ||
			evt.Keysym.Mod&sdl.KMOD_CAPS != 0 || evt.Keysym.Mod&sdl.KMOD_NUM != 0) {
		// lets let eventTextInput() handle this so we don't need to think
		// about keyboard layouts and shift chars like if shift+2 == '@' or '"'
		return
	}

	mod := keyEventModToCodesModifier(evt.Keysym.Mod)
	keyCode := sr.keyCodeLookup(evt.Keysym.Sym)

	b := codes.GetAnsiEscSeq(sr.keyboardMode.Get(), keyCode, mod)
	if len(b) > 0 {
		sr.term.Reply(b)
	}
}

// SDL doesn't name these, so lets name it ourselves for convenience
const (
	_MOUSE_BUTTON_LEFT = 1 + iota
	_MOUSE_BUTTON_MIDDLE
	_MOUSE_BUTTON_RIGHT
	_MOUSE_BUTTON_X1
	_MOUSE_BUTTON_X2
)

func (tw *termWidgetT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	posCell := sr.convertPxToCellXY(evt.X, evt.Y)

	if config.Config.Tmux.Enabled && sr.windowTabs != nil &&
		(evt.Y-sr.border)/sr.glyphSize.Y == sr.term.GetSize().Y+sr.footer-1 {
		// window tab bar
		if evt.State == sdl.PRESSED {
			return
		}

		x := ((evt.X - sr.border) / sr.glyphSize.X) - sr.windowTabs.offset.X
		for i := range sr.windowTabs.boundaries {
			if x < sr.windowTabs.boundaries[i] {
				sr._switchWindow(i - 1)
				return
			}
		}

		return
	}

	if evt.State == sdl.RELEASED {
		sr.term.MouseClick(posCell, evt.Button, evt.Clicks, false, func() {})
		return
	}

	switch evt.Button {
	case _MOUSE_BUTTON_LEFT:
		sr.term.MouseClick(posCell, evt.Button, evt.Clicks, true, func() {
			highlighterStart(sr, evt)
			sr.highlighter.setMode(_HIGHLIGHT_MODE_LINE_RANGE)
		})

	case _MOUSE_BUTTON_MIDDLE:
		sr.term.MouseClick(posCell, evt.Button, evt.Clicks, true, sr.clipboardPasteText)

	case _MOUSE_BUTTON_RIGHT:
		sr.term.MouseClick(posCell, evt.Button, evt.Clicks, true, sr.clipboardPasteText)

	case _MOUSE_BUTTON_X1:
		sr.term.MouseClick(posCell, evt.Button, evt.Clicks, true, func() {})
	}
}

var _highlighterStartFooterText = fmt.Sprintf(
	"Copy to clipboard: [%s] Square region  |  [%s] Text region  |  [%s] Entire line(s)  |  [%s] PNG",
	types.KEY_STR_CTRL, types.KEY_STR_SHIFT, types.KEY_STR_ALT, types.KEY_STR_META,
)

func highlighterStart(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	sr.footerText = _highlighterStartFooterText

	sr.highlighter = &highlighterT{
		button: evt.Button,
		rect:   &sdl.Rect{X: evt.X, Y: evt.Y},
	}
	if sr.keyModifier != 0 {
		sr.highlighter.modifier(sr.keyModifier)
	}
	sr.keyModifier = 0
}

func (tw *termWidgetT) eventMouseWheel(sr *sdlRender, evt *sdl.MouseWheelEvent) {
	mouseX, mouseY, _ := sdl.GetMouseState()

	if evt.Direction == sdl.MOUSEWHEEL_FLIPPED {
		sr.term.MouseWheel(sr.convertPxToCellXY(mouseX, mouseY), &types.XY{X: evt.X, Y: -evt.Y})
	} else {
		sr.term.MouseWheel(sr.convertPxToCellXY(mouseX, mouseY), &types.XY{X: evt.X, Y: evt.Y})
	}
}

func (tw *termWidgetT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	if config.Config.Tmux.Enabled && sr.windowTabs != nil {

		if (evt.Y-sr.border)/sr.glyphSize.Y == sr.term.GetSize().Y+sr.footer-1 {
			x := ((evt.X - sr.border) / sr.glyphSize.X) - sr.windowTabs.offset.X
			for i := range sr.windowTabs.boundaries {
				if x < sr.windowTabs.boundaries[i] {
					sr.windowTabs.mouseOver = i - 1
					return
				}
			}
		}

		sr.windowTabs.mouseOver = -1
	}

	sr.term.MouseMotion(
		sr.convertPxToCellXY(evt.X, evt.Y),
		&types.XY{
			X: evt.XRel / sr.glyphSize.X,
			Y: evt.YRel / sr.glyphSize.Y,
		},
		sr._termMouseMotionCallback,
	)
}

func (sr *sdlRender) _termMouseMotionCallback() {
	sr.footerText = "[left click] Copy  |  [right click] Paste  |  [wheel] Scrollback buffer"
}

func (sr *sdlRender) _switchWindow(winIndex int) {
	// update old
	sr.windowTabs.windows[sr.windowTabs.active].Active = false
	sr.windowTabs.windows[winIndex].ActivePane().Term().MakeVisible(false)

	// select new
	sr.windowTabs.active = winIndex

	// update new
	sr.windowTabs.windows[winIndex].Active = true
	sr.windowTabs.windows[winIndex].ActivePane().Term().MakeVisible(true)
	command := fmt.Sprintf("%s -t %s", tmux.CMD_SELECT_WINDOW, sr.windowTabs.windows[winIndex].Id)
	_, err := sr.tmux.SendCommand([]byte(command))
	if err != nil {
		sr.DisplayNotification(types.NOTIFY_ERROR, err.Error())
	}
	sr.term = sr.windowTabs.windows[winIndex].ActivePane().Term()
}

func (sr *sdlRender) RefreshWindowList() {
	sr.limiter.Lock()
	sr.windowTabs = nil
	sr.limiter.Unlock()
}
