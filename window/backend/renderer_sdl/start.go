package rendersdl

import (
	"log"
	"sync/atomic"

	"github.com/lmorg/mxtty/assets"
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	border int32 = 5
	width  int32 = 1024
	height int32 = 768
)

var focused *sdlRender

func Initialise(fontName string, fontSize int) types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	focused = new(sdlRender)
	err = focused.createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open(fontName, fontSize)
	if err != nil {
		panic(err.Error())
	}

	focused.setTypeFace(font)

	return focused
}

func (sr *sdlRender) createWindow(caption string) error {
	var err error

	// Create a window for us to draw the text on
	sr.window, err = sdl.CreateWindow(
		caption,                                          // window title
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, // window pos
		width, height, // window dimensions
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE, // window properties
	)
	if err != nil {
		return err
	}

	rwops, err := sdl.RWFromMem(assets.Get(assets.ICON_BMP))
	if err != nil {
		return err
	}

	icon, err := sdl.LoadBMPRW(rwops, true)
	if err != nil {
		return err
	}

	sr.window.SetIcon(icon)

	sr.surface, err = sr.window.GetSurface()
	return err
}

func (sr *sdlRender) setTypeFace(f *ttf.Font) {
	sr.font = f
	sr.glyphSize = typeface.GetSize()
	sr.termSize = sr.getTermSize()
}

func (sr *sdlRender) getTermSize() *types.Rect {
	x, y := sr.window.GetSize()

	return &types.Rect{
		X: (x - (border * 2)) / sr.glyphSize.X,
		Y: (y - (border * 2)) / sr.glyphSize.Y,
	}
}

func Start(term *virtualterm.Term) {
	c := virtualterm.SGR_COLOUR_BLACK
	pixel := sdl.MapRGBA(focused.surface.Format, c.Red, c.Green, c.Blue, 255)
	err := focused.surface.FillRect(&sdl.Rect{W: focused.surface.W, H: focused.surface.H}, pixel)
	if err != nil {
		log.Printf("error drawing background: %s", err.Error())
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TextInputEvent:
				term.Pty.Secondary.WriteString(evt.GetText())

			case *sdl.KeyboardEvent:
				if evt.State == sdl.RELEASED {
					break
				}

				switch evt.Keysym.Sym {
				case sdl.K_ESCAPE:
					term.Pty.Secondary.Write([]byte{codes.AsciiEscape})
				case sdl.K_TAB:
					term.Pty.Secondary.Write([]byte{'\t'})
				case sdl.K_RETURN:
					term.Pty.Secondary.Write([]byte{'\n'})
				case sdl.K_BACKSPACE:
					term.Pty.Secondary.Write([]byte{codes.IsoBackspace})

				case sdl.K_UP:
					term.Pty.Secondary.Write(codes.AnsiUp)
				case sdl.K_DOWN:
					term.Pty.Secondary.Write(codes.AnsiDown)
				case sdl.K_LEFT:
					term.Pty.Secondary.Write(codes.AnsiBackwards)
				case sdl.K_RIGHT:
					term.Pty.Secondary.Write(codes.AnsiForwards)

				case sdl.K_PAGEDOWN:
					term.Pty.Secondary.Write(codes.AnsiPageDown)
				case sdl.K_PAGEUP:
					term.Pty.Secondary.Write(codes.AnsiPageUp)

				// F-Keys

				case sdl.K_F1:
					term.Pty.Secondary.Write(codes.AnsiF1VT100)
				case sdl.K_F2:
					term.Pty.Secondary.Write(codes.AnsiF2VT100)
				case sdl.K_F3:
					term.Pty.Secondary.Write(codes.AnsiF3VT100)
				case sdl.K_F4:
					term.Pty.Secondary.Write(codes.AnsiF4VT100)
				case sdl.K_F5:
					term.Pty.Secondary.Write(codes.AnsiF5)
				case sdl.K_F6:
					term.Pty.Secondary.Write(codes.AnsiF6)
				case sdl.K_F7:
					term.Pty.Secondary.Write(codes.AnsiF7)
				case sdl.K_F8:
					term.Pty.Secondary.Write(codes.AnsiF8)
				case sdl.K_F9:
					term.Pty.Secondary.Write(codes.AnsiF9)
				case sdl.K_F10:
					term.Pty.Secondary.Write(codes.AnsiF10)
				case sdl.K_F11:
					term.Pty.Secondary.Write(codes.AnsiF11)
				case sdl.K_F12:
					term.Pty.Secondary.Write(codes.AnsiF12)

				}
			}
		}

		sdl.Delay(5)
		term.Render()

		if atomic.CompareAndSwapInt32(&focused.updateTitle, 1, 0) {
			focused.window.SetTitle(focused.title)
		}

		err = focused.window.UpdateSurface()
		if err != nil {
			log.Printf("error in renderer: %s", err.Error())
		}
	}
}
