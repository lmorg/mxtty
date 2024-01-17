package rendersdl

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/lmorg/mxtty/types"
	"github.com/veandco/go-sdl2/sdl"
)

var notifyColour = map[int]*types.Colour{
	types.NOTIFY_DEBUG:   {Red: 200, Green: 200, Blue: 200},
	types.NOTIFY_INFO:    {Red: 0, Green: 0, Blue: 255},
	types.NOTIFY_WARNING: {Red: 0, Green: 255, Blue: 255},
	types.NOTIFY_ERROR:   {Red: 255, Green: 0, Blue: 0},
}

type notifyT struct {
	stack []*notificationT
	mutex sync.Mutex
}

type notificationT struct {
	Type    types.NotificationType
	Message string
	wait    <-chan time.Time
	end     time.Time
}

func (n *notifyT) _wait() {
	for {
		if len(n.stack) == 0 {
			return
		}

		<-n.stack[0].wait
		n.remove()
	}
}

func (n *notifyT) add(notification *notificationT) {
	d := 5 * time.Second
	notification.end = time.Now().Add(d)
	notification.wait = time.After(d)

	n.mutex.Lock()
	n.stack = append(n.stack, notification)

	if len(n.stack) > 0 {
		go n._wait()
	}
	n.mutex.Unlock()

	log.Printf("NOTIFICATION: %s", notification.Message)
}

func (n *notifyT) remove() {
	n.mutex.Lock()
	n.stack = n.stack[1:]
	n.mutex.Unlock()
}

func (n *notifyT) get() []*notificationT {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if len(n.stack) == 0 {
		return nil
	}

	notifications := make([]*notificationT, len(n.stack))
	copy(notifications, n.stack)

	return notifications
}

func (sr *sdlRender) DisplayNotification(notificationType types.NotificationType, message string) {
	notification := &notificationT{
		Type:    notificationType,
		Message: message,
	}
	sr.notifications.add(notification)
}

func (sr *sdlRender) renderNotification(windowRect *sdl.Rect) {
	notifications := sr.notifications.get()
	if notifications == nil {
		return
	}

	surface, err := sdl.CreateRGBSurfaceWithFormat(0, windowRect.W, windowRect.H, 32, uint32(sdl.PIXELFORMAT_RGBA32))
	if err != nil {
		panic(err) //TODO: don't panic!
	}
	defer surface.Free()

	sr.setFontStyle(types.SGR_BOLD)

	padding := sr.border * 2
	var offset int32
	for _, notification := range notifications {
		// generate text
		s := strconv.Itoa(int(time.Until(notification.end)/time.Second) + 1)
		countdown, err := sr.font.RenderUTF8Blended(s, sdl.Color{R: 255, G: 255, B: 255, A: 200})
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer countdown.Free()

		text, err := sr.font.RenderUTF8BlendedWrapped(notification.Message, sdl.Color{R: 0, G: 0, B: 0, A: 255}, int(windowRect.W-padding-padding-countdown.W))
		if err != nil {
			panic(err) // TODO: don't panic!
		}
		defer text.Free()

		// draw border
		rect := sdl.Rect{
			X: sr.border,
			Y: sr.border + offset,
			W: windowRect.W - padding,
			H: text.H + padding,
		}
		c := notifyColour[int(notification.Type)]
		sr.renderer.SetDrawColor(c.Red, c.Green, c.Blue, 255)
		sr.renderer.DrawRect(&rect)

		// fill background
		rect = sdl.Rect{
			X: sr.border + 1,
			Y: sr.border + 1 + offset,
			W: sr.surface.W - padding - 2,
			H: text.H + padding - 2,
		}
		sr.renderer.SetDrawColor(c.Red, c.Green, c.Blue, 190)
		sr.renderer.FillRect(&rect)

		// render text
		rect = sdl.Rect{
			X: windowRect.W - padding - countdown.W,
			Y: padding + offset,
			W: countdown.W,
			H: countdown.H,
		}
		err = countdown.Blit(nil, surface, &rect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}

		texture, err := sr.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err) //TODO: don't panic!
		}

		err = sr.renderer.Copy(texture, &rect, &rect)
		if err != nil {
			panic(err) //TODO: don't panic!
		}

		rect = sdl.Rect{
			X: padding,
			Y: padding + offset,
			W: sr.surface.W - padding - 2 - countdown.W,
			H: text.H + padding - 2,
		}
		err = text.Blit(nil, surface, &rect)
		if err != nil {
			panic(err) // TODO: don't panic!
		}

		texture, err = sr.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			panic(err) //TODO: don't panic!
		}

		err = sr.renderer.Copy(texture, &rect, &rect)
		if err != nil {
			panic(err) //TODO: don't panic!
		}

		offset += text.H + (sr.border * 3)
	}
}
