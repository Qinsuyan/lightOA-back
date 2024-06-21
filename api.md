# <center>API</center>

## 1. 用户登录

- POST /api/login
- 备注：除了登录请求外，其他请求需携带 token
  > token 放在`Header`中的`LTOAToken`字段
- 请求 body：json

```typescript
{
  username: string;
  password: string;
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data {
        token string //token字符串
        user User //用户信息
    }
}
```

## 2. 用户登出

- GET /api/login
- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 3. 新增用户

- POST /api/user
- 请求 body：json

```typescript
{
  username: string;
  phone: string;
  password: string;
  passwordConfirm: string;
  role: int; //角色ID
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 4. 修改自己的用户信息

- PUT /api/user
- 请求 body：json

```typescript
{
  username: string; //必填
  phone: string; //必填
  password: string; //可选，修改密码时传入
  passwordConfirm: string; //可选，修改密码时传入
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 5. 修改他人的用户信息

- PUT /api/user/:id
- 请求 body：json

```typescript
{
  username: string; //必填
  phone: string; //必填
  password: string; //可选，修改密码时传入
  passwordConfirm: string; //可选，修改密码时传入
  role: int; //可选，修改角色时传入
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 6. 删除用户

- DELETE /api/user/:id

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 7. 列出用户

- GET /api/user/list
- 请求 query

```typescript
{
  size: int; //每页数量
  index: int; //页码
  order: string; //排序字段,有多项时用","隔开
  sort: string; //排序方式,asc/desc
  username: string; //用户名
  phone: string; //手机号
  role: int; //角色
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data:{
        totol:int //总条数
        list:User[]
    }
}
```

## 8. 添加角色

- POST /api/role
- 请求 body:json

```typescript
{
   name:string //角色名称
   desc:string //角色描述
   resources:Resource[] //角色权限
}
```

Resource 类型：

```typescript
    {
        id:number,
        alias:string,
        name:string,
        type:number,
        parentId:number,
        children:Resource[]
    }
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data:{
        totol:int //总条数
        list:User[]
    }
}
```

## 9. 编辑角色

- PUT /api/role

- 请求 body:json

```typescript
{
   id:number //角色ID
   name:string //角色名称
   desc:string //角色描述
   resources:Resource[] //角色权限（可以简化结构，只包含ID）
}
```

Resource 类型：

```typescript
    {
        id:number,
        alias:string,
        name:string,
        type:number,
        parentId:number,
        children:Resource[]
    }
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data:{
        totol:int //总条数
        list:User[]
    }
}
```

## 10. 删除角色

> 没有用户使用时才可以删除

- DELETE /api/role/:id

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```

## 11. 列出角色

- GET /api/role/list
- 请求 query

```typescript
{
    size: int; //每页数量
    index: int; //页码
    sort: string; //排序方式,asc/desc（只能按name排序，所以不用指定字段）
    name: string; //用户名
}
```

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data:{
        totol:int //总条数
        list:{
                id:number //角色ID
                name:string //角色名称
                desc:string //角色描述
                resources:Resource[] //角色权限
            }[]
    }
}
```

## 12. 列出所有资源

> 虽然path里面包含`role`，但是与具体角色无关

- GET /api/role/resources

- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
    data:Resource
}
```
