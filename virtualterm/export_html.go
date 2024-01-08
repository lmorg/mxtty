package virtualterm

import (
	"html"
	"strings"
)

// ExportHTML returns a HTML render of the virtual terminal
func (t *Term) ExportHtml() string {
	s := `<span class="">`

	lastSgr := &sgr{}
	var lastChar rune = 0

	t.mutex.Lock()

	for y := range t.cells {
		for x := range t.cells[y] {
			sgr := t.cells[y][x].sgr
			char := t.cells[y][x].char

			if t.cells[y][x].differs(lastChar, lastSgr) {
				s += `</span><span class="` + sgrHtmlClassLookup(sgr) + `">`
			}

			if char != 0 { // if cell contains no data then lets assume it's a space character
				s += html.EscapeString(string(char))
			} else {
				s += " "
			}

			lastSgr = sgr
			lastChar = char
		}
		s += "\n"
	}

	t.mutex.Unlock()

	return s + "</span>"
}

func sgrHtmlClassLookup(sgr *sgr) string {
	classes := make([]string, 0)

	for bit, class := range sgrHtmlClassNames {
		if sgr.checkFlag(bit) {
			classes = append(classes, class)
		}
	}

	if sgr.checkFlag(sgrFgColour4) {
		classes = append(classes, sgrColourHtmlClassNames[sgr.fg.Red])
	}

	return strings.Join(classes, " ")
}
