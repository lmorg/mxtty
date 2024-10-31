package codes

import (
	"fmt"

	"github.com/lmorg/mxtty/debug"
)

type KeyboardMode int32
type FunctionKey int

const (
	KeysNormal KeyboardMode = 0 + iota
	KeysApplication
	KeysVT220
	KeysVT52
)

const (
	AnsiUp FunctionKey = 0 + iota
	AnsiDown
	AnsiRight
	AnsiLeft
	AnsiInsert
	AnsiHome
	AnsiEnd
	AnsiDelete
	AnsiPageUp
	AnsiPageDown

	AnsiKeyPadSpace
	AnsiKeyPadTab
	AnsiKeyPadEnter
	AnsiKeyPadMultiply
	AnsiKeyPadAdd
	AnsiKeyPadComma
	AnsiKeyPadMinus
	AnsiKeyPadPeriod
	AnsiKeyPadDivide
	AnsiKeyPad0
	AnsiKeyPad1
	AnsiKeyPad2
	AnsiKeyPad3
	AnsiKeyPad4
	AnsiKeyPad5
	AnsiKeyPad6
	AnsiKeyPad7
	AnsiKeyPad8
	AnsiKeyPad9
	AnsiKeyPadEqual

	AnsiShiftTab

	AnsiOptUp
	AnsiOptDown
	AnsiOptLeft
	AnsiOptRight

	AnsiCtrlUp
	AnsiCtrlDown
	AnsiCtrlLeft
	AnsiCtrlRight

	AnsiF1
	AnsiF2
	AnsiF3
	AnsiF4
	AnsiF5
	AnsiF6
	AnsiF7
	AnsiF8
	AnsiF9
	AnsiF10
	AnsiF11
	AnsiF12
	AnsiF13
	AnsiF14
	AnsiF15
	AnsiF16
	AnsiF17
	AnsiF18
	AnsiF19
	AnsiF20

	AnsiShiftF1
	AnsiShiftF2
	AnsiShiftF3
	AnsiShiftF4
	AnsiShiftF5
	AnsiShiftF6
	AnsiShiftF7
	AnsiShiftF8
	AnsiShiftF9
	AnsiShiftF10
	AnsiShiftF11
	AnsiShiftF12
)

const esc = 27

/*
	Reference documentation used:
	- ASCII table: https://upload.wikimedia.org/wikipedia/commons/thumb/1/1b/ASCII-Table-wide.svg/1280px-ASCII-Table-wide.svg.png
	- xterm: https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys
*/

var (
	Ss2 = []byte{esc, 'N'}
	Ss3 = []byte{esc, 'O'}
	Csi = []byte{esc, '['}
)

func ss3(b ...byte) []byte { return append(Ss3, b...) }
func csi(b ...byte) []byte { return append(Csi, b...) }

