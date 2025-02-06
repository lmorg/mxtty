package virtualterm

import (
	"github.com/lmorg/mxtty/config"
	"github.com/lmorg/mxtty/types"
)

func (term *Term) Render() {
	if !term.visible {
		return
	}

	term._mutex.Lock()

	screen := term.visibleScreen()

	if !config.Config.Terminal.TypeFace.Ligatures || term._mouseButtonDown || term._searchHighlight {
		term._renderCells(screen)
	} else {
		term._renderLigatures(screen)
	}

	term._renderOutputBlockChrome(screen)

	term._renderCursor()

	term._mutex.Unlock()
}

func (term *Term) _renderCells(screen types.Screen) {
	pos := new(types.XY)
	elementStack := make(map[types.Element]bool) // no duplicates

	for ; pos.Y < term.size.Y; pos.Y++ {
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			switch {
			case screen[pos.Y].Cells[pos.X].Element != nil:
				_, ok := elementStack[screen[pos.Y].Cells[pos.X].Element]
				if !ok {
					elementStack[screen[pos.Y].Cells[pos.X].Element] = true
					offset := getElementXY(screen[pos.Y].Cells[pos.X].Char)
					screen[pos.Y].Cells[pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}

			case screen[pos.Y].Cells[pos.X].Char == 0:
				continue

			case screen[pos.Y].Cells[pos.X].Sgr == nil:
				continue

			default:
				if screen[pos.Y].Cells[pos.X].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCell(screen[pos.Y].Cells[pos.X], pos)
			}
		}
	}
}

func (term *Term) _renderLigatures(screen types.Screen) {
	var (
		pos          = new(types.XY)
		elementStack = make(map[types.Element]bool) // no duplicates
		hash         uint64
		defaultHash  = types.SGR_DEFAULT.HashValue()
	)

	for ; pos.Y < term.size.Y; pos.Y++ {
		if screen[pos.Y].Cells[0].Sgr == nil {
			hash = defaultHash
		} else {
			hash = screen[pos.Y].Cells[0].Sgr.HashValue()
		}

		var start int32
		for pos.X = 0; pos.X < term.size.X; pos.X++ {
			newHash := defaultHash // ^uint64(0)
			if screen[pos.Y].Cells[pos.X].Sgr != nil {
				newHash = screen[pos.Y].Cells[pos.X].Sgr.HashValue()
			}

			if screen[pos.Y].Cells[pos.X].Element != nil {
				_, ok := elementStack[screen[pos.Y].Cells[pos.X].Element]
				if !ok {
					elementStack[screen[pos.Y].Cells[pos.X].Element] = true
					offset := getElementXY(screen[pos.Y].Cells[pos.X].Char)
					screen[pos.Y].Cells[pos.X].Element.Draw(nil, &types.XY{X: pos.X - offset.X, Y: pos.Y - offset.Y})
				}
				newHash = defaultHash //^uint64(0)
			}

			if hash != newHash {
				if screen[pos.Y].Cells[start].Sgr != nil && screen[pos.Y].Cells[start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
					continue // blink
				}
				term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:pos.X], &types.XY{X: start, Y: pos.Y})
				hash = newHash
				start = pos.X
			}

			if screen[pos.Y].Cells[pos.X].Char == 0 || screen[pos.Y].Cells[pos.X].Element != nil {
				term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:pos.X], &types.XY{X: start, Y: pos.Y})
				start = pos.X + 1
			}
		}

		if start < pos.X {
			if screen[pos.Y].Cells[start].Sgr != nil && screen[pos.Y].Cells[start].Sgr.Bitwise.Is(types.SGR_SLOW_BLINK) && !term._slowBlinkState {
				continue // blink
			}
			term.renderer.PrintCellBlock(screen[pos.Y].Cells[start:], &types.XY{X: start, Y: pos.Y})
		}
	}
}

func (term *Term) _renderOutputBlockChrome(screen types.Screen) {
	var (
		foundEnd   bool
		i          int32
		errorBlock bool
	)

	term._cacheBlock = [][]int32{}

	for y := int32(len(screen)) - 1; y >= 0; y-- {
		i++
		if len(screen[y].Hidden) != 0 {
			term.renderer.DrawOutputBlockChrome(y, 1, types.COLOUR_FOLDED, true)
		}
		if screen[y].Meta.Is(types.ROW_OUTPUT_BLOCK_END) {
			i = 0
			errorBlock = false
			foundEnd = true
		}
		if screen[y].Meta.Is(types.ROW_OUTPUT_BLOCK_ERROR) {
			i = 0
			errorBlock = true
			foundEnd = true
		}

		if screen[y].Meta.Is(types.ROW_OUTPUT_BLOCK_BEGIN) {
			if !foundEnd {
				_, row, err := term.outputBlocksFindStartEnd(int32(len(term._scrollBuf)-term._scrollOffset) + y)
				if err != nil {
					continue
				}
				i--
				errorBlock = row[1].Meta.Is(types.ROW_OUTPUT_BLOCK_ERROR)
			}

			_renderOutputBlockChrome(term, y, i, errorBlock)
			foundEnd = false
			i = 0
		}
	}

	if foundEnd {
		_renderOutputBlockChrome(term, 0, i, errorBlock)
	}

	if len(term._cacheBlock) == 0 {
		_, row, err := term.outputBlocksFindStartEnd(int32(len(term._scrollBuf) - term._scrollOffset))
		if err != nil {
			return
		}

		errorBlock = row[1].Meta.Is(types.ROW_OUTPUT_BLOCK_ERROR)
		_renderOutputBlockChrome(term, 0, int32(len(screen))-1, errorBlock)
	}
}

func _renderOutputBlockChrome(term *Term, start, end int32, errorBlock bool) {
	end++
	if errorBlock {
		term.renderer.DrawOutputBlockChrome(start, end, types.COLOUR_ERROR, false)
	} else {
		term.renderer.DrawOutputBlockChrome(start, end, types.COLOUR_OK, false)
	}
	term._cacheBlock = append(term._cacheBlock, []int32{start, end})
}
