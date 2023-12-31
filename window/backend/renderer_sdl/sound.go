package rendersdl

import (
	"log"

	"github.com/lmorg/mxtty/assets"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func (sr *sdlRender) initBell() {
	err := sdl.Init(sdl.INIT_AUDIO)
	if err != nil {
		log.Printf("ERROR: could not initialise audio sound system: %s", err.Error())
	}

	rwops, err := sdl.RWFromMem(assets.Get(assets.BELL))
	if err != nil {
		log.Printf("ERROR: could not load %s from resource pack: %s", assets.BELL, err.Error())
	}

	err = mix.Init(mix.INIT_MP3)
	if err != nil {
		log.Printf("ERROR: could not initialise sound mixer: %s", err.Error())
		return
	}

	err = mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 1, 4096)
	if err != nil {
		log.Printf("ERROR: could not open audio: %s", err.Error())
		return
	}

	sr.bell, err = mix.LoadMUSRW(rwops, 0)
	if err != nil {
		log.Printf("ERROR: could not load %s from memory: %s", assets.BELL, err.Error())
		return
	}
}

func (sr *sdlRender) Bell() {
	err := sr.bell.Play(1)
	if err != nil {
		log.Printf("ERROR: could not play %s from memory: %s", assets.BELL, err.Error())
		return
	}

	sdl.Delay(5000)
}
