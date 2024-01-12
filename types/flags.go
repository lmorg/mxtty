package types

type SgrFlag uint32

// Flags
const (
	SGR_NORMAL SgrFlag = 0

	SGR_BOLD SgrFlag = 1 << iota
	SGR_ITALIC
	SGR_UNDERLINE
	SGR_STRIKETHROUGH
	SGR_SLOW_BLINK
	SGR_INVERT

	APC_ELEMENT
	APC_BEGIN_ELEMENT
)

func (f SgrFlag) Is(flag SgrFlag) bool {
	return f&flag != 0
}

func (f *SgrFlag) Set(flag SgrFlag) {
	*f |= flag
}

func (f *SgrFlag) Unset(flag SgrFlag) {
	*f &^= flag
}
