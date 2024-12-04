package virtualterm

import (
	"fmt"
	"strings"

	"github.com/lmorg/mxtty/types"
)

const _SEARCH_OFFSET = 3

func (term *Term) Search() {
	if term.IsAltBuf() {
		term.renderer.DisplayNotification(types.NOTIFY_WARN, "Search is not supported in alt buffer")
		return
	}

	term.renderer.DisplayInputBox("Value to search for:", term._searchLastString, term.searchBuf)
}

func (term *Term) searchBuf(search string) {
	term._searchLastString = search

	if search == "" {
		term._searchHighlight = false
		for _, cell := range term._searchHlHistory {
			cell.Sgr.Bitwise.Unset(types.SGR_HIGHLIGHT_SEARCH_RESULT)
		}
		term._searchHlHistory = []*types.Cell{}
		term._scrollOffset = 0
		term.updateScrollback()
		return
	}

	term._mutex.Lock()
	defer term._mutex.Unlock()

	_, normOk := term._searchBuf(term._normBuf, search)
	offset, scrollOk := term._searchBuf(term._scrollBuf, search)

	term._searchHighlight = term._searchHighlight || normOk || scrollOk

	if normOk {
		return
	}

	if scrollOk {
		// +_SEARCH_OFFSET to add some before context
		term._scrollOffset = len(term._scrollBuf) - offset + _SEARCH_OFFSET
		term.ShowCursor(false)
		term.updateScrollback()
		return
	}

	term.renderer.DisplayNotification(types.NOTIFY_WARN, fmt.Sprintf("Search string not found: '%s'", search))
}

func (term *Term) _searchBuf(buf [][]types.Cell, search string) (int, bool) {
	firstMatch := -1
	for y := len(buf) - 1; y >= 0; y-- {
		for x := len(buf[y]) - 1; x >= 0; x-- {
			if buf[y][x].Phrase == nil {
				continue
			}
			s := string(*buf[y][x].Phrase)
			if strings.Contains(s, search) {
				i, j, l := 0, 0, 0
			highlight:
				for ; x+i < len(buf[y]); i++ {
					buf[y+j][x+i].Sgr = buf[y+j][x+i].Sgr.Copy()
					buf[y+j][x+i].Sgr.Bitwise.Set(types.SGR_HIGHLIGHT_SEARCH_RESULT)
					term._searchHlHistory = append(term._searchHlHistory, &buf[y+j][x+i])
					l++
				}
				if l < len(*buf[y][x].Phrase) {
					i = 0
					j++
					goto highlight
				}

				if firstMatch == -1 {
					firstMatch = y
				}

				//return y, true
			}
		}
	}
	return firstMatch, firstMatch != -1
}
