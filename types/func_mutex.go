package types

import "sync"

type FuncMutex struct {
	fn    func()
	mutex sync.Mutex
}

func (fm *FuncMutex) Set(fn func()) {
	fm.mutex.Lock()
	fm.fn = fn
	fm.mutex.Unlock()
}

func (fm *FuncMutex) Call() bool {
	fm.mutex.Lock()
	fn := fm.fn
	fm.mutex.Unlock()
	if fn != nil {
		fn()
	}
	return fn != nil
}
