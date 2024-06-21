package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"

	"github.com/labstack/echo/v4"
)

func handleResourceList(c echo.Context) error {
	authed, _, err := checkAuth(c, "resources:list")
	if err != nil || !authed {
		return err
	}
	list, err := db.GetAllResources()
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "获取资源配置失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	return c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "获取成功",
		Data:   list,
		Prompt: entity.SILENT,
	})
}
