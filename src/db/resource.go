package db

import (
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
)

func GetAllResources() (*entity.Resource, error) {
	session := con.Table(entity.ResourceRaw{})
	defer session.Close()
	var resources []entity.ResourceRaw
	err := session.Find(&resources)
	if err != nil {
		return nil, err
	}
	return util.FormResources(resources), nil
}

func GetResourceIdByAlias(alias string) (int, error) {
	session := con.Table(entity.ResourceRaw{})
	defer session.Close()
	var resource entity.ResourceRaw
	_, err := session.Where("alias = ?", alias).Get(&resource)
	if err != nil {
		return 0, err
	}
	return resource.Id, nil
}
