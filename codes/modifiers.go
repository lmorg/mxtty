package codes

/*
	  Code     Modifiers
	---------+---------------------------
	   2     | Shift
	   3     | Alt
	   4     | Shift + Alt
	   5     | Control
	   6     | Shift + Control
	   7     | Alt + Control
	   8     | Shift + Alt + Control
	   9     | Meta
	   10    | Meta + Shift
	   11    | Meta + Alt
	   12    | Meta + Alt + Shift
	   13    | Meta + Ctrl
	   14    | Meta + Ctrl + Shift
	   15    | Meta + Ctrl + Alt
	   16    | Meta + Ctrl + Alt + Shift
	---------+---------------------------

	For example, shift-F5 would be sent as CSI 1 5 ; 2 ~
*/

type Modifier int

const (
	MOD_SHIFT Modifier = 1 << iota
	MOD_ALT
	MOD_CTRL
	MOD_META
)

func translateModToCode(mod Modifier) []byte {
	switch mod {
	case MOD_SHIFT:
		return []byte{';', '2'}

	case MOD_ALT:
		return []byte{';', '3'}

	case MOD_SHIFT | MOD_ALT:
		return []byte{';', '4'}

	case MOD_CTRL:
		return []byte{';', '5'}

	case MOD_SHIFT | MOD_CTRL:
		return []byte{';', '6'}

	case MOD_ALT | MOD_CTRL:
		return []byte{';', '7'}

	case MOD_SHIFT | MOD_ALT | MOD_CTRL:
		return []byte{';', '8'}

	case MOD_META:
		return []byte{';', '9'}

	case MOD_META | MOD_SHIFT:
		return []byte{';', '1', '0'}

	case MOD_META | MOD_ALT:
		return []byte{';', '1', '1'}

	case MOD_META | MOD_ALT | MOD_SHIFT:
		return []byte{';', '1', '2'}

	case MOD_META | MOD_CTRL:
		return []byte{';', '1', '3'}

	case MOD_META | MOD_CTRL | MOD_SHIFT:
		return []byte{';', '1', '4'}

	case MOD_META | MOD_CTRL | MOD_ALT:
		return []byte{';', '1', '5'}

	case MOD_META | MOD_CTRL | MOD_ALT | MOD_SHIFT:
		return []byte{';', '1', '6'}

	default:
		panic("invalid modifier")
	}
}
