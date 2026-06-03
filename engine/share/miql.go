package share

import (
	"errors"
	mson "misakadb/engine/Mson"
	"misakadb/network/context"
)

func RunMson(msonParse *mson.MsonParse, serviceContext *context.ServiceConnContext) error {
	if msonParse == nil {
		return errors.New("mson is nil")
	}

	switch msonParse.Active {
	case "cre-dat":
		return MiqlCreateDB(msonParse, serviceContext)
	case "drp-dat":
		return MiqlDropDB(msonParse, serviceContext)
	default:
		return serviceContext.Send("[err]unknown miql!")
	}
}
