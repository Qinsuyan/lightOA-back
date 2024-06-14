# <center>API</center>

1. 用户登录

- POST /api/login
- 备注：除了登录请求外，其他请求需携带 token
    > token放在`Header`中的`LTOAToken`字段
- 请求 body：json

```typescript
{
  username:string
  password:string
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

2. 用户登出

- GET /api/login
- 响应 body：json

```go
{
    code int       //200-成功 非200-失败
    msg string //提示信息
}
```