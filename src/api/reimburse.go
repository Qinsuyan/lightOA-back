package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/storage"
	"strconv"

	"github.com/labstack/echo/v4"
)

func handleReimburseAddFile(c echo.Context) error {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if f.Size == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "文件不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if f.Size > 1024*1024*50 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "文件大小不能超过50M",
			Prompt: entity.ERROR,
		})
		return nil
	}
	name := c.FormValue("name")
	if name == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "文件名不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	reimburse := c.FormValue("reimburse")
	if reimburse == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "未传入申报信息ID",
			Prompt: entity.ERROR,
		})
		return nil
	}
	reimburseId, err := strconv.Atoi(reimburse)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}

	//检查权限
	reimburseInfo, err := db.GetReimburseById(reimburseId)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburseInfo == nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "未找到报销信息",
			Prompt: entity.ERROR,
		})
		return nil
	}
	user, err := getUser(c)
	if err != nil {
		return err
	}
	if user == nil {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "登录状态已失效",
			Prompt: entity.ERROR,
		})
		return nil
	}
	//自己是createdBy或applicant，且状态不是已经审核，才能上传
	authed := false
	if reimburseInfo.Applicant == user.Id || reimburseInfo.CreatedBy == user.Id {
		if reimburseInfo.Status != 3 {
			authed = true
		}
	}
	if !authed {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限为此报销信息添加文件",
			Prompt: entity.ERROR,
		})
		return nil
	}
	uuid, err := storage.SaveFile(f, map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
	})
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "文件保存失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if uuid == "invalid format" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "仅支持上传PDF或图片",
			Prompt: entity.ERROR,
		})
		return nil
	}
	//保存到数据库
	id, err := db.AddFileInfo(&entity.File{Name: name, UUID: uuid, Owner: user.Id, Type: 2})
	if err != nil || id == 0 {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "文件保存失败",
			Prompt: entity.ERROR,
		})
		return err
	}

	//与reimburse关联
	err = db.AddReimburseFile([]int{id}, reimburseId)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "添加失败",
			Prompt: entity.ERROR,
		})
		return err
	}

	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "文件保存成功",
		Prompt: entity.SUCCESS,
		Data:   id,
	})
	return nil
}

func handleReimburseAdd(c echo.Context) error {
	authed, user, err := checkAuth(c, "reimburse:add")
	if err != nil || !authed {
		return err
	}
	reimburse := &entity.ReimbursePayload{}
	err = c.Bind(reimburse)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.Title == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "标题不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.Amount == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "金额不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.Desc == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "描述不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.HappenedAt.IsZero() {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "发生时间不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.Applicant == 0 || reimburse.Applicant == user.Id {
		//未传入申请者时，用户自己是申请者
		reimburse.Applicant = user.Id
	} else {
		authed, _, err := getAuth(c, "reimburse:audit")
		if err != nil {
			return err
		}
		if !authed {
			c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
				Code:   ERROR_AUTH,
				Msg:    "无权限为他人申报",
				Prompt: entity.ERROR,
			})
			return nil
		}
	}

	//添加到reimburse表
	reimburseTableData := entity.Reimburse{
		Applicant:  reimburse.Applicant,
		Amount:     reimburse.Amount,
		Auditor:    0,
		Title:      reimburse.Title,
		Desc:       reimburse.Desc,
		Status:     1,
		HappenedAt: reimburse.HappenedAt,
		CreatedBy:  user.Id,
	}
	reimburseId, err := db.AddReimburse(&reimburseTableData)
	if err != nil || reimburseId == 0 {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "添加失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[int]{
		Code:   OK,
		Msg:    "添加成功",
		Prompt: entity.SUCCESS,
		Data:   reimburseId,
	})
	return nil
}

