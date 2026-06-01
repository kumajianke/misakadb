package share

import (
	"encoding/json"
	"misakadb/clilog"
	"misakadb/network/context"
)

type MsonParse struct {
	Active string `json:"active"`
	Name   string `json:"name"`
	Engine string `json:"engine"`
}

func NewMsonParse(json_string string) *MsonParse {
	msonParse := &MsonParse{}
	if err := json.Unmarshal([]byte(json_string), msonParse); err != nil {
		clilog.Error(
			"MSON 序列化错误： " + json_string,
		)
		return nil
	}

	return msonParse
}

func (msonParse *MsonParse) MsonToString() ([]byte, error) {
	data, err := json.Marshal(msonParse)
	if err != nil {
		clilog.Warning("无法解析的mson对象")
		return nil, err
	}
	return data, nil
}

func (msonParse *MsonParse) RunnableMsonFectory(service *context.ServiceConnContext) {
	active := msonParse.Active
	switch active {
	case "cre-dat":
		MiqlCreateDB(msonParse, service)
	default:
		service.Send("[err]unknow miql!")
	}

}