var ansiEscapeSeq = map[KeyboardMode]map[FunctionKey][]byte{
	KeysNormal: {
		AnsiUp:          csi('A'),
		AnsiDown:        csi('B'),
		AnsiRight:       csi('C'),
		AnsiLeft:        csi('D'),
		AnsiHome:        csi('H'),
		AnsiEnd:         csi('E'),
		AnsiKeyPadSpace: []byte{' '},
		AnsiKeyPadTab:   []byte{'\t'},
		AnsiKeyPadEnter: []byte{'\r'},
		AnsiShiftTab:    csi('Z'),
		AnsiOptUp:       []byte{esc, esc, '[', 'A'},
		AnsiOptDown:     []byte{esc, esc, '[', 'B'},
		AnsiOptLeft:     []byte{esc, esc, '[', 'D'},
		AnsiOptRight:    []byte{esc, esc, '[', 'C'},
		AnsiCtrlUp:      csi('1', ';', '5', 'A'),
		AnsiCtrlDown:    csi('1', ';', '5', 'B'),
		AnsiCtrlLeft:    csi('1', ';', '5', 'D'),
		AnsiCtrlRight:   csi('1', ';', '5', 'C'),

		AnsiF1:  ss3('P'),
		AnsiF2:  ss3('Q'),
		AnsiF3:  ss3('R'),
		AnsiF4:  ss3('S'),
		AnsiF5:  csi('1', '5', '~'),
		AnsiF6:  csi('1', '7', '~'),
		AnsiF7:  csi('1', '8', '~'),
		AnsiF8:  csi('1', '9', '~'),
		AnsiF9:  csi('2', '0', '~'),
		AnsiF10: csi('2', '1', '~'),
		AnsiF11: csi('2', '3', '~'),
		AnsiF12: csi('2', '4', '~'),

		AnsiShiftF1:  csi('1', ';', '2', 'P'),
		AnsiShiftF2:  csi('1', ';', '2', 'Q'),
		AnsiShiftF3:  csi('1', ';', '2', 'R'),
		AnsiShiftF4:  csi('1', ';', '2', 'S'),
		AnsiShiftF5:  csi('1', '5', ';', '2', '~'),
		AnsiShiftF6:  csi('1', '7', ';', '2', '~'),
		AnsiShiftF7:  csi('1', '8', ';', '2', '~'),
		AnsiShiftF8:  csi('1', '9', ';', '2', '~'),
		AnsiShiftF9:  csi('2', '0', ';', '2', '~'),
		AnsiShiftF10: csi('2', '1', ';', '2', '~'),
		AnsiShiftF11: csi('2', '3', ';', '2', '~'),
		AnsiShiftF12: csi('2', '4', ';', '2', '~'),
	},

	KeysApplication: {
		AnsiUp:    ss3('A'),
		AnsiDown:  ss3('B'),
		AnsiRight: ss3('C'),
		AnsiLeft:  ss3('D'),
		AnsiHome:  ss3('H'),
		AnsiEnd:   ss3('E'),
	},

	KeysVT220: {
		AnsiHome:     csi('1', '~'),
		AnsiInsert:   csi('2', '~'),
		AnsiDelete:   csi('3', '~'),
		AnsiEnd:      csi('4', '~'),
		AnsiPageUp:   csi('5', '~'),
		AnsiPageDown: csi('6', '~'),

		AnsiKeyPadSpace:    ss3(' '),
		AnsiKeyPadTab:      ss3('I'),
		AnsiKeyPadEnter:    ss3('M'),
		AnsiKeyPadMultiply: ss3('j'),
		AnsiKeyPadAdd:      ss3('k'),
		AnsiKeyPadComma:    ss3('l'),
		AnsiKeyPadMinus:    ss3('m'),
		AnsiKeyPadPeriod:   ss3('n'),
		AnsiKeyPadDivide:   ss3('o'),
		AnsiKeyPad0:        ss3('p'),
		AnsiKeyPad1:        ss3('q'),
		AnsiKeyPad2:        ss3('r'),
		AnsiKeyPad3:        ss3('s'),
		AnsiKeyPad4:        ss3('t'),
		AnsiKeyPad5:        ss3('u'),
		AnsiKeyPad6:        ss3('v'),
		AnsiKeyPad7:        ss3('w'),
		AnsiKeyPad8:        ss3('x'),
		AnsiKeyPad9:        ss3('y'),
		AnsiKeyPadEqual:    ss3('X'),

		AnsiF13: csi('2', '5', '~'),
		AnsiF14: csi('2', '6', '~'),
		AnsiF15: csi('2', '8', '~'),
		AnsiF16: csi('2', '9', '~'),
		AnsiF17: csi('3', '1', '~'),
		AnsiF18: csi('3', '2', '~'),
		AnsiF19: csi('3', '3', '~'),
		AnsiF20: csi('3', '4', '~'),
	},

	KeysVT52: {
		AnsiUp:    []byte{esc, 'A'},
		AnsiDown:  []byte{esc, 'B'},
		AnsiRight: []byte{esc, 'C'},
		AnsiLeft:  []byte{esc, 'D'},

		AnsiKeyPadSpace:    []byte{esc, '?', ' '},
		AnsiKeyPadTab:      []byte{esc, '?', '\t'},
		AnsiKeyPadEnter:    []byte{esc, '?', 'M'},
		AnsiKeyPadMultiply: []byte{esc, '?', 'j'},
		AnsiKeyPadAdd:      []byte{esc, '?', 'k'},
		AnsiKeyPadComma:    []byte{esc, '?', 'l'},
		AnsiKeyPadMinus:    []byte{esc, '?', 'm'},
		AnsiKeyPadPeriod:   []byte{esc, '?', 'n'},
		AnsiKeyPadDivide:   []byte{esc, '?', 'o'},
		AnsiKeyPad0:        []byte{esc, '?', 'p'},
		AnsiKeyPad1:        []byte{esc, '?', 'q'},
		AnsiKeyPad2:        []byte{esc, '?', 'r'},
		AnsiKeyPad3:        []byte{esc, '?', 's'},
		AnsiKeyPad4:        []byte{esc, '?', 't'},
		AnsiKeyPad5:        []byte{esc, '?', 'u'},
		AnsiKeyPad6:        []byte{esc, '?', 'v'},
		AnsiKeyPad7:        []byte{esc, '?', 'w'},
		AnsiKeyPad8:        []byte{esc, '?', 'x'},
		AnsiKeyPad9:        []byte{esc, '?', 'y'},
		AnsiKeyPadEqual:    []byte{esc, '?', 'X'},
	},
}

func getAnsiEscSeq(keySet KeyboardMode, keyPress FunctionKey) []byte {
	b, ok := ansiEscapeSeq[keySet][keyPress]
	if !ok {
		if keySet != KeysNormal {
			return getAnsiEscSeq(KeysNormal, keyPress)
		}

		debug.Log(fmt.Sprintf("No sequence available for %d in %d", keyPress, keySet))
		return b
	}

	return b
}

// TODO:
// As a special case, the SS3  sent before F1 through F4 is altered to CSI when
// sending a function key modifier as a parameter.
func GetAnsiEscSeq(keySet KeyboardMode, keyPress FunctionKey, modifier Modifier) []byte {
	b := getAnsiEscSeq(keySet, keyPress)
	if len(b) == 0 || modifier == 0 {
		return b
	}

	ending := b[len(b)-1]
	seq := append(b[:len(b)-1], translateModToCode(modifier)...)
	return append(seq, ending)
}
