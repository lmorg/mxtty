package assets

import (
	"embed"
)

//go:embed *
var embedFs embed.FS

var assets map[string][]byte

func init() {
	assets = make(map[string][]byte)

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

		assets[name] = b
	}
}

func Get(name string) []byte {
	return assets[name]
}
