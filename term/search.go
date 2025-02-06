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

	term.renderer.DisplayInputBox("Value to search for", term._searchLastString, term.searchBuf)
}

func (term *Term) searchBuf(search string) {
	search = strings.ToLower(search)
	term._searchLastString = search

	if search == "" {
		term._searchHighlight = false
		for _, cell := range term._searchHlHistory {
			if cell != nil && cell.Sgr != nil {
				cell.Sgr.Bitwise.Unset(types.SGR_HIGHLIGHT_SEARCH_RESULT)
			}
		}
		term._searchHlHistory = []*types.Cell{}
		term._searchResults = nil
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
		//term.ShowCursor(false)
		term.updateScrollback()
		return
	}

	term.renderer.DisplayNotification(types.NOTIFY_WARN, fmt.Sprintf("Search string not found: '%s'", search))
}

func (term *Term) _searchBuf(buf types.Screen, search string) (int, bool) {
	firstMatch := -1
	for y := len(buf) - 1; y >= 0; y-- {
		for x := len(buf[y].Cells) - 1; x >= 0; x-- {
			if buf[y].Cells[x].Phrase == nil {
				continue
			}

			s := strings.ToLower(string(*buf[y].Cells[x].Phrase))
			if strings.Contains(s, search) {
				term._searchResults = append(term._searchResults, searchResult{
					rowId:  buf[y].Id,
					phrase: string(*buf[y].Cells[x].Phrase),
				})
				i, j, l := 0, 0, 0
			highlight:
				for ; x+i < len(buf[y].Cells); i++ {
					buf[y+j].Cells[x+i].Sgr = buf[y+j].Cells[x+i].Sgr.Copy()
					buf[y+j].Cells[x+i].Sgr.Bitwise.Set(types.SGR_HIGHLIGHT_SEARCH_RESULT)
					term._searchHlHistory = append(term._searchHlHistory, buf[y+j].Cells[x+i])
					l++
				}
				if l < len(*buf[y].Cells[x].Phrase) {
					i = 0
					j++
					goto highlight
				}

				if firstMatch == -1 {
					firstMatch = y
				}
			}
		}
	}
	return firstMatch, firstMatch != -1
}

func (term *Term) ShowSearchResults() {
	offset := term._scrollOffset
	sr := make([]searchResult, len(term._searchResults))
	results := make([]string, len(term._searchResults))
	j := len(term._searchResults) - 1

	for i := range term._searchResults {
		sr[j] = term._searchResults[i]
		results[j] = term._searchResults[i].phrase
		j--
	}

	cbHighlight := func(i int) {
		term.scrollToRowId(sr[i].rowId, _SEARCH_OFFSET)
	}
	cbCancel := func(int) {
		term._scrollOffset = offset
		term.updateScrollback()
	}
	cbSelect := func(int) {}
	term.renderer.DisplayMenu("Search results", results, cbHighlight, cbSelect, cbCancel)

	/*if term._scrollOffset != 0 {
		term.ShowCursor(false)
	}*/
}
