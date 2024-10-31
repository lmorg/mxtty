package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

type termWidgetT struct{}

func (tw *termWidgetT) eventTextInput(sr *sdlRender, evt *sdl.TextInputEvent) {
	sr.term.Reply([]byte(evt.GetText()))
}

func (tw *termWidgetT) eventKeyPress(sr *sdlRender, evt *sdl.KeyboardEvent) {
	sr.keyModifier = evt.Keysym.Mod

	switch evt.Keysym.Sym {
	default:
		if evt.Keysym.Sym < 10000 {
			if tw.eventKeyPressMod(sr, evt) {
				return
			}
		}

	case sdl.K_ESCAPE:
		sr.term.Reply([]byte{codes.AsciiEscape})
	case sdl.K_TAB:
		sr.term.Reply([]byte{'\t'})
	case sdl.K_RETURN:
		sr.term.Reply([]byte{'\n'})
	case sdl.K_BACKSPACE:
		sr.term.Reply([]byte{codes.AsciiBackspace})

	case sdl.K_UP:
		sr.ansiReply(codes.AnsiUp, evt.Keysym.Mod)
	case sdl.K_DOWN:
		sr.ansiReply(codes.AnsiDown, evt.Keysym.Mod)
	case sdl.K_LEFT:
		sr.ansiReply(codes.AnsiLeft, evt.Keysym.Mod)
	case sdl.K_RIGHT:
		sr.ansiReply(codes.AnsiRight, evt.Keysym.Mod)

	case sdl.K_PAGEDOWN:
		sr.ansiReply(codes.AnsiPageDown, evt.Keysym.Mod)
	case sdl.K_PAGEUP:
		sr.ansiReply(codes.AnsiPageUp, evt.Keysym.Mod)

	// F-Keys

	case sdl.K_F1:
		sr.ansiReply(codes.AnsiF1, evt.Keysym.Mod)
	case sdl.K_F2:
		sr.ansiReply(codes.AnsiF2, evt.Keysym.Mod)
	case sdl.K_F3:
		sr.ansiReply(codes.AnsiF3, evt.Keysym.Mod)
	case sdl.K_F4:
		sr.ansiReply(codes.AnsiF4, evt.Keysym.Mod)
	case sdl.K_F5:
		sr.ansiReply(codes.AnsiF5, evt.Keysym.Mod)
	case sdl.K_F6:
		sr.ansiReply(codes.AnsiF6, evt.Keysym.Mod)
	case sdl.K_F7:
		sr.ansiReply(codes.AnsiF7, evt.Keysym.Mod)
	case sdl.K_F8:
		sr.ansiReply(codes.AnsiF8, evt.Keysym.Mod)
	case sdl.K_F9:
		sr.ansiReply(codes.AnsiF9, evt.Keysym.Mod)
	case sdl.K_F10:
		sr.ansiReply(codes.AnsiF10, evt.Keysym.Mod)
	case sdl.K_F11:
		sr.ansiReply(codes.AnsiF11, evt.Keysym.Mod)
	case sdl.K_F12:
		sr.ansiReply(codes.AnsiF12, evt.Keysym.Mod)
	case sdl.K_F13:
		sr.ansiReply(codes.AnsiF13, evt.Keysym.Mod)
	case sdl.K_F14:
		sr.ansiReply(codes.AnsiF14, evt.Keysym.Mod)
	case sdl.K_F15:
		sr.ansiReply(codes.AnsiF15, evt.Keysym.Mod)
	case sdl.K_F16:
		sr.ansiReply(codes.AnsiF16, evt.Keysym.Mod)
	case sdl.K_F17:
		sr.ansiReply(codes.AnsiF17, evt.Keysym.Mod)
	case sdl.K_F18:
		sr.ansiReply(codes.AnsiF18, evt.Keysym.Mod)
	case sdl.K_F19:
		sr.ansiReply(codes.AnsiF19, evt.Keysym.Mod)
	case sdl.K_F20:
		sr.ansiReply(codes.AnsiF20, evt.Keysym.Mod)
	}
}

func (tw *termWidgetT) eventKeyPressMod(sr *sdlRender, evt *sdl.KeyboardEvent) (ok bool) {
	switch {
	case evt.Keysym.Mod&sdl.KMOD_CTRL != 0:
		fallthrough
	case evt.Keysym.Mod&sdl.KMOD_LCTRL != 0:
		fallthrough
	case evt.Keysym.Mod&sdl.KMOD_RCTRL != 0:
		if evt.Keysym.Sym > '`' && evt.Keysym.Sym < 'z' {
			sr.term.Reply([]byte{byte(evt.Keysym.Sym) - 0x60})
			return true
		}
	}

	return
}

const (
	_MOUSE_BUTTON_LEFT = 1 + iota
	_MOUSE_BUTTON_MIDDLE
	_MOUSE_BUTTON_RIGHT
	_MOUSE_BUTTON_X1
	_MOUSE_BUTTON_X2
)

func (tw *termWidgetT) eventMouseButton(sr *sdlRender, evt *sdl.MouseButtonEvent) {
	if evt.State == sdl.RELEASED {
		return
	}

	switch evt.Button {
	case _MOUSE_BUTTON_LEFT:
		sr.term.MouseClick(evt.Button, &types.XY{
			X: (evt.X - sr.border) / sr.glyphSize.X,
			Y: (evt.Y - sr.border) / sr.glyphSize.Y,
		}, func() {
			highlighterStart(sr, evt)
			sr.highlighter.setMode(_HIGHLIGHT_MODE_LINE_RANGE)
		})

	case _MOUSE_BUTTON_MIDDLE:
		sr.clipboardPasteText()

	case _MOUSE_BUTTON_RIGHT:
		highlighterStart(sr, evt)

	case _MOUSE_BUTTON_X1:
		highlighterStart(sr, evt)
		sr.highlighter.setMode(_HIGHLIGHT_MODE_SQUARE)
	}
}

func highlighterStart(sr *sdlRender, evt *sdl.MouseButtonEvent) {
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
	if evt.Direction == sdl.MOUSEWHEEL_FLIPPED {
		sr.term.MouseWheel(int(-evt.Y))
	} else {
		sr.term.MouseWheel(int(evt.Y))
	}
}

func (tw *termWidgetT) eventMouseMotion(sr *sdlRender, evt *sdl.MouseMotionEvent) {
	// do nothing
}
