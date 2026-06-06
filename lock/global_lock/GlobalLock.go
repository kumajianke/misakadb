package global_lock

// building...

import (
	"errors"
	"misakadb/clilog"
	"sync"
	"sync/atomic"
	"time"

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
		globalLocksPool = &GlobalLocksPool{
			youngPool: *xsync.NewMapOf[string, *GlobalLocks](),
			oldPool:   *xsync.NewMapOf[string, *GlobalLocks](),
		}
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
Redo:
	// 尝试直接在新池中获取锁
	lock, cache_hit := pools.youngPool.Load(lock_name)
	var lockRefCounter int32
	if cache_hit {
		lockRefCounter = lock.RefCounter.Load()
		if lockRefCounter < 0 {
			// 已经被标记到了墓碑 晦气晦气!
			cache_hit = false
		}
	}
	if !cache_hit {
		// 如果缓存没有命中 查看是否旧池中有这个锁 有的话获取对应的旧值并给他删掉
		lock, cache_hit = pools.oldPool.Load(lock_name)
		if cache_hit {
			lockRefCounter = lock.RefCounter.Load()
			if lockRefCounter < 0 {
				cache_hit = false
			} else {
				lock, _ = pools.youngPool.LoadOrStore(lock_name, lock) // 升级锁
				pools.oldPool.Delete(lock_name)
			}
		}
		if !cache_hit {
			lock, _ = pools.youngPool.LoadOrStore(lock_name, &GlobalLocks{})
			// 新旧池都没有 就新建一个
		}
	}

	lock_success := false
	switch lock_method {
	case "lock", "l":
		if !lock.RefCounter.CompareAndSwap(lockRefCounter, lockRefCounter+1) {
			goto Redo
		} else {
			lock.Lock.Lock()
		}
		lock_success = true
	case "tl", "try_lock":
		if !lock.RefCounter.CompareAndSwap(lockRefCounter, lockRefCounter+1) {
			// 我刚提的新车!!(获取的锁被其他协程用了)
			goto Redo
		} else {
			lock_success = lock.Lock.TryLock()
			if !lock_success {
				lock.RefCounter.Add(-1) // 美美把玩 CAS +1 需要注意碰到墓碑 防止被中途回收  -1 要是给的墓碑就当随份子了
			}
		}

	default:
		return nil, nil, errors.New("invalid lock_method: must be 'lock', 'l', 'tl', or 'try_lock'")
	}

	return lock, func() {
		if lock_method == "" || !lock_success {
			return
		}
		lock.Lock.Unlock()      // 先解锁 防止Add - 1为-1被回收了
		lock.RefCounter.Add(-1) // 回收引用计数
	}, nil
}

func lockPoolsGCThread() {
	// 因为 old 有回到 young的机会 所以我们先加紧处理old的淘汰
	// 再去做 young 的 old 化
	pools := GetGlobalLockPool()
	pools.oldPool.Range(func(key string, value *GlobalLocks) bool {
		lock_ref := value.RefCounter.Load()
		if lock_ref == 0 { // 引用计数为0表示已经被解锁了
			pools.oldPool.Delete(key) // 老东西 你已经没用了!
		}
		return true
	})

	pools.youngPool.Range(func(key string, value *GlobalLocks) bool {
		if value.RefCounter.CompareAndSwap(0, -1) {
			_, loaded := pools.oldPool.LoadOrStore(key, value)
			if loaded {
				clilog.Warning("[GlobalLocksPool] before that Young Pool GC the lock，old pool had the lock even.")
			}
			pools.youngPool.Delete(key) // 你也曾经是一个年轻的锁 直到膝盖中了一箭
		}
		return true
	})
}

func StartLoPoolGC() {
	go func() {
		for {
			lockPoolsGCThread()
			time.Sleep(time.Second * 50)
		}
	}()
}
