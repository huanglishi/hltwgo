package ordermanag

import (
	"huling/app/model"

	"github.com/gohouse/gorose/v2"
)

func DB() gorose.IOrm {
	return model.DB.NewOrm()
}
