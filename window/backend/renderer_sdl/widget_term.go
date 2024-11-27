package rendersdl

import (
	"fmt"
	"time"

	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

type termWidgetT struct{}

func (tw *termWidgetT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	sr.footerText = ""
	b := []byte(evt.GetText())

	if len(b) == 1 {
		switch b[0] {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'+', '-', '*', '/':
			go func() {
				select {
				case ignore := <-sr.keyIgnore:
					if ignore {
						return
					}
					sr.term.Reply(b)

				case <-time.After(5 * time.Millisecond):
					sr.term.Reply(b)
				}

			}()
			return
		}
	}

	sr.term.Reply(b)
}

func (tw *termWidgetT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	go func() {
		switch evt.Keysym.Sym {
		case sdl.K_KP_1, sdl.K_KP_2, sdl.K_KP_3, sdl.K_KP_4, sdl.K_KP_5,
			sdl.K_KP_6, sdl.K_KP_7, sdl.K_KP_8, sdl.K_KP_9, sdl.K_KP_0,
			sdl.K_KP_PLUS, sdl.K_KP_MINUS, sdl.K_KP_MULTIPLY, sdl.K_KP_DIVIDE:
			go func() {
				sr.keyIgnore <- true
			}()
		}
		//log.Printf("key: %d", evt.Keysym.Sym)

	}()
	tw._eventKeyPress(sr, evt)
}

func (tw *termWidgetT) _eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
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

	switch {
	case keyCode == codes.AnsiF3 && mod == 0:
		sr.term.Search()
		return
	}

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
				switch evt.Clicks {
				case 1:
					sr.selectWindow(i - 1)
				case 2:
					sr.DisplayInputBox("Please enter a new name for this window:", sr.windowTabs.windows[i-1].Name, func(name string) {
						err := sr.windowTabs.windows[i-1].Rename(name)
						if err != nil {
							sr.DisplayNotification(types.NOTIFY_ERROR, err.Error())
						}
					})
				}
				return
			}
		}
		if evt.Clicks == 2 {
			sr.tmux.NewWindow()
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
				if x >= 0 && x < sr.windowTabs.boundaries[i] {
					sr.windowTabs.mouseOver = i - 1
					sr.footerText = fmt.Sprintf("[click]  Switch to window '%s' (%s)", sr.windowTabs.windows[i-1].Name, sr.windowTabs.windows[i-1].Id)
					return
				}
			}
			sr.footerText = "[2x click]  Start new window"
			sr.windowTabs.mouseOver = -1
			return
		}

		sr.windowTabs.mouseOver = -1
		sr.footerText = ""
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

func (sr *sdlRender) selectWindow(winIndex int) {
	if winIndex < 0 || winIndex >= len(sr.windowTabs.windows) {
		return
	}

	winId := sr.windowTabs.windows[winIndex].Id
	err := sr.tmux.SelectWindow(winId)
	if err != nil {
		sr.DisplayNotification(types.NOTIFY_ERROR, err.Error())
	}
}

func (sr *sdlRender) RefreshWindowList() {
	if sr.tmux == nil {
		return
	}

	sr.limiter.Lock()

	sr.windowTabs = nil
	sr.term = sr.tmux.ActivePane().Term()

	sr.limiter.Unlock()
}
