package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"

	"github.com/labstack/echo/v4"
)

// 添加role
func handleRoleAdd(c echo.Context) error {
	authed, _, err := checkAuth(c, "role:add")
	if err != nil || !authed {
		return err
	}
	var role entity.Role
	err = c.Bind(&role)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if role.Name == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "角色名不能为空",
			Prompt: entity.WARN,
		})
		return nil
	}
	if len(role.Resources) == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "无法添加没有任何权限的角色",
			Prompt: entity.WARN,
		})
		return nil
	}
	//检查角色姓名是否存在
	has, err := db.IsRoleNameExist(role.Name)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "添加失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if has {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "角色名已存在",
			Prompt: entity.WARN,
		})
		return nil
	}
	ok, err := db.AddRole(&entity.RoleRaw{
		Name: role.Name,
		Desc: role.Desc,
	})
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
			Msg:    "角色名已存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	//添加资源
	roleId, err := db.GetRoleIdByName(role.Name)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "资源配置失败",
			Prompt: entity.WARN,
		})
		return err
	}
	if roleId == 0 {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "添加失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	ok, err = db.AddRoleResource(roleId, role.Resources)
	if !ok || err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "资源配置失败",
			Prompt: entity.WARN,
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

// 编辑role
func handleRoleEdit(c echo.Context) error {
	authed, _, err := checkAuth(c, "role:edit")
	if err != nil || !authed {
		return err
	}
	var role entity.Role
	err = c.Bind(&role)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if role.Id == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "角色id不能为空",
			Prompt: entity.WARN,
		})
	}
	if role.Name == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "角色名不能为空",
			Prompt: entity.WARN,
		})
	}
	if len(role.Resources) == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "必须指定权限",
			Prompt: entity.WARN,
		})
		return nil
	}
	err = db.EditRole(&entity.RoleRaw{
		Name: role.Name,
		Desc: role.Desc,
	})
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "编辑失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	//删除旧的资源
	err = db.DeleteRoleResource(role.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "资源配置失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	//重新添加资源
	ok, err := db.AddRoleResource(role.Id, role.Resources)
	if !ok || err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "资源配置失败",
			Prompt: entity.WARN,
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
