package assets

import "embed"

const BELL = "bell.mp3"

//go:embed bell.mp3
var embedFsBell embed.FS

func init() {

	embedFs := embedFsBell

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
