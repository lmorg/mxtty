package rendersdl

import (
	"bytes"
	"fmt"
	stdlib_image "image"
	"image/png"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
	"golang.design/x/clipboard"
)

func (sr *sdlRender) copySurfaceToClipboard() {
	err := sr._copySurfaceToClipboard(sr.surface, sr.highlighter.rect)
	if err != nil {
		sr.DisplayNotification(types.NOTIFY_ERROR, fmt.Sprintf("Could not copy to clipboard: %s", err.Error()))
	} else {
		sr.DisplayNotification(types.NOTIFY_INFO, "Copied to clipboard as PNG")
	}
	sr.highlighter = nil
}

func (sr *sdlRender) _copySurfaceToClipboard(src *sdl.Surface, rect *sdl.Rect) error {
	dstRect := sdl.Rect{
		W: rect.W,
		H: rect.H,
	}
	surf, err := sdl.CreateRGBSurfaceWithFormat(0, rect.W, rect.H, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		return err
	}

	err = src.Blit(rect, surf, &dstRect)
	if err != nil {
		return err
	}

	img := stdlib_image.NewRGBA(stdlib_image.Rect(0, 0, int(rect.W), int(rect.H)))
	copy(img.Pix, surf.Pixels())

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtImage, buf.Bytes())
	return nil
}

func (sr *sdlRender) clipboardPasteText() {
	sr.highlighter = nil
	b := clipboard.Read(clipboard.FmtText)
	if len(b) != 0 {
		sr.term.Reply(b)
	} else {
		sr.DisplayNotification(types.NOTIFY_WARN, "Clipboard does not contain text to paste")
	}
}
