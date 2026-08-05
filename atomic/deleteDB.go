package atomic

import (
	"errors"
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"path"
)

// 原子： 删除数据库操作
func AtomicDropDB(dbname string) error {
	atom_center := atomic_work_center.NewAtomicWorkCenter()
	task := atomic_work_center.NewTask(nil) // 创建一个作业

	remove_task_type, ok := pluginsloader.ResolveTaskType("misaka.removefile@用于删除文件的tasktype")
	if !ok {
		return errors.New("faild to get the tasktype of others plugins, please check the tasktype alias!")
	}
	task.TaskBooks = tasktype.NewShipBuilder().Add(
		remove_task_type, path.Join("./db-datas", dbname),
	).Build()
	// 添加删除数据库的TaskType 组合TaskShip，并交给原子中心获得对应的任务ID

	atom_center.AddTask(task, 3)
	return nil
}
