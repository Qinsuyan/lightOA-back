package db

import (
	"errors"
	"lightOA-end/src/entity"
)

func AddFileInfo(file *entity.File) (int, error) {
	_, err := con.InsertOne(file)
	if err != nil {
		return 0, err
	}
	return file.Id, err
}

func DeleteFileInfo(ids []int) ([]string, error) {
	var files []entity.File
	err := con.In("id", ids).Find(&files)
	if err != nil {
		return nil, err
	}

	// 提取 UUID 字段
	var uuids []string
	for _, file := range files {
		uuids = append(uuids, file.UUID)
	}
	// 删除记录
	_, err = con.In("id", ids).Delete(&entity.File{})
	if err != nil {
		return nil, err
	}

	return uuids, nil
}

func CheckFileExist(ids []int) error {
	var count int64
	_, err := con.Table(&entity.Reimburse{}).In("id", ids).Count(&count)
	if err != nil {
		return nil
	}
	if count != int64(len(ids)) {
		return errors.New("false")
	}
	return nil
}

func GetFileInfoByIds(ids []int) ([]entity.File, error) {
	var files []entity.File
	err := con.Table(&entity.File{}).Cols("id", "name").In("id", ids).Find(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetFileById(id int) (*entity.File, error) {
	file := entity.File{
		Id: id,
	}
	_, err := con.Table(&entity.File{}).Get(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}
