package rendersdl

import (
	"bytes"
	"fmt"
	stdlib_image "image"
	"image/png"
	"unsafe"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"golang.design/x/clipboard"
)

func (sr *sdlRender) copyRendererToClipboard() {
	defer func() {
		sr.highlighter = nil
		sr.renderer.SetRenderTarget(nil)
		sr.TriggerRedraw()
	}()

	w, h := sr.window.GetSize()
	debug.Log(fmt.Sprintf("w:%d, h:%d", w, h))

	pitch := w * 4
	pixelData := make([]uint8, pitch*h)

	debug.Log("readpixels")
	err := sr.renderer.ReadPixels(&sdl.Rect{W: w, H: h}, uint32(sdl.PIXELFORMAT_RGBA8888), unsafe.Pointer(&pixelData), int(pitch))
	if err != nil {
		debug.Log(err)
		sr.DisplayNotification(types.NOTIFY_ERROR, fmt.Sprintf("Could not copy to clipboard: %s", err.Error()))
		return
	}
	return
	img := stdlib_image.NewRGBA(stdlib_image.Rect(0, 0, int(w), int(h)))
	copy(img.Pix, pixelData)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		sr.DisplayNotification(types.NOTIFY_ERROR, fmt.Sprintf("Could not copy to clipboard: %s", err.Error()))
		return
	}

	clipboard.Write(clipboard.FmtImage, buf.Bytes())
	sr.DisplayNotification(types.NOTIFY_INFO, "Copied to clipboard as PNG")
}

/*func (sr *sdlRender) _copySurfaceToClipboard(src *sdl.Surface, rect *sdl.Rect) error {
	dstRect := sdl.Rect{
		W: rect.W,
		H: rect.H,
	}
	surf, err := sdl.CreateRGBSurfaceWithFormat(0, rect.W, rect.H, 32, uint32(sdl.PIXELFORMAT_RGBA8888))
	if err != nil {
		return err
	}
	defer surf.Free()

	err = src.Blit(rect, surf, &dstRect)
	if err != nil {
		return err
	}

	img := stdlib_image.NewRGBA(stdlib_image.Rect(0, 0, int(rect.W), int(rect.H)))
	copy(img.Pix, surf.Pixels())
	//copy(img.Pix, pixels)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtImage, buf.Bytes())
	return nil
}*/

func (sr *sdlRender) clipboardPasteText() {
	sr.highlighter = nil
	b := clipboard.Read(clipboard.FmtText)
	if len(b) != 0 {
		sr.term.Reply(b)
	} else {
		sr.DisplayNotification(types.NOTIFY_WARN, "Clipboard does not contain text to paste")
	}
}
