package integrations


import (
	"embed"
	"fmt"
)

const (
	BELL          = "bell.mp3"
	ICON_APP      = "icon.bmp"
	TYPEFACE      = "Hasklig-Regular.ttf"
	ICON_DEBUG    = "icon-message.bmp"
	ICON_INFO     = "icon-info.bmp"
	ICON_WARN     = "icon-warn.bmp"
	ICON_ERROR    = "icon-error.bmp"
	ICON_DOWN     = "icon-down.bmp"
	ICON_QUESTION = "icon-question.bmp"
)

//go:embed shell.*
var embedFs embed.FS

var integrations map[string][]byte

func init() {
	integrations = make(map[string][]byte)

	dir, err := embedFs.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for i := range dir {
		name := dir[i].Name()

		b, err := embedFs.ReadFile(name)
		if err != nil {
			// not a bug in murex
			panic(err)
		}

		integrations[name] = b
	}
}

func Get(name string) []byte {
	b, ok := integrations[name]
	if !ok {
		panic(fmt.Sprintf("no asset found named '%s'", name))
	}
	return b
}
