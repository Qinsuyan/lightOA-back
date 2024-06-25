package db

import (
	"errors"
	"lightOA-end/src/entity"
	"strings"
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
func GetUnscopedUserRaw(user *entity.UserRaw) (bool, error) {
	return con.Unscoped().Get(user)
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
	userByPhone := &entity.UserRaw{
		Phone: online.Phone,
	}
	userByName := &entity.UserRaw{
		Username: online.Phone,
	}
	byPhoneExist, pErr := con.Get(userByPhone)
	byNameExist, eErr := con.Get(userByName)
	if pErr != nil && eErr != nil {
		return nil, errors.New("cannot get user")
	}
	if !byPhoneExist && !byNameExist {
		return nil, nil
	}
	if byPhoneExist {
		if !userByPhone.DeletedAt.IsZero() {
			return nil, nil
		}
		return userByPhone, nil
	} else {
		if !userByName.DeletedAt.IsZero() {
			return nil, nil
		}
		return userByName, nil
	}
}
func LoginUser(record *entity.Online) (string, error) {
	search := &entity.Online{
		Phone: record.Phone,
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
		_, err = session.Where("phone = ?", record.Phone).Cols("expire").Update(record)
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
func LogoutUserByToken(token string) {
	session := con.Table(entity.Online{})
	defer func() {
		session.Close()
	}()
	record := &entity.Online{Token: token, Expire: time.Now()}
	session.Where("token = ?", token).Update(record)
}
func LogoutUserByPhone(phone string) {
	session := con.Table(entity.Online{})
	defer func() {
		session.Close()
	}()
	record := &entity.Online{Phone: phone, Expire: time.Now()}
	session.Where("phone = ?", phone).Update(record)
}

var NaturalPriviledges = []string{"user:self"}

func contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
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
	roleRaw, err := GetRoleRawById(user.Role)
	if err != nil {
		return false, nil, err
	}
	if !roleRaw.DeletedAt.IsZero() {
		return false, nil, nil
	}
	if contains(NaturalPriviledges, alias) {
		return true, user, nil
	}
	//获取到用户的所有资源
	//如果连root资源都没有，则具有全部权限
	//否则需要判断用户是否具有该资源权限
	resourceIds := make([]*entity.RoleResource, 0)
	idSession := con.Table(entity.RoleResource{}).Where("roleId = ?", user.Role)
	defer func() {
		idSession.Close()
	}()
	err = idSession.Find(&resourceIds)
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

func ExistUserOfRole(roleId int) (bool, error) {
	session := con.Table(entity.UserRaw{})
	defer func() {
		session.Close()
	}()
	return session.Where("role = ?", roleId).Exist()
}

// AddUser 新增用户
func AddUser(user *entity.UserRaw) (bool, error) {
	// search := &entity.UserRaw{
	// 	Username: user.Username,
	// }
	// exist, err := con.Get(search)
	// if err != nil {
	// 	return true, err
	// }
	// if exist {
	// 	return false, nil
	// }
	search := &entity.UserRaw{
		Phone: user.Phone,
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

// 编辑用户
func EditUser(userId int, user *entity.UserRaw) error {
	_, err := con.Table(entity.UserRaw{}).Where("id = ?", userId).Update(user)
	return err
}

// 删除用户
func DeleteUser(userId int) error {
	_, err := con.Table(entity.UserRaw{}).Where("id = ?", userId).Delete(entity.UserRaw{})
	return err
}

func ListUser(filter *entity.UserListFilter) (int64, []entity.UserInfo, error) {
	session := con.Table(entity.UserRaw{})
	if filter.Username != "" {
		session.Where("username like ?", "%"+filter.Username+"%")
	}
	if filter.Phone != "" {
		session.Where("phone like ?", "%"+filter.Phone+"%")
	}
	if filter.Role != 0 {
		session.Where("role = ?", filter.Role)
	}
	if filter.Sort != "" && filter.Order != "" {
		cols := []string{}
		cols = append(cols, strings.Split(filter.Order, ",")...)
		if filter.Sort == "desc" {
			session.Desc(cols...)
		} else {
			session.Asc(cols...)
		}
	}
	if filter.PageSize != 0 && filter.PageNum != 0 {
		session.Limit(filter.PageSize, (filter.PageNum-1)*filter.PageSize)
	}
	users := []entity.UserRaw{}
	total, err := session.FindAndCount(&users)
	if err != nil {
		return 0, nil, err
	}
	var result []entity.UserInfo
	for _, user := range users {
		role, err := GetUserRoleByRoleId(user.Role)
		if err != nil {
			return 0, nil, err
		}
		result = append(result, entity.UserInfo{
			Department: user.Department,
			Role:       *role,
			Id:         user.Id,
			Username:   user.Username,
			Phone:      user.Phone,
		})
	}
	return total, result, nil
}
