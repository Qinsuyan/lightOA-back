package db

import (
	"lightOA-end/src/entity"
)

func AddReimburse(reimburse *entity.Reimburse) (int, error) {
	_, err := con.InsertOne(reimburse)
	return reimburse.Id, err
}

func AddReimburseFile(fileIds []int, reimburseId int) error {
	for _, id := range fileIds {
		_, err := con.InsertOne(&entity.ReimbuiseFiles{ReimburseId: reimburseId, FileId: id})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetReimburseById(id int) (*entity.Reimburse, error) {
	search := &entity.Reimburse{Id: id}
	if has, err := con.Get(search); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return search, nil
}

func GetReimburseFileIdsByReimburseId(id int) ([]int, error) {
	var list []entity.ReimbuiseFiles

	err := con.Where("reimburseId = ?", id).Find(&list)

	if err != nil {
		return nil, err
	}
	var ids []int
	for _, item := range list {
		ids = append(ids, item.FileId)
	}
	return ids, nil
}

func GetReimburseCommentsByReimburseId(id int) ([]entity.ReimburseComment, error) {
	var list []entity.ReimburseComments

	err := con.Where("reimburseId = ?", id).Find(&list)

	if err != nil {
		return nil, err
	}
	var comments []entity.ReimburseComment
	for _, c := range list {
		newComment := entity.ReimburseComment{
			Id:      c.Id,
			Comment: c.Comment,
			Time:    c.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		creator := &entity.UserRaw{Id: c.Creator}
		exist, err := GetUnscopedUserRaw(creator)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, nil
		}
		newComment.Creator = entity.UserInfo{
			Id:         creator.Id,
			Username:   creator.Username,
			Department: creator.Department,
		}
		comments = append(comments, newComment)
	}
	return comments, nil
}