func handleReimburseDetail(c echo.Context) error {
	//审核人员或者申请者可以查看详情
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	reimburse, err := db.GetReimburseById(id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "无法获取报销信息",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse == nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "报销信息不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	authed, user, err := getAuth(c, "reimburse:audit")
	if (reimburse.Applicant != user.Id && reimburse.CreatedBy != user.Id) && !authed {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.WARN,
		})
		return err
	}
	applicant := &entity.UserRaw{Id: reimburse.Applicant}
	exist, err := db.GetUnscopedUserRaw(applicant)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取申请人信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if !exist {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "申请人不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	creator := &entity.UserRaw{Id: reimburse.CreatedBy}
	if reimburse.CreatedBy != reimburse.Applicant {
		exist, err := db.GetUnscopedUserRaw(creator)
		if err != nil {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "获取创建人信息失败",
				Prompt: entity.ERROR,
			})
			return err
		}
		if !exist {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "创建人不存在",
				Prompt: entity.ERROR,
			})
			return nil
		}
	} else {
		creator = applicant
	}

	response := entity.ReimburseInfo{
		Id:     reimburse.Id,
		Title:  reimburse.Title,
		Desc:   reimburse.Desc,
		Amount: reimburse.Amount,
		Applicant: entity.UserInfo{
			Id:         applicant.Id,
			Username:   applicant.Username,
			Department: applicant.Department,
		},
		CreatedBy: entity.UserInfo{
			Id:         creator.Id,
			Username:   creator.Username,
			Department: creator.Department,
		},
		Status:     reimburse.Status,
		HappenedAt: reimburse.HappenedAt,
		CreatedAt:  reimburse.CreatedAt,
		UpdatedAt:  reimburse.UpdatedAt,
	}
	//添加auditor,files,comments
	if reimburse.Auditor != 0 {
		auditor := &entity.UserRaw{
			Id: reimburse.Auditor,
		}
		exist, err := db.GetUnscopedUserRaw(auditor)
		if err != nil {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "获取审核人信息失败",
				Prompt: entity.ERROR,
			})
			return err
		}
		if !exist {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "审核人不存在",
				Prompt: entity.ERROR,
			})
			return nil
		}
		response.Auditor = entity.UserInfo{
			Id:         auditor.Id,
			Username:   auditor.Username,
			Department: auditor.Department,
		}
	}
	fileIds, err := db.GetReimburseFileIdsByReimburseId(reimburse.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取报销文件失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if len(fileIds) > 0 {
		files, err := db.GetFileInfoByIds(fileIds)
		if err != nil {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "获取报销文件失败",
				Prompt: entity.ERROR,
			})
			return err
		}
		var fileInfos []entity.FileInfo
		for _, file := range files {
			fileInfos = append(fileInfos, entity.FileInfo{
				Id:   file.Id,
				Name: file.Name,
			})
		}
		response.Files = fileInfos
	}
	comments, err := db.GetReimburseCommentsByReimburseId(reimburse.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取报销审核信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	response.Comments = comments
	c.JSON(OK, entity.HttpResponse[entity.ReimburseInfo]{
		Code:   OK,
		Msg:    "获取报销信息成功",
		Prompt: entity.SILENT,
		Data:   response,
	})
	return nil
}

func handleReimburseEdit(c echo.Context) error {

	reimburse := &entity.ReimbursePayload{}
	err := c.Bind(reimburse)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.Id == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.Amount == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "报销金额不能为0",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.Title == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "标题不能为空",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.Desc == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "描述不能为空",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse.HappenedAt.IsZero() {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "发生时间不能为空",
			Prompt: entity.ERROR,
		})
		return err
	}
	user, err := getUser(c)
	if err != nil {
		return err
	}
	if user == nil {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "登录状态已失效",
			Prompt: entity.ERROR,
		})
		return nil
	}
	reimburseInfo, err := db.GetReimburseById(reimburse.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取报销信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburseInfo == nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "报销信息不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburseInfo.Status == 3 {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限修改已审核通过的报销信息",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if user.Id != reimburseInfo.Applicant && user.Id != reimburseInfo.CreatedBy {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.Applicant != 0 && reimburseInfo.Applicant != reimburse.Applicant {
		//修改申请人
		auth, _, err := checkAuth(c, "reimburse:audit")
		if err != nil {
			return err
		}
		if !auth {
			c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
				Code:   ERROR_AUTH,
				Msg:    "无权限修改申请人",
				Prompt: entity.ERROR,
			})
			return nil
		}
	}
	err = db.EditReimburse(&entity.Reimburse{
		Id:         reimburse.Id,
		Applicant:  reimburse.Applicant,
		Amount:     reimburse.Amount,
		Status:     reimburseInfo.Status,
		Auditor:    reimburseInfo.Auditor,
		Title:      reimburse.Title,
		Desc:       reimburse.Desc,
		HappenedAt: reimburse.HappenedAt,
		CreatedBy:  reimburseInfo.CreatedBy,
	})
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "修改报销信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "修改成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleReimburseDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	reimburse, err := db.GetReimburseById(id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除报销信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if reimburse == nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "报销信息不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if reimburse.Status == 3 {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限删除已审核通过的报销信息",
			Prompt: entity.ERROR,
		})
		return nil
	}
	user, err := getUser(c)
	if err != nil {
		return err
	}
	if user == nil {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "登录状态已失效",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if user.Id != reimburse.Applicant && user.Id != reimburse.CreatedBy {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.ERROR,
		})
		return nil
	}
	err = db.DeleteReimburse(reimburse.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除失败",
			Prompt: entity.ERROR,
		})
		return nil
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "删除成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleReimburseHistory(c echo.Context) error {
	var filter entity.ReimburseListFilter
	err := c.Bind(&filter)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	user, err := getUser(c)
	if err != nil {
		return err
	}
	if user == nil {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "登录状态已失效",
			Prompt: entity.ERROR,
		})
		return nil
	}
}
