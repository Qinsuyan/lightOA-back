# <center>lightOA 后端项目</center>

## 已设计的表模式

###  User 
用户
###  Role 
角色
### Resource 
资源
### RoleResource 
角色资源关系
按`RoleId`查询到的所有`Resource`，即为这个角色所拥有的权限
### Online 
登录状态

小范围使用，直接用数据库保存登录状态
### Log 
用户操作日志

