package types

type NotificationType int

const (
	NOTIFY_DEBUG = iota
	NOTIFY_INFO
	NOTIFY_WARN
	NOTIFY_ERROR
	NOTIFY_SCROLL
)

type Notification interface {
	SetMessage(string)
	Close()
}
