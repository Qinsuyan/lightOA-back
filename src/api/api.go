package api

import (
	"fmt"
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/log"

	"github.com/labstack/echo/v4"
)

// 错误码
const OK = 200
const ERROR_INTERNAL = 500
const ERROR_AUTH = 401
const ERROR_INVALID_PARAM = 400

func online(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("LTOAToken")
		//只检查登录状态，权限和用户是否禁用在接口中检查
		logged, err := db.IsUserOnline(token)
		if err != nil {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    fmt.Sprintf("%v", err),
				Prompt: entity.ERROR,
			})
			return nil
		}
		if logged {
			return next(c)
		} else {
			c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
				Code:   ERROR_AUTH,
				Msg:    "登录状态已失效",
				Prompt: entity.WARN,
			})
			return nil
		}
	}
}

func Start(port string, dist string) error {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := ERROR_INTERNAL
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == 404 {
			err := c.File(dist + "/index.html")
			if err != nil {
				log.Err(err).Msgf("err while %s %s", c.Request().Method, c.Request().URL.Path)
			}
		}
	}
	e.Static("/", dist)
	api := e.Group("/api")
	// 用户登录，不鉴权
	api.POST("/token", handleUserLogin) //1.用户登录
	api.GET("/token", handleUserLogout) //2.用户登出
	{
		onlines := api.Group("")
		onlines.Use(online)
		//用户操作
		user := onlines.Group("/user")
		{
			user.POST("", handleUserAdd)          //3.新增用户
			user.PUT("", handleSelfModify)        //4.修改自身用户信息
			user.PUT("/:id", handleUserModify)    //5.修改其他用户的信息
			user.DELETE("/:id", handleUserDelete) //6.删除用户
			user.GET("/list", handleUserList)     //7.列出用户
		}
		//角色操作
		role := onlines.Group("/role")
		{
			role.POST("", handleRoleAdd)               //8.新增角色
			role.PUT("", handleRoleEdit)               //9.编辑角色
			role.DELETE("/:roleId", handleRoleDelete)  //10.删除角色
			role.GET("/list", handleRoleList)          //11.列出角色
			role.GET("/resources", handleResourceList) //12.列出所有的资源
		}
	}
	err := e.Start(port)
	return err
}

func checkAuth(c echo.Context, auth string) (bool, *entity.UserRaw, error) {
	token := c.Request().Header.Get("LTOAToken")
	authorized, user, err := db.IsUserAuthorized(auth, token)
	if err != nil {
		log.Err(err).Msg("err while getting user")
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return false, nil, err
	}
	if !authorized {
		c.JSON(ERROR_AUTH, entity.HttpResponse[any]{
			Code:   ERROR_AUTH,
			Msg:    "无权限",
			Prompt: entity.WARN,
		})
		return false, user, nil
	}
	return true, user, nil
}
