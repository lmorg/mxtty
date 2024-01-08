package virtualterm

type sgrFlag uint32

// Flags
const (
	sgrReset sgrFlag = 0

	sgrBold sgrFlag = 1 << iota
	sgrItalic
	sgrUnderscore
	sgrBlink
	sgrInvert
)
