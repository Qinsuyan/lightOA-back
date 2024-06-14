package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
	"time"

	"github.com/labstack/echo/v4"
)

// 用户登录
func handleUserLogin(c echo.Context) error {
	var payload entity.UserPayload
	err := c.Bind(&payload)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "登录参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if payload.Password == "" || payload.Username == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "缺少参数",
			Prompt: entity.WARN,
		})
		return nil
	}
	userRaw :=
		&entity.UserRaw{Username: payload.Username, Password: util.Sha256(payload.Password)}
	exist, err := db.GetUserRaw(userRaw)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "登录失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if !exist || !userRaw.DeletedAt.IsZero() {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "用户名或密码错误",
			Prompt: entity.WARN,
		})
		return nil
	}
	token := util.FormToken(userRaw.Username)
	on := &entity.Online{
		Username: userRaw.Username,
		Token:    token,
		Expire:   time.Now().Add(24 * time.Hour),
	}
	trueToken, err := db.LoginUser(on)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "登录失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	userInfo := &entity.UserInfo{Id: userRaw.Id, Username: userRaw.Username}
	userRole, err := db.GetUserRoleByRoleId(userRaw.Role)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "登录失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if userRole != nil {
		userInfo.Role = *userRole
	}
	c.JSON(OK, entity.HttpResponse[entity.UserToken]{
		Code: OK,
		Data: entity.UserToken{
			User:  *userInfo,
			Token: trueToken,
		},
	})
	return nil
}

// 注销登录
func handleUserLogout(c echo.Context) error {
	online := &entity.Online{
		Token: c.Request().Header.Get("LTOAToken"),
	}
	db.LogoutUser(online)
	c.JSON(OK, entity.HttpResponse[any]{
		Code: OK,
		Msg:  "已退出登录",
	})
	return nil
}

// 添加用户
func handleUserAdd(c echo.Context) error {
	authed, _, err := checkAuth(c, "user:add")
	if err != nil || !authed {
		return err
	}
	var user entity.UserPayload
	err = c.Bind(&user)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if user.Username == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户名不能为空",
			Prompt: entity.WARN,
		})
		return nil
	}
	if user.Password == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "密码不能为空",
			Prompt: entity.WARN,
		})
		return nil
	}
	if user.Role == 0 {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "必须指定角色",
			Prompt: entity.WARN,
		})
		return nil
	}
	ok, err := db.AddUser(&entity.UserRaw{
		Username: user.Username,
		Password: util.Sha256(user.Password),
		Role:     user.Role,
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
			Msg:    "用户已存在",
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
