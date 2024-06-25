package entity

// Prompt code
const ERROR = 0
const WARN = 1
const INFO = 3
const SUCCESS = 4
const SILENT = 5

// http请求响应
type HttpResponse[T any] struct {
	Prompt int    `json:"prompt,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Data   T      `json:"data,omitempty"`
	Code   int    `json:"code"`
}
type ListResponse[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}
type ListRequest struct {
	PageSize int `query:"size"`
	PageNum  int `query:"index"`
}
