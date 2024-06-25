package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/storage"
	"strconv"

	"github.com/labstack/echo/v4"
)

func handleGetFileData(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	file, err := db.GetFileById(id)
	if err != nil || file == nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "获取失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	// 按文件类型检查权限
	authed := false
	if file.Type == 1 {
		//展示文件
		authed = true
	}
	if file.Type == 2 {
		//发票
		//申请者或审核人员有权限
		authed, user, err := getAuth(c, "reimburse:audit")
		if err == nil {
			if authed || user.Id == file.Owner {
				authed = true
			}
		}
	}
	if !authed {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.WARN,
		})
		return err
	}
	data, err := storage.GetFile(file)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+file.Name)
	c.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	c.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(len(data)))

	return c.Blob(OK, "application/octet-stream", data)
}

func handleGetFileInfo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	file, err := db.GetFileById(id)
	if err != nil || file == nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "获取失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	// 按文件类型检查权限
	authed := false
	if file.Type == 1 {
		//展示文件
		authed = true
	}
	if file.Type == 2 {
		//发票
		//申请者或审核人员有权限
		authed, user, err := getAuth(c, "reimburse:audit")
		if err == nil {
			if authed || user.Id == file.Owner {
				authed = true
			}
		}
	}
	if !authed {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.WARN,
		})
		return err
	}
	data, err := storage.GetFile(file)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[entity.FileInfo]{
		Code:   OK,
		Msg:    "获取成功",
		Prompt: entity.SILENT,
		Data: entity.FileInfo{
			Id:   file.Id,
			Name: file.Name,
			File: data,
		},
	})
	return nil
}
