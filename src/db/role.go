package db

import (
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
)

func GetUserRoleByRoleId(id int) (*entity.Role, error) {
	roleRaw := entity.RoleRaw{Id: id}
	roleExist, err := con.Get(&roleRaw)
	if err != nil {
		return nil, err
	}
	if !roleExist {
		return nil, nil
	}
	resourceIds := make([]entity.RoleResource, 0)
	idSession := con.Table(entity.RoleResource{})
	err = idSession.Where("roleId = ?", id).Find(&resourceIds)
	defer func() {
		idSession.Close()
	}()
	if err != nil {
		return nil, err
	}
	resources := make([]*entity.ResourceRaw, 0)
	session := con.Table(entity.ResourceRaw{})
	ids := make([]int, len(resources))
	for r := range resourceIds {
		ids = append(ids, resourceIds[r].ResourceId)
	}
	defer func() {
		session.Close()
	}()
	err = session.In("Id", ids).Find(&resources)
	if err != nil {
		return nil, err
	}
	//组装树结构
	role := &entity.Role{Id: roleRaw.Id, Name: roleRaw.Name}
	util.FormUserRole(role, resources)
	return role, nil
}

func GetRoleRawById(id int) (*entity.RoleRaw, error) {
	var roleRaw entity.RoleRaw
	has, err := con.Where("id = ?", id).Get(&roleRaw)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &roleRaw, nil
}

func IsRoleNameExist(name string) (bool, error) {
	var role entity.RoleRaw
	has, err := con.Where("name = ?", name).Get(&role)
	return has, err
}

func AddRole(role *entity.RoleRaw) (bool, error) {
	search := &entity.RoleRaw{
		Name: role.Name,
	}
	exist, err := con.Get(search)
	if err != nil {
		return true, err
	}
	if exist {
		return false, nil
	}
	_, err = con.InsertOne(role)
	if err != nil {
		return true, err
	}
	return true, nil
}

func EditRole(role *entity.RoleRaw) error {
	_, err := con.Where("id = ?", role.Id).Update(role)
	return err
}

func GetRoleIdByName(name string) (int, error) {
	var role entity.RoleRaw
	has, err := con.Where("name = ?", name).Get(&role)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, nil
	}
	return role.Id, nil
}

func AddRoleResource(roleId int, resource []*entity.Resource) (bool, error) {
	if len(resource) == 0 {
		return false, nil
	}
	for _, r := range resource {
		_, err := con.InsertOne(&entity.RoleResource{
			RoleId:     roleId,
			ResourceId: r.Id,
		})
		if err != nil {
			return false, err
		}
		if r.Children != nil {
			ok, err := AddRoleResource(roleId, r.Children)
			if !ok || err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func DeleteRoleResource(roleId int) error {
	_, err := con.Where("roleId = ?", roleId).Delete(&entity.RoleResource{})
	return err
}

func DeleteRole(id int) error {
	_, err := con.Delete(&entity.RoleRaw{Id: id})
	return err
}
func ListRole(filter *entity.RoleListFilter) ([]entity.Role, error) {
	session := con.Table(entity.RoleRaw{})
	if filter.Name != "" {
		session.Where("name like ?", "%"+filter.Name+"%")
	}

	if filter.Sort != "" {
		if filter.Sort == "desc" {
			session.Desc("name")
		} else {
			session.Asc("name")
		}
	}
	if filter.PageSize != 0 && filter.PageNum != 0 {
		session.Limit(filter.PageSize, (filter.PageNum-1)*filter.PageSize)
	}
	roleRaws := []entity.RoleRaw{}
	if err := session.Find(&roleRaws); err != nil {
		return nil, err
	}
	result := []entity.Role{}
	for _, r := range roleRaws {
		role, err := GetUserRoleByRoleId(r.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, *role)
	}
	return result, nil
}
