package global_lock

// building...

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/puzpuzpuz/xsync/v3"
)

var once sync.Once
var globalLocksPool *GlobalLocksPool

// 全局锁池
type GlobalLocks struct {
	RefCounter atomic.Int32
	Lock       sync.Mutex
}

type GlobalLocksPool struct {
	youngPool xsync.MapOf[string, *GlobalLocks] // 当前热门的锁
	oldPool   xsync.MapOf[string, *GlobalLocks] // 池化分代技术；即将淘汰的锁用于GC
}

func GetGlobalLockPool() *GlobalLocksPool {
	once.Do(func() {
		globalLocksPool = &GlobalLocksPool{}
	})
	return globalLocksPool
}

/**
* 获取一个锁，如果锁是第一次创建，自动注册到全局锁池
 */
func GetOrStoreGlobalLock(lock_name string, lock_method string) (*GlobalLocks, func(), error) {
	pools := GetGlobalLockPool()

	if lock_name == "" {
		return nil, nil, errors.New("lock_name can't be empty string!")
	}

	// 尝试直接在新池中获取锁
	lock, cache_hit := pools.youngPool.Load(lock_name)
	if !cache_hit {
		// 如果缓存没有命中 查看是否旧池中有这个锁 有的话获取对应的旧值并给他删掉
		lock, cache_hit = pools.oldPool.Load(lock_name)
		if cache_hit {
			lock, _ = pools.youngPool.LoadOrStore(lock_name, lock) // 升级锁
			pools.oldPool.Delete(lock_name)
		} else {
			lock, _ = pools.youngPool.LoadOrStore(lock_name, &GlobalLocks{})
			// 新旧池都没有 就新建一个
		}
	}

	lock_success := false
	switch lock_method {
	case "lock", "l":
		lock.Lock.Lock()
		lock.RefCounter.Add(1)
		lock_success = true
	case "tl", "try_lock":
		lock_success = lock.Lock.TryLock()
		if lock_success {
			lock.RefCounter.Add(1)
		} else {
			return nil, nil, errors.New("Get lock but lock is locked by other, can't lock the lock!")
		}
	default:
		return nil, nil, errors.New("invalid lock_method: must be 'lock', 'l', 'tl', or 'try_lock'")
	}

	return lock, func() {
		if lock_method == "" || !lock_success {
			return
		}
		lock.RefCounter.Add(-1) // 回收引用计数
		lock.Lock.Unlock()
	}, nil
}

func LocksPoolGCThrea() {

}
