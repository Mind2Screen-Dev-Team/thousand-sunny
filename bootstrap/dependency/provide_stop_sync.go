package dependency

import "sync"

type WaitGroupStop struct {
	sync.WaitGroup
}

func ProvideSyncStop() *WaitGroupStop {
	return &WaitGroupStop{}
}
