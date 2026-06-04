package global_lock

// building...

import (
	"sync"
	"sync/atomic"

	"github.com/puzpuzpuz/xsync/v3"
)

// 全局锁池
type GlobalLocks struct {
	RefCounter atomic.Int32
	Lock       sync.Mutex
}

type GlobalLocksPool struct {
	LocksPool xsync.MapOf[string, *GlobalLocks]
}
