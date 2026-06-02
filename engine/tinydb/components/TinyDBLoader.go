package components

import (
	mson "misakadb/engine/Mson"
	engine_base "misakadb/engine/base"
)

/**
* 记录器对象 用来JSON序列化
 */
type TinyDBRecorder struct {
}

type TinyDBLoaderImp struct {
	engine_base.BaseLoaderCore
}

func (this *TinyDBLoaderImp) WriteLoader(log mson.MsonParse) error {
	return nil
}

func (this *TinyDBLoaderImp) ReadLoader(log mson.MsonParse) error {
	return nil
}

func (this *TinyDBLoaderImp) InitLoader(log mson.MsonParse) error {
	return nil
}
