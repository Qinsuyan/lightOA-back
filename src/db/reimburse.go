package db

import (
	"lightOA-end/src/entity"
	"strings"
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

func DeleteReimburseFile(fileIds []int, reimburseId int) error {
	_, err := con.Table(entity.ReimbuiseFiles{}).In("fileId", fileIds).Delete(&entity.ReimbuiseFiles{})
	return err
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

func EditReimburse(reimburse *entity.Reimburse) error {
	_, err := con.Where("id = ?", reimburse.Id).Update(reimburse)
	return err
}

func DeleteReimburse(id int) error {
	_, err := con.Where("id = ?", id).Delete(&entity.Reimburse{})
	return err
}

func GetReimburseHistoryList(userId int, filter *entity.ReimburseListFilter) (int64, []entity.ReimburseInfo, error) {
	session := con.Table(entity.Reimburse{})
	session.Where("createdBy = ? or applicant = ?", userId, userId)
	if filter.Title != "" {
		session.Where("title like ?", "%"+filter.Title+"%")
	}
	if filter.MaxAmount > 0 {
		session.Where("amount <= ?", filter.MaxAmount)
	}
	if filter.MinAmount > 0 {
		session.Where("amount >= ?", filter.MinAmount)
	}
	if filter.Auditor > 0 {
		session.Where("auditor = ?", filter.Auditor)
	}
	if filter.Status > 0 {
		session.Where("status = ?", filter.Status)
	}
	if filter.Applicant > 0 {
		session.Where("applicant = ?", filter.Applicant)
	}
	if !filter.HappenedAtStart.IsZero() {
		session.Where("happenedAt >= ?", filter.HappenedAtStart)
	}
	if !filter.HappenedAtEnd.IsZero() {
		session.Where("happenedAt <= ?", filter.HappenedAtStart)
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
	if filter.PageNum > 0 && filter.PageSize > 0 {
		session.Limit(filter.PageSize, (filter.PageNum-1)*filter.PageSize)
	}
	var reimburseList []entity.Reimburse
	total, err := session.FindAndCount(&reimburseList)
	if err != nil {
		return 0, nil, err
	}
	var reimburseInfoList []entity.ReimburseInfo
	for _, reimburse := range reimburseList {
		r := entity.ReimburseInfo{
			Id:         reimburse.Id,
			Title:      reimburse.Title,
			Desc:       reimburse.Desc,
			Amount:     reimburse.Amount,
			Status:     reimburse.Status,
			HappenedAt: reimburse.HappenedAt,
			UpdatedAt:  reimburse.UpdatedAt,
			CreatedAt:  reimburse.CreatedAt,
		}
		//补充auditor
		if reimburse.Auditor != 0 {
			auditor := &entity.UserRaw{
				Id: reimburse.Auditor,
			}
			has, err := GetUnscopedUserRaw(auditor)
			if !has || err != nil {
				return 0, nil, err
			}
			r.Auditor = entity.UserInfo{
				Id:         auditor.Id,
				Username:   auditor.Username,
				Department: auditor.Department,
			}
		}
		//补充createdBy
		creator := &entity.UserRaw{Id: reimburse.CreatedBy}
		has, err := GetUnscopedUserRaw(creator)
		if !has || err != nil {
			return 0, nil, err
		}
		r.CreatedBy = entity.UserInfo{
			Id:         creator.Id,
			Username:   creator.Username,
			Department: creator.Department,
		}
		//补充applicant
		applicant := &entity.UserRaw{Id: reimburse.Applicant}
		has, err = GetUnscopedUserRaw(applicant)
		if !has || err != nil {
			return 0, nil, err
		}
		r.Applicant = entity.UserInfo{
			Id:         applicant.Id,
			Username:   applicant.Username,
			Department: applicant.Department,
		}
		//补充files
		fileIds, err := GetReimburseFileIdsByReimburseId(reimburse.Id)
		if err != nil {
			return 0, nil, err
		}
		files, err := GetFileInfoByIds(fileIds)
		if err != nil {
			return 0, nil, err
		}
		var fileInfos []entity.FileInfo
		for _, file := range files {
			fileInfos = append(fileInfos, entity.FileInfo{
				Id:   file.Id,
				Name: file.Name,
			})
		}
		r.Files = fileInfos
		//补充comments
		comments, err := GetReimburseCommentsByReimburseId(reimburse.Id)
		if err != nil {
			return 0, nil, err
		}
		r.Comments = comments
		reimburseInfoList = append(reimburseInfoList, r)
	}
	return total, reimburseInfoList, nil
}

func GetReimburseAllList(filter *entity.ReimburseListFilter) (int64, []entity.ReimburseInfo, error) {
	session := con.Table(entity.Reimburse{})
	if filter.Title != "" {
		session.Where("title like ?", "%"+filter.Title+"%")
	}
	if filter.MaxAmount > 0 {
		session.Where("amount <= ?", filter.MaxAmount)
	}
	if filter.MinAmount > 0 {
		session.Where("amount >= ?", filter.MinAmount)
	}
	if filter.Auditor > 0 {
		session.Where("auditor = ?", filter.Auditor)
	}
	if filter.CreatedBy > 0 {
		session.Where("createdBy = ?", filter.CreatedBy)
	}
	if filter.Status > 0 {
		session.Where("status = ?", filter.Status)
	}
	if filter.Applicant > 0 {
		session.Where("applicant = ?", filter.Applicant)
	}
	if !filter.HappenedAtStart.IsZero() {
		session.Where("happenedAt >= ?", filter.HappenedAtStart)
	}
	if !filter.HappenedAtEnd.IsZero() {
		session.Where("happenedAt <= ?", filter.HappenedAtStart)
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
	if filter.PageNum > 0 && filter.PageSize > 0 {
		session.Limit(filter.PageSize, (filter.PageNum-1)*filter.PageSize)
	}
	var reimburseList []entity.Reimburse
	total, err := session.FindAndCount(&reimburseList)
	if err != nil {
		return 0, nil, err
	}
	var reimburseInfoList []entity.ReimburseInfo
	for _, reimburse := range reimburseList {
		r := entity.ReimburseInfo{
			Id:         reimburse.Id,
			Title:      reimburse.Title,
			Desc:       reimburse.Desc,
			Amount:     reimburse.Amount,
			Status:     reimburse.Status,
			HappenedAt: reimburse.HappenedAt,
			UpdatedAt:  reimburse.UpdatedAt,
			CreatedAt:  reimburse.CreatedAt,
		}
		//补充auditor
		if reimburse.Auditor != 0 {
			auditor := &entity.UserRaw{
				Id: reimburse.Auditor,
			}
			has, err := GetUnscopedUserRaw(auditor)
			if !has || err != nil {
				return 0, nil, err
			}
			r.Auditor = entity.UserInfo{
				Id:         auditor.Id,
				Username:   auditor.Username,
				Department: auditor.Department,
			}
		}
		//补充createdBy
		creator := &entity.UserRaw{Id: reimburse.CreatedBy}
		has, err := GetUnscopedUserRaw(creator)
		if !has || err != nil {
			return 0, nil, err
		}
		r.CreatedBy = entity.UserInfo{
			Id:         creator.Id,
			Username:   creator.Username,
			Department: creator.Department,
		}
		//补充applicant
		applicant := &entity.UserRaw{Id: reimburse.Applicant}
		has, err = GetUnscopedUserRaw(applicant)
		if !has || err != nil {
			return 0, nil, err
		}
		r.Applicant = entity.UserInfo{
			Id:         applicant.Id,
			Username:   applicant.Username,
			Department: applicant.Department,
		}
		//补充files
		fileIds, err := GetReimburseFileIdsByReimburseId(reimburse.Id)
		if err != nil {
			return 0, nil, err
		}
		files, err := GetFileInfoByIds(fileIds)
		if err != nil {
			return 0, nil, err
		}
		var fileInfos []entity.FileInfo
		for _, file := range files {
			fileInfos = append(fileInfos, entity.FileInfo{
				Id:   file.Id,
				Name: file.Name,
			})
		}
		r.Files = fileInfos
		//补充comments
		comments, err := GetReimburseCommentsByReimburseId(reimburse.Id)
		if err != nil {
			return 0, nil, err
		}
		r.Comments = comments
		reimburseInfoList = append(reimburseInfoList, r)
	}
	return total, reimburseInfoList, nil
}

func AuditReimburse(audition *entity.ReimburseAudition, auditor int) (bool, error) {
	reimburse, err := GetReimburseById(audition.Id)
	if err != nil {
		return false, err
	}
	if reimburse.Status == 3 {
		return false, nil
	}
	if audition.Approve {
		reimburse.Status = 3
	} else {
		reimburse.Status = 2
	}
	reimburse.Auditor = auditor
	_, err = con.Table(entity.Reimburse{}).Update(reimburse)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ReapplyReimburse(id int) (bool, error) {
	reimburse, err := GetReimburseById(id)
	if err != nil {
		return false, err
	}
	if reimburse.Status != 2 {
		return false, nil
	}
	reimburse.Status = 1
	_, err = con.Table(entity.Reimburse{}).Update(reimburse)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetReimburseAuditors() ([]entity.UserInfo, error) {
	resourceId, err := GetResourceIdByAlias("reimburse:audit")
	if err != nil {
		return nil, err
	}
	if resourceId == 0 {
		return nil, nil
	}
	roles, err := GetRoleIdsByResourceId(resourceId)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	var users []entity.UserInfo
	err = con.Table(entity.UserRaw{}).In("role", roles).Find(&users)
	if err != nil {
		return nil, err
	}
	if users == nil {
		return nil, nil
	}
	var userList []entity.UserInfo
	for _, user := range users {
		userList = append(userList, entity.UserInfo{
			Id:         user.Id,
			Username:   user.Username,
			Department: user.Department,
		})
	}
	return userList, nil
}
