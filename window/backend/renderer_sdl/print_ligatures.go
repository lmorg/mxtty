package rendersdl

import (
	"fmt"
	"time"

	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
	"github.com/lmorg/mxtty/window/backend/renderer_sdl/layer"
	"github.com/veandco/go-sdl2/sdl"
)

type cachedLigaturesT struct {
	cache map[uint64]map[string]*cachedLigatureT
}

type cachedLigatureT struct {
	texture *sdl.Texture
	rect    *sdl.Rect
	ttl     time.Time
}

func newCachedLigatures(renderer *sdlRender) *cachedLigaturesT {
	cl := &cachedLigaturesT{
		cache: make(map[uint64]map[string]*cachedLigatureT),
	}

	go func() {
		for {
			time.Sleep(60 * time.Second)
			if !renderer.limiter.TryLock() {
				continue // we don't want to be flushing the cache while rendering to screen
			}

			for _, m := range cl.cache {
				for s, cache := range m {
					if cache.ttl.Before(time.Now()) {
						cache.texture.Destroy()
						delete(m, s)
					}
				}
			}

			renderer.limiter.Unlock()
		}
	}()

	return cl
}

// Get returns nil if unsuccessful
func (cl *cachedLigaturesT) Get(hash uint64, text string) *cachedLigatureT {
	m, ok := cl.cache[hash]
	if !ok {
		return nil
	}

	cache, ok := m[text]
	if ok {
		cache.ttl = time.Now().Add(5 * time.Second)
	}
	return cache
}

func (cl *cachedLigaturesT) Store(hash uint64, text string, texture *sdl.Texture, rect *sdl.Rect) {
	m, ok := cl.cache[hash]
	if !ok {
		cl.cache[hash] = make(map[string]*cachedLigatureT)
		m = cl.cache[hash]
	}
	m[text] = &cachedLigatureT{
		texture: texture,
		rect:    rect,
		ttl:     time.Now().Add(20 * time.Second),
	}
}

// PrintCellBlock is much slower because it doesn't cache textures
func (sr *sdlRender) PrintCellBlock(cells []types.Cell, cellPos *types.XY) {
	r := make([]rune, len(cells))
	for i := 0; i < len(r); i++ {
		r[i] = cells[i].Char
	}

	if len(r) == 0 {
		return
	}

	s := string(r)

	hash := cells[0].Sgr.HashValue()
	cache := sr.ligCache.Get(hash, s)
	if cache != nil {
		dstRect := &sdl.Rect{
			X: (sr.glyphSize.X * cellPos.X) + sr.border,
			Y: (sr.glyphSize.Y * cellPos.Y) + sr.border,
			W: cache.rect.W,
			H: cache.rect.H,
		}
		sr.AddToElementStack(&layer.RenderStackT{cache.texture, cache.rect, dstRect, false})
		return
	}

	surface := _newFontSurface(sr.glyphSize, int32(len(cells)))
	defer surface.Free()

	sgr := cells[0].Sgr
	if sgr == nil {
		sgr = types.SGR_DEFAULT
	}

	sr.font.SetStyle(fontStyle(sgr.Bitwise))

	fg, bg := sgrOpts(sgr)

	cellBlockRect := &sdl.Rect{
		W: sr.glyphSize.X * int32(len(cells)),
		H: sr.glyphSize.Y,
	}

	// render background colour

	if bg != nil {
		pixel := sdl.MapRGBA(surface.Format, bg.Red, bg.Green, bg.Blue, 255)
		err := surface.FillRect(cellBlockRect, pixel)
		if err != nil {
			panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
		}
	}

	if config.Config.Terminal.TypeFace.DropShadow && bg == nil {
		shadowRect := &sdl.Rect{
			X: cellBlockRect.X + dropShadowOffset,
			Y: cellBlockRect.Y + dropShadowOffset,
			W: cellBlockRect.W,
			H: cellBlockRect.H,
		}

		c := textShadow
		shadowText, err := sr.font.RenderUTF8Blended(s, c)
		if err != nil {
			panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
		}
		defer shadowText.Free()

		err = shadowText.Blit(nil, surface, shadowRect)
		if err != nil {
			panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
		}
	}

	// render cell char
	text, err := sr.font.RenderUTF8Blended(s, sdl.Color{R: fg.Red, G: fg.Green, B: fg.Blue, A: 255})
	if err != nil {
		panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
	}
	defer text.Free()

	err = text.Blit(nil, surface, cellBlockRect)
	if err != nil {
		panic(err) // TODO: better error handling please!
	}
	if sgr.Bitwise.Is(types.SGR_BOLD) {
		text.SetBlendMode(sdl.BLENDMODE_ADD)
		err = text.Blit(nil, surface, cellBlockRect)
		if err != nil {
			panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
		}
	}

	texture, err := sr.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(fmt.Sprintf("error printing '%s' (%d): %v", s, len(s), err)) // TODO: better error handling please!
	}

	dstRect := &sdl.Rect{
		X: (sr.glyphSize.X * cellPos.X) + sr.border,
		Y: (sr.glyphSize.Y * cellPos.Y) + sr.border,
		W: cellBlockRect.W, // + dropShadowOffset,
		H: cellBlockRect.H, // + dropShadowOffset,
	}

	sr.AddToElementStack(&layer.RenderStackT{texture, cellBlockRect, dstRect, false})
	sr.ligCache.Store(hash, s, texture, cellBlockRect)
}
