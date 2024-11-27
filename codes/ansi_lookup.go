package codes

import (
	"fmt"

	"github.com/lmorg/mxtty/debug"
	"github.com/lmorg/mxtty/types"
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

var _retryLookUpTable = map[types.KeyboardMode]*[]types.KeyboardMode{
	types.KeysNormal:      {types.KeysNormal, types.KeysVT220},
	types.KeysApplication: {types.KeysApplication, types.KeysNormal, types.KeysVT220},
	types.KeysVT52:        {types.KeysVT52, types.KeysNormal},
	types.KeysVT220:       {types.KeysVT220, types.KeysNormal},
	types.KeysTmuxClient:  {types.KeysTmuxClient},
}

var _ansiLookUpTable = map[types.KeyboardMode]map[KeyCode][]byte{
	types.KeysNormal: {
		AnsiUp:    csi('A'),
		AnsiDown:  csi('B'),
		AnsiRight: csi('C'),
		AnsiLeft:  csi('D'),
		AnsiHome:  csi('H'),
		AnsiEnd:   csi('E'),

		AnsiKeyPadSpace: []byte{' '},
		AnsiKeyPadTab:   []byte{'\t'},
		AnsiKeyPadEnter: []byte{'\r'},

		/*AnsiShiftTab:    csi('Z'),
		AnsiOptUp:       []byte{esc, esc, '[', 'A'},
		AnsiOptDown:     []byte{esc, esc, '[', 'B'},
		AnsiOptLeft:     []byte{esc, esc, '[', 'D'},
		AnsiOptRight:    []byte{esc, esc, '[', 'C'},
		AnsiCtrlUp:      csi('1', ';', '5', 'A'),
		AnsiCtrlDown:    csi('1', ';', '5', 'B'),
		AnsiCtrlLeft:    csi('1', ';', '5', 'D'),
		AnsiCtrlRight:   csi('1', ';', '5', 'C'),*/

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

		/*AnsiShiftF1:  csi('1', ';', '2', 'P'),
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
		AnsiShiftF12: csi('2', '4', ';', '2', '~'),*/
	},

	types.KeysApplication: {
		AnsiUp:    ss3('A'),
		AnsiDown:  ss3('B'),
		AnsiRight: ss3('C'),
		AnsiLeft:  ss3('D'),
		AnsiHome:  ss3('H'),
		AnsiEnd:   ss3('E'),
	},

	types.KeysVT220: {
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

	types.KeysVT52: {
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

	types.KeysTmuxClient: {
		8:    _tmuxKeyResponce("BSpace"),
		'\t': _tmuxKeyResponce("Tab"),
		'\r': _tmuxKeyResponce("Enter"),
		'\n': _tmuxKeyResponce("Enter"),
		' ':  _tmuxKeyResponce("Space"),
		esc:  _tmuxKeyResponce("Escape"),
		'"':  _tmuxKeyResponce(`'"'`),
		'\'': _tmuxKeyResponce(`"'"`),
		128:  _tmuxKeyResponce("BSpace"),

		AnsiUp:       _tmuxKeyResponce("Up"),
		AnsiDown:     _tmuxKeyResponce("Down"),
		AnsiRight:    _tmuxKeyResponce("Right"),
		AnsiLeft:     _tmuxKeyResponce("Left"),
		AnsiHome:     _tmuxKeyResponce("Home"),
		AnsiEnd:      _tmuxKeyResponce("End"),
		AnsiInsert:   _tmuxKeyResponce("Insert"),
		AnsiDelete:   _tmuxKeyResponce("Delete"),
		AnsiPageUp:   _tmuxKeyResponce("PageUp"),
		AnsiPageDown: _tmuxKeyResponce("PageDown"),

		AnsiKeyPadSpace:    _tmuxKeyResponce("Space"),
		AnsiKeyPadTab:      _tmuxKeyResponce("Tab"),
		AnsiKeyPadEnter:    _tmuxKeyResponce("KPEnter"),
		AnsiKeyPadMultiply: _tmuxKeyResponce("KP*"),
		AnsiKeyPadAdd:      _tmuxKeyResponce("KP+"),
		AnsiKeyPadComma:    _tmuxKeyResponce("KP,"),
		AnsiKeyPadMinus:    _tmuxKeyResponce("KP-"),
		AnsiKeyPadPeriod:   _tmuxKeyResponce("KP."),
		AnsiKeyPadDivide:   _tmuxKeyResponce("KP/"),
		AnsiKeyPad0:        _tmuxKeyResponce("KP0"),
		AnsiKeyPad1:        _tmuxKeyResponce("KP1"),
		AnsiKeyPad2:        _tmuxKeyResponce("KP2"),
		AnsiKeyPad3:        _tmuxKeyResponce("KP3"),
		AnsiKeyPad4:        _tmuxKeyResponce("KP4"),
		AnsiKeyPad5:        _tmuxKeyResponce("KP5"),
		AnsiKeyPad6:        _tmuxKeyResponce("KP6"),
		AnsiKeyPad7:        _tmuxKeyResponce("KP7"),
		AnsiKeyPad8:        _tmuxKeyResponce("KP8"),
		AnsiKeyPad9:        _tmuxKeyResponce("KP9"),
		AnsiKeyPadEqual:    _tmuxKeyResponce("KP="),

		AnsiF1:  _tmuxKeyResponce("F1"),
		AnsiF2:  _tmuxKeyResponce("F2"),
		AnsiF3:  _tmuxKeyResponce("F3"),
		AnsiF4:  _tmuxKeyResponce("F4"),
		AnsiF5:  _tmuxKeyResponce("F5"),
		AnsiF6:  _tmuxKeyResponce("F6"),
		AnsiF7:  _tmuxKeyResponce("F7"),
		AnsiF8:  _tmuxKeyResponce("F8"),
		AnsiF9:  _tmuxKeyResponce("F9"),
		AnsiF10: _tmuxKeyResponce("F10"),
		AnsiF11: _tmuxKeyResponce("F11"),
		AnsiF12: _tmuxKeyResponce("F12"),
		// F13->F20 are not supported in tmux
	},
}

func getAnsiEscSeq(keySet types.KeyboardMode, keyPress KeyCode) []byte {
	b, ok := _ansiLookUpTable[keySet][keyPress]
	if !ok {
		if keySet != types.KeysNormal {
			return getAnsiEscSeq(types.KeysNormal, keyPress)
		}

		debug.Log(fmt.Sprintf("No sequence available for %d in %d", keyPress, keySet))
		return b
	}

	return b
}

// TODO:
// As a special case, the SS3 sent before F1 through F4 is altered to CSI when
// sending a function key modifier as a parameter.
func GetAnsiEscSeq(keySet types.KeyboardMode, keyPress KeyCode, modifier Modifier) []byte {
	// check for hardcoded exceptions
	b := specialCaseSequences(keySet, keyPress, modifier)
	if len(b) != 0 {
		return b
	}

	// fallback to generalized formats
	lookupSet := _retryLookUpTable[keySet]
	for _, set := range *lookupSet {
		b = getAnsiEscSeq(set, keyPress)
		if len(b) != 0 {
			// no modifiers
			if modifier == 0 {
				return b
			}

			// contains modifiers
			return spliceKeysAndModifiers(b, modifier)
		}
	}

	return b
}

func _tmuxKeyResponce(keyName string) []byte {
	return append([]byte{0}, []byte(keyName+" ")...)
}
