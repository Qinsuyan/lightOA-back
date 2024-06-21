package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"strconv"

	"github.com/labstack/echo/v4"
)

func handleDepartmentAdd(c echo.Context) error {
	authed, _, err := checkAuth(c, "department:add")
	if err != nil || !authed {
		return err
	}
	var dep entity.Department
	err = c.Bind(&dep)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if dep.Name == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "部门名称不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if !dep.CreatedAt.IsZero() || !dep.DeletedAt.IsZero() {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return nil
	}
	ok, err := db.AddDepartment(&dep)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "添加失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if !ok {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "部门已存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "添加成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleDepartmentEdit(c echo.Context) error {
	authed, _, err := checkAuth(c, "department:edit")
	if err != nil || !authed {
		return err
	}
	var dep entity.Department
	err = c.Bind(&dep)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if dep.Id == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "部门id不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	if !dep.CreatedAt.IsZero() || !dep.DeletedAt.IsZero() {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return nil
	}
	err = db.EditDepartment(&dep)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "编辑失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "编辑成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleDepartmentDelete(c echo.Context) error {
	authed, _, err := checkAuth(c, "department:delete")
	if err != nil || !authed {
		return err
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "部门id不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	depId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "无法解析部门id",
			Prompt: entity.WARN,
		})
		return err
	}
	//检查部门下是否有员工
	exist, err := db.ExistUserOfDepartment(int(depId))
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if exist {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "部门下存在员工",
			Prompt: entity.ERROR,
		})
		return nil
	}
	err = db.DeleteDepartment(int(depId))
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "删除成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleDepartmentList(c echo.Context) error {
	authed, _, err := checkAuth(c, "department:list")
	if err != nil || !authed {
		return err
	}
	var filter entity.ListRequest
	err = c.Bind(&filter)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	total, list, err := db.ListDepartment(&filter)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[entity.ListResponse[entity.Department]]{
		Code:   OK,
		Msg:    "获取成功",
		Prompt: entity.SUCCESS,
		Data: entity.ListResponse[entity.Department]{
			List:  list,
			Total: total,
		},
	})
	return nil
}
