package base_unloader

import tasktype "misakadb/atomic/atomicWorkCenter/TaskType"

func AddTaskType(allTaskTypelst []tasktype.TaskType) []tasktype.TaskType {
	const (
		TaskRemoveFile   tasktype.TaskType = "remove_file"
		TaskModFile      tasktype.TaskType = "mod_file"
		TaskRemoveFolder tasktype.TaskType = "remove_folder"
	)
	allTaskTypelst = append(allTaskTypelst, TaskRemoveFile)
	allTaskTypelst = append(allTaskTypelst, TaskModFile)
	allTaskTypelst = append(allTaskTypelst, TaskRemoveFolder)
	return allTaskTypelst
}
