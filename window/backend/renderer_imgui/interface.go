package rendererimgui

import "github.com/lmorg/mxtty/types"

type imguiRender struct{}

func (renderer *imguiRender) Size() *types.Rect {
	return termSize
}

func (renderer *imguiRender) Close() {
	//typeface.Close()
	//window.Destroy()
	//sdl.Quit()
}

func (renderer *imguiRender) SetWindowTitle(title string) {
	/*
		unsupported in SDL due to:
		NSWindow geometry should only be modified on the main thread!
	*/

	//window.SetTitle(fmt.Sprintf("%s - %s", title, app.Name))
}
