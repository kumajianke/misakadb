package runner

import work_center_serializer "misakadb/atomic/atomicWorkCenter/WorkCenterSerializer"

func SaveCenter() {
	work_center_serializer.FastInitWorkCenterSerializer()
}
