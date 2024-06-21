package db

import (
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
)

func GetAllResources() (*entity.Resource, error) {
	session := con.Table(entity.ResourceRaw{})
	defer session.Close()
	var resources []*entity.ResourceRaw
	err := session.Find(&resources)
	if err != nil {
		return nil, err
	}
	return util.FormResources(resources), nil
}
