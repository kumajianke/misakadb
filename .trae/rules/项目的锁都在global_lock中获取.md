本项目的锁不再单独创建 而是使用global_lock里的全局锁池进行获取，锁的命名方式是: `动作-资源[-参数]`如：
- `write-resource-name` 就是资源source为name的写锁
- `read-resource-name` 就是资源source为name的读锁
- `rw-resource-name` 就是资源source为name的读写锁
- `rw-resource` 就是资源source的读写锁
