package api

import (
	"lightOA-end/src/db"
	"lightOA-end/src/entity"
	"lightOA-end/src/util"
	"strconv"
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
	if payload.Password == "" || payload.Phone == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "缺少参数",
			Prompt: entity.WARN,
		})
		return nil
	}
	userRaw :=
		&entity.UserRaw{Phone: payload.Username, Password: util.Sha256(payload.Password)}
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
		Phone:  userRaw.Phone,
		Token:  token,
		Expire: time.Now().Add(24 * time.Hour),
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
	db.LogoutUserByToken(c.Request().Header.Get("LTOAToken"))
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
	if user.Phone == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "手机号不能为空",
			Prompt: entity.WARN,
		})
		return nil
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
	if user.Password != user.PasswordConfirm {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "两次密码输入不匹配",
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
	payload := &entity.UserRaw{
		Username: user.Username,
		Password: util.Sha256(user.Password),
		Role:     user.Role,
		Phone:    user.Phone,
	}

	ok, err := db.AddUser(payload)
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

// 修改自身用户信息
func handleSelfModify(c echo.Context) error {
	authed, user, err := checkAuth(c, "user:self")
	if err != nil || !authed {
		return err
	}
	//只能更改 username, phone, password且必须都带上
	var userNew entity.UserPayload
	err = c.Bind(&userNew)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if userNew.Username == "" || userNew.Phone == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "缺少参数",
			Prompt: entity.WARN,
		})
		return nil
	}
	if userNew.Password != userNew.PasswordConfirm {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "两次密码输入不匹配",
			Prompt: entity.WARN,
		})
		return nil
	}
	payload := &entity.UserRaw{
		Username: userNew.Username,
		//Password: util.Sha256(userNew.Password),
		Phone: userNew.Phone,
	}
	if userNew.Password != "" {
		payload.Password = util.Sha256(userNew.Password)
	}
	err = db.EditUser(user.Id, payload)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "修改用户信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "修改用户信息成功",
		Prompt: entity.SUCCESS,
	})
	return nil
}

func handleUserModify(c echo.Context) error {
	authed, _, err := checkAuth(c, "user:edit")
	if err != nil || !authed {
		return err
	}
	userId := c.Param("id")

	if userId == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户ID不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}

	userIdNum, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户ID错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	user := &entity.UserRaw{Id: int(userIdNum)}
	exist, err := db.GetUserRaw(user)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "修改失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if !exist {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	userNew := entity.UserPayload{}
	err = c.Bind(&userNew)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return err
	}

	if userNew.Username == "" || userNew.Phone == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "缺少参数",
			Prompt: entity.WARN,
		})
		return nil
	}
	if userNew.Password != userNew.PasswordConfirm {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "两次密码输入不一致",
			Prompt: entity.WARN,
		})
	}
	payload := entity.UserRaw{
		Username: userNew.Username,
		Phone:    userNew.Phone,
	}
	if userNew.Password != "" {
		payload.Password = util.Sha256(userNew.Password)
	}
	if userNew.Role != 0 {
		role, err := db.GetRoleRawById(userNew.Role)
		if err != nil {
			c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
				Code:   ERROR_INTERNAL,
				Msg:    "指定的角色不存在",
				Prompt: entity.WARN,
			})
			return err
		}
		payload.Role = role.Id
	}
	err = db.EditUser(int(userIdNum), &payload)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "修改失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	return c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "修改成功",
		Prompt: entity.SUCCESS,
	})
}

func handleUserDelete(c echo.Context) error {
	authed, _, err := checkAuth(c, "user:del")
	if err != nil || !authed {
		return err
	}
	userId := c.Param("id")
	if userId == "" {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户ID不能为空",
			Prompt: entity.ERROR,
		})
		return nil
	}
	userIdNum, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户ID错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	user := entity.UserRaw{Id: int(userIdNum)}
	exist, err := db.GetUserRaw(&user)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	if !exist {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "用户不存在",
			Prompt: entity.ERROR,
		})
		return nil
	}
	//下线用户
	db.LogoutUserByPhone(user.Phone)
	//删除用户
	err = db.DeleteUser(user.Id)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "删除失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	return c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "删除成功",
		Prompt: entity.SUCCESS,
	})
}

func handleUserList(c echo.Context) error {
	authed, _, err := checkAuth(c, "user:list")
	if err != nil || !authed {
		return err
	}
	var filter entity.UserListFilter
	if err := c.Bind(&filter); err != nil {
		c.JSON(ERROR_INVALID_PARAM, entity.HttpResponse[any]{
			Code:   ERROR_INVALID_PARAM,
			Msg:    "参数错误",
			Prompt: entity.ERROR,
		})
		return err
	}
	// if filter.PageSize == 0 {
	// 	filter.PageSize = 10
	// }
	// if filter.PageNum == 0 {
	// 	filter.PageNum = 1
	// }
	users, err := db.ListUser(&filter)
	if err != nil {
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "查询用户信息失败",
			Prompt: entity.ERROR,
		})
		return err
	}
	return c.JSON(OK, entity.HttpResponse[any]{
		Code:   OK,
		Msg:    "查询用户信息成功",
		Prompt: entity.SUCCESS,
		Data: entity.ListResponse[entity.UserInfo]{
			Total: int64(len(users)),
			List:  users,
		},
	})
}
