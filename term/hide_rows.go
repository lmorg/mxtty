package virtualterm

import (
	"errors"
	"fmt"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
)

func clone[T any](src []T) []T {
	s := make([]T, len(src))
	copy(s, src)
	return s
}

func (term *Term) HideRows(start int32, end int32) error {
	if term.IsAltBuf() {
		return errors.New("this feature is not supported in alt buffer")
	}

	term._mutex.Lock()
	defer term._mutex.Unlock()

	newBuf := term._scrollBuf
	newBuf = append(newBuf, term._normBuf...)

	if len(newBuf[start-1].Hidden) != 0 {
		return errors.New("this row already contains hidden rows")
	}

	newBuf[start-1].Hidden = clone(newBuf[start:end])
	debug.Log(newBuf[start-1].Hidden.String())
	length := len(newBuf[start-1].Hidden)
	newBuf = append(newBuf[:start], newBuf[end:]...)

	if len(newBuf) < int(term.size.Y) {
		newBuf = append(term.makeScreen(), newBuf...)
	}

	if term._scrollOffset > 0 {
		term._scrollOffset -= int(end - start)
	}
	term.updateScrollback()

	term._normBuf = clone(newBuf[len(newBuf)-int(term.size.Y):])
	term._scrollBuf = clone(newBuf[:len(newBuf)-int(term.size.Y)])

	term.renderer.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%d rows have been hidden", length))

	return nil
}

func (term *Term) UnhideRows(pos int32) error {
	if term.IsAltBuf() {
		return errors.New("this feature is not supported in alt buffer")
	}

	term._mutex.Lock()
	defer term._mutex.Unlock()

	tmp := term._scrollBuf
	tmp = append(tmp, term._normBuf...)

	length := len(tmp[pos].Hidden)
	debug.Log(tmp[pos].Hidden.String())
	newBuf := append(clone(tmp[:pos+1]), tmp[pos].Hidden...)
	tmp[pos].Hidden = nil
	newBuf = clone(append(newBuf, tmp[pos+1:]...))

	term._normBuf = clone(newBuf[len(newBuf)-int(term.size.Y):])
	term._scrollBuf = clone(newBuf[:len(newBuf)-int(term.size.Y)])

	term.renderer.DisplayNotification(types.NOTIFY_INFO, fmt.Sprintf("%d rows have been unhidden", length))

	return nil
}
