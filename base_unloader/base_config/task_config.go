package baseconfig

import tasktype "misakadb/atomic/atomicWorkCenter/TaskType"

const (
	ModName = "misaka basic mode @misaka"
)
const (
	TaskRemoveFile   tasktype.TaskType = "remove_file"
	TaskModFile      tasktype.TaskType = "mod_file"
	TaskRemoveFolder tasktype.TaskType = "remove_folder"
)
