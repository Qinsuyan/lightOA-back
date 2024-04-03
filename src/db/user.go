package db

import (
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
	"time"
)

func IsUserOnline(token string) (bool, error) {
	pack := &entity.Online{
		Token: token,
	}
	exist, err := con.Get(pack)
	if err != nil {
		return false, err
	}
	if exist {
		if pack.Expire.After(time.Now()) {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

func GetUserRaw(user *entity.UserRaw) (bool, error) {
	return con.Get(user)
}

func GetUserRawByToken(token string) (*entity.UserRaw, error) {
	online := &entity.Online{
		Token: token,
	}
	_, err := con.Get(online)
	//外层已经确保记录存在
	if err != nil {
		return nil, err
	}
	user := &entity.UserRaw{
		Username: online.Username,
	}
	exist, err := con.Get(user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	if !user.DeletedAt.IsZero() {
		return nil, nil
	}
	return user, nil
}

func LoginUser(record *entity.Online) (string, error) {
	search := &entity.Online{
		Username: record.Username,
	}
	exist, err := con.Get(search)
	if err != nil {
		return "", err
	}
	if exist {
		session := con.Table(entity.Online{})
		defer func() {
			session.Close()
		}()
		_, err = session.Where("username = ?", record.Username).Cols("expire").Update(record)
		if err != nil {
			return "", err
		}
		return search.Token, nil
	} else {
		_, err = con.InsertOne(record)
		if err != nil {
			return "", err
		}
		return record.Token, nil
	}
}

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

func LogoutUser(record *entity.Online) {
	record.Expire = time.Now()
	session := con.Table(entity.Online{})
	defer func() {
		session.Close()
	}()
	session.Where("token = ?", record.Token).Update(record)
}

func IsUserAuthorized(alias string, token string) (bool, *entity.UserRaw, error) {
	user, err := GetUserRawByToken(token)
	if err != nil {
		return false, nil, err
	}
	if user == nil {
		return false, nil, nil
	}
	if !user.DeletedAt.IsZero() {
		return false, nil, nil
	}
	//获取到用户的所有资源
	//如果连root资源都没有，则具有全部权限
	//否则需要判断用户是否具有该资源权限
	resourceIds := make([]*entity.RoleResource, 0)
	idSession := con.Where("roleId = ?", user.Role)
	defer func() {
		idSession.Close()
	}()
	err = idSession.Find(resourceIds)
	if err != nil {
		return false, nil, err
	}
	if len(resourceIds) == 0 {
		return true, user, nil
	}
	resources := make([]*entity.ResourceRaw, 0)
	session := con.Table(entity.ResourceRaw{})
	defer func() {
		idSession.Close()
	}()
	ids := make([]int, len(resourceIds))
	for r := range resourceIds {
		ids = append(ids, resourceIds[r].ResourceId)
	}
	err = session.In("Id", ids).Find(&resources)
	if err != nil {
		return false, nil, err
	}
	for r := range resources {
		if resources[r].Alias == alias {
			return true, user, nil
		}
	}
	return false, user, nil
}

// AddUser 新增用户
func AddUser(user *entity.UserRaw) (bool, error) {
	search := &entity.UserRaw{
		Username: user.Username,
	}
	exist, err := con.Get(search)
	if err != nil {
		return true, err
	}
	if exist {
		return false, nil
	}
	_, err = con.InsertOne(user)
	if err != nil {
		return true, err
	}
	return true, nil
}
