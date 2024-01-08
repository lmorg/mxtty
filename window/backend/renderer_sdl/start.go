package rendersdl

import (
	"github.com/lmorg/mxtty/codes"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
	"github.com/lmorg/mxtty/window/backend/typeface"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	window    *sdl.Window
	surface   *sdl.Surface
	font      *ttf.Font
	glyphSize *types.Rect
	termSize  *types.Rect
	border    int32 = 5
	width     int32 = 1024
	height    int32 = 768
)

func Initialise() types.Renderer {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err.Error())
	}

	err = createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open("hasklig.ttf", 14)
	//font, err := typeface.Open("monaco.ttf", 16)
	if err != nil {
		panic(err.Error())
	}

	setTypeFace(font)

	return new(sdlRender)
}

func createWindow(caption string) error {
	var err error

	// Create a window for us to draw the text on
	window, err = sdl.CreateWindow(caption, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}

	surface, err = window.GetSurface()
	return err
}

func setTypeFace(f *ttf.Font) {
	font = f
	glyphSize = typeface.GetSize()
	termSize = getTermSize()
}

func getTermSize() *types.Rect {
	x, y := window.GetSize()

	return &types.Rect{
		X: (x - (border * 2)) / glyphSize.X,
		Y: (y - (border * 2)) / glyphSize.Y,
	}
}

func Start(term *virtualterm.Term) {
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

		sdl.Delay(15)
	}
}
