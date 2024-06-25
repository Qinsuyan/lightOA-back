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
		} else {
			log.Err(err).Msgf("err while %s %s", c.Request().Method, c.Request().URL.Path)
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
		//部门管理
		department := onlines.Group("/department")
		{
			department.POST("", handleDepartmentAdd)          //13.新增部门
			department.PUT("", handleDepartmentEdit)          //14.编辑部门
			department.DELETE("/:id", handleDepartmentDelete) //15.删除部门
			department.GET("/list", handleDepartmentList)     //16.列出部门
		}
		//报销管理
		reimburse := onlines.Group("/reimburse")
		{
			//管理
			reimburse.POST("", handleReimburseAdd)            //17.新增报销
			reimburse.POST("/files", handleReimburseAddFile)  //18.新增报销文件
			reimburse.PUT("", handleReimburseEdit)            //19.编辑报销信息（只能申请者编辑未审核、审核失败的报销）
			reimburse.DELETE("/:id", handleReimburseDelete)   //20.删除报销（只能申请者删除未审核的报销）
			reimburse.GET("/history", handleReimburseHistory) //21.列出自己的报销
			// reimburse.GET("/list", handleReimburseList)           //22.列出所有人的报销
			// reimburse.POST("/audit", handleReimburseAudit)        //23.审核报销信息
			// reimburse.GET("/auditors", handleReimburseAudit)        //24.获取有审核权限的人员（包括已删除的人员）
			// reimburse.POST("/reapply", handleReimburseReapply)        //25.重新提交报销信息
			reimburse.GET("/:id", handleReimburseDetail) //26.查看报销详情
			// reimburse.GET("/sheet", handleReimburseExportAsSheet) //27.导出报销表格
			// reimburse.GET("/zip", handleReimburseExportAsZip)     //28.导出报销的所有文件（日期-人员-事由-文件；表格）
			// 统计
			// summary:
			// reimburse.GET("/statistic", handleReimburseStatistic) //29.报销数据统计
			// total_requests: 总报销请求数。
			// reimbursed_requests: 已报销请求数。
			// unreimbursed_requests: 未报销请求数。
			// total_amount: 报销总金额。
			// monthly_totals:
			// 每个月的总报销金额和报销请求数量。
			// daily_totals:
			// 每天的总报销金额和报销请求数量。
			// department_totals:
			// 每个部门的总报销金额和报销请求数量。
			// category_totals:
			// 每个报销类别的总报销金额和报销请求数量。
			// individual_requests:
			// 每个报销请求的详细信息，包括ID、金额、日期、状态、部门和类别。
		}
		//文件管理
		files := onlines.Group("/files")
		{
			files.GET("/bin/:id", handleGetFileData) //30.获取文件数据（不返回fileinfo，直接返回二进制数据）
			files.GET("/:id", handleGetFileInfo)     //31.获取文件数据（返回fileinfo）
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

func getAuth(c echo.Context, auth string) (bool, *entity.UserRaw, error) {
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
	return authorized, user, nil
}

func getUser(c echo.Context) (*entity.UserRaw, error) {
	token := c.Request().Header.Get("LTOAToken")
	user, err := db.GetUserRawByToken(token)
	if err != nil {
		log.Err(err).Msg("err while getting user")
		c.JSON(ERROR_INTERNAL, entity.HttpResponse[any]{
			Code:   ERROR_INTERNAL,
			Msg:    "参数解析失败",
			Prompt: entity.ERROR,
		})
		return nil, err
	}
	return user, nil
}
