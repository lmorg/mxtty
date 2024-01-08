package main

const stuff = `func main() {
	defer typeface.Close()
	defer window.Close()

	err := window.Create("mxtty - Multimedia Terminal Emulator")
	if err != nil {
		panic(err.Error())
	}

	font, err := typeface.Open("hasklig.ttf", 14)
	if err != nil {
		panic(err.Error())
	}

	err = window.PrintText(font, out "hello world" -> grep 'world')
	if err != nil {
		panic(err.Error())
	}

	err = window.Update()
	if err != nil {
		panic(err.Error())
	}

	// Run infinite loop until user closes the window
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		sdl.Delay(16)
	}
}`
