package engine

/**
* 数据库日志核心，用于记录每个记录指定存储等信息
**/
type BaseLoaderCore struct {
	Path string // dbdata的存储路径
}

/**
 * 数据库核心 不同的数据库指向了一个核心
 */
type BaseEngineCore struct {
	DBLoader BaseLoaderCore
}
