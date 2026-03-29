package errs

import (
	"log"
	"sync"
)

var (
	recoverHandleMu sync.RWMutex
	h               RecoverHandle
)

type RecoverHandle func(any, *Stack)

func defaultRecoverHandle(err any, stack *Stack) {
	recoverHandleMu.RLock()
	handler := h
	recoverHandleMu.RUnlock()

	if handler != nil {
		handler(err, stack)
		return
	}
	log.Printf("panic: %v %+v\n", err, stack)
}

// Recover recover 处理，打印堆栈
// 直接 defer errs.Recover() 而不能defer func(){errs.Recover()}
func Recover() {
	RecoverFunc(defaultRecoverHandle)
}

func RecoverFunc(fn RecoverHandle) {
	err := recover()
	if err == nil {
		return
	}
	stack := Callers(2, cMaxStackDepth)
	fn(err, stack)
}

// RegisterRecoverHandle 自定义的recover处理函数
func RegisterRecoverHandle(fn RecoverHandle) {
	recoverHandleMu.Lock()
	defer recoverHandleMu.Unlock()
	h = fn
}
