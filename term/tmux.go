package virtualterm

import "github.com/lmorg/mxtty/codes"

func (term *Term) tmuxRenameWindow() {
	var title []rune
	for {
		r := term.Pty.Read()
		if r == codes.AsciiEscape {
			if term.Pty.Read() == '\\' {
				break
			}
			continue
		}
		title = append(title, r)
	}

	term.renderer.SetWindowTitle(string(title))
}
