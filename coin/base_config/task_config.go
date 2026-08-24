package baseconfig

import tasktype "misakadb/atomic/atomicWorkCenter/TaskType"

const (
	ModName = "misaka basic mode on v0.1.8@misaka"
)
const (
	TaskRemoveFile   tasktype.TaskType = "remove_file"
	TaskModFile      tasktype.TaskType = "mod_file"
	TaskRemoveFolder tasktype.TaskType = "remove_folder"
)
