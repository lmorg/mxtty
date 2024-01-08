package rendererimgui

import (
	//"github.com/lmorg/mxtty/imgui"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/virtualterm"
)

var (
	backend imgui.Backend[imgui.SDLWindowFlags] //[imgui.GLFWWindowFlags]

	glyphSize *types.Rect
	termSize  *types.Rect
	border    int32 = 5
	width     int   = 1024
	height    int   = 768
)

func Initialise(fontName string, fontSize int) types.Renderer {
	err := createWindow("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	setTypeFace()

	return new(imguiRender)

	/*font, err := typeface.Open(fontName, fontSize)
	if err != nil {
		panic(err.Error())
	}

	setTypeFace(font)

	return new(sdlRender)*/
}

func createWindow(caption string) error {
	/*var err error

	// Create a window for us to draw the text on
	window, err = sdl.CreateWindow(
		caption,                                          // window title
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, // window pos
		width, height, // window dimensions
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE, // window properties
	)
	if err != nil {
		return err
	}

	surface, err = window.GetSurface()
	return err*/

	//backend = imgui.CreateBackend(imgui.NewGLFWBackend())
	backend = imgui.CreateBackend(imgui.NewSDLBackend())
	//backend.SetAfterCreateContextHook(afterCreateContext)
	//backend.SetBeforeDestroyContextHook(beforeDestroyContext)

	backend.CreateWindow(caption, width, height)

	return nil
}

func setTypeFace() {
	//font = f
	glyphSize = &types.Rect{16, 16}
	termSize = getTermSize()
}

func getTermSize() *types.Rect {
	x, y := backend.DisplaySize()

	return &types.Rect{
		X: (x - (border * 2)) / glyphSize.X,
		Y: (y - (border * 2)) / glyphSize.Y,
	}
}

func colour(colour *types.Colour) imgui.Vec4 {
	return imgui.NewVec4(float32(colour.Red)/0xFF, float32(colour.Green)/0xFF, float32(colour.Blue)/0xFF, 1)
}

func Start(term *virtualterm.Term) {
	imgui.Begin("bg")
	backend.SetBgColor(colour(virtualterm.SGR_COLOUR_BLACK))
	imgui.End()
	for {
		/*c := virtualterm.SGR_COLOUR_BLACK
		pixel := sdl.MapRGBA(surface.Format, c.Red, c.Green, c.Blue, 255)
		err := surface.FillRect(&sdl.Rect{W: surface.W, H: surface.H}, pixel)
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

			sdl.Delay(15)
		}*/
		time.Sleep(15 * time.Millisecond)
	}
}
