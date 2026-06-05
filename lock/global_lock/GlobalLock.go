package global_lock

// building...

import (
	"sync"
	"sync/atomic"

	"github.com/puzpuzpuz/xsync/v3"
)

var once sync.Once
var globalLockPools *GlobalLocksPool

// 全局锁池
type GlobalLocks struct {
	RefCounter atomic.Int32
	Lock       sync.Mutex
}

type GlobalLocksPool struct {
	YoungPool xsync.MapOf[string, *GlobalLocks] // 当前热门的锁
	OldPool   xsync.MapOf[string, *GlobalLocks] // 分代锁；即将淘汰的锁用于GC
}

func GetGlobalLockPool() *GlobalLocksPool {
	once.Do(func() {
		globalLockPools = &GlobalLocksPool{}
	})
	return globalLockPools
}

/**
* 获取一个锁，如果锁是第一次创建，自动注册到全局锁池
 */
func GetOrStoreGlobalLock(lock_name string, lockit bool) (*GlobalLocks, func()) {
	pools := GetGlobalLockPool()

	// 尝试直接在新池中获取锁
	lock, cache_hit := pools.YoungPool.Load(lock_name)
	if !cache_hit {
		// 如果缓存没有命中 查看是否旧池中有这个锁 有的话获取对应的旧值并给他删掉
		lock, cache_hit = pools.OldPool.Load(lock_name)
		if cache_hit {
			pools.YoungPool.LoadOrStore(lock_name, lock) // 升级锁
			pools.OldPool.Delete(lock_name)
		} else {
			pools.YoungPool.LoadOrStore(lock_name, &GlobalLocks{}) // 新旧池都没有 就新建一个
		}
	}
	//TODO
	return nil, nil
}
