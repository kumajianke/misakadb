package share

import (
	"errors"
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	"misakadb/config"
	mson "misakadb/engine/Mson"
	engine_dispatch "misakadb/engine/dispatch"
	"misakadb/lock/global_lock"
	"misakadb/miusers"
	"misakadb/network/context"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"path/filepath"
)

func MiqlCreateDB(msonParse *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonParse.Active != "cre-dat" {
		return errors.New("Error Dispatch!")
	}

	engineName := msonParse.Engine                                    // 获取到对应的引擎名字
	dbEngine := engine_dispatch.NewEngine(engineName, msonParse.Name) // 数据库引擎

	if dbEngine == nil {
		clilog.Error("未知的引擎诉求")
		if err := serviceContext.Send("[err]未知的引擎诉求"); err != nil {
			return err
		}
		return errors.New("unknown engine request")
	}

	err := dbEngine.DBLoader().InitLoader(*msonParse) // 选择对应的数据库引擎进行初始化

	if err != nil {
		err_string := err.Error()
		serviceContext.Send("[err]" + err_string) // 错误信息的返回
		return err
	}

	serviceContext.Send("[ok]create db is ok!")
	return nil
}

func MiqlDropDB(msonPaese *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonPaese.Active != "drp-dat" {
		return errors.New("Error Dispatch!")
	}

	dbname := msonPaese.Name
	username := serviceContext.LoginUser

	// 鉴权操作
	err := miusers.NewUserManager().VerifyRole(username, "root")
	if err != nil {
		serviceContext.Send("[err]" + err.Error())
		return err
	}

	// 加锁操作
	_, unlock, err := global_lock.GetOrStoreGlobalLock(
		global_lock.GetLocksPrefix().DBFile+dbname,
		"l",
	)

	if err != nil {
		serviceContext.Send("[err]" + err.Error())
		return err
	}
	defer unlock()

	// 具体删除

	dropDBTaskType, ok := pluginsloader.ResolveTaskType("misaka.removefolder")
	if !ok {
		return errors.New("drop db task type is not registered")
	}

	work_center := atomic_work_center.NewAtomicWorkCenter() // 原子任务中心
	task_ship := tasktype.NewShipBuilder().Add(
		dropDBTaskType, filepath.Join(config.GetGlobalMisakaConfigure().Private.Storage.Path, dbname),
	).Build()

	remove_task := atomic_work_center.NewTask(task_ship)
	remove_task.TaskBooks = task_ship
	ok, task_key := work_center.AddTask(remove_task, 3)
	if !ok {
		return errors.New("create task failed")
	}
	if err, ok := work_center.DoSustain(task_key); err != nil || !ok {
		return err
	}
	serviceContext.Send("[ok]drop db is ok!")

	return nil
}
