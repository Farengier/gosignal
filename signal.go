package gosignal

import (
	"os"
	"os/signal"
	"sync"
)

var (
	once  sync.Once
	mtx   sync.RWMutex
	calls = make(map[os.Signal][]func())
	ch    = make(chan os.Signal, 1)
)

// OnSignal registers callback to be processed on signal
func OnSignal(sig os.Signal, callback func()) {
	signal.Notify(ch, sig)
	mtx.Lock()
	defer mtx.Unlock()
	calls[sig] = append(calls[sig], callback)
	once.Do(processor)
}

func processor() {
	go func() {
		for {
			sig := <-ch
			processSignal(sig)
		}
	}()
}

func processSignal(sig os.Signal) {
	mtx.RLock()
	defer mtx.RUnlock()
	if callbacks, ok := calls[sig]; ok {
		for _, call := range callbacks {
			call()
		}
	}
}
