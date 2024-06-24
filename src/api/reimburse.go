package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/storage"
	"strconv"

	"github.com/labstack/echo/v4"
)

func handleReimburseAddFile(c echo.Context) error {
	authed, user, err := checkAuth(c, "reimburse:add")
	if err != nil || !authed {
		return err
	}
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
	if len(reimburse.Files) == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "文件不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	//检查每一个文件是不是都在数据库中
	err = db.CheckFileExist(reimburse.Files)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "文件不存在",
			Prompt: entity.ERROR,
		})
		return err
	}
	//添加到reimburse表和reimburseFiles表中
	reimburseTableData := entity.Reimburse{
		Applicant: user.Id,
		Amount:    reimburse.Amount,
		Auditor:   0,
		Title:     reimburse.Title,
		Desc:      reimburse.Desc,
		Status:    0,
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
	err = db.AddReimburseFile(reimburse.Files, reimburseId)
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
		Msg:    "添加成功",
		Prompt: entity.SUCCESS,
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
	if reimburse.Applicant != user.Id && !authed {
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
		Status: reimburse.Status,
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
