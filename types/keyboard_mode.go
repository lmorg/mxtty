package types

type KeyboardMode int32

const (
	KeysNormal KeyboardMode = 0 + iota
	KeysApplication
	KeysVT220
	KeysVT52
)
