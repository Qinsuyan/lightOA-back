package db

import "lightOA-end/src/entity"

func AddDepartment(dep *entity.Department) (bool, error) {
	search := &entity.Department{
		Name: dep.Name,
	}
	exist, err := con.Get(search)
	if err != nil {
		return true, err
	}
	if exist {
		return false, nil
	}
	_, err = con.InsertOne(dep)
	if err != nil {
		return true, err
	}
	return true, nil
}

func DeleteDepartment(id int) error {
	_, err := con.Delete(&entity.Department{
		Id: id,
	})
	return err
}

func ExistUserOfDepartment(depId int) (bool, error) {
	session := con.Table(entity.UserRaw{})
	defer func() {
		session.Close()
	}()
	return session.Where("department = ?", depId).Exist()
}

func EditDepartment(dep *entity.Department) error {
	_, err := con.Update(dep, &entity.Department{
		Id: dep.Id,
	})
	return err
}

func ListDepartment(filter *entity.ListRequest) (int64, []entity.Department, error) {
	var list []entity.Department
	session := con.Table(entity.Department{})
	if filter.PageNum != 0 && filter.PageSize != 0 {
		session.Limit(filter.PageSize, (filter.PageNum-1)*filter.PageSize)
	}
	total, err := session.FindAndCount(&list)
	return total, list, err
}

func GetDepartmentById(id int) (*entity.Department, error) {
	var dep entity.Department
	has, err := con.Where("id = ?", id).Get(&dep)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &dep, nil
}
