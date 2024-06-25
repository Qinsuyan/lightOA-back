package entity

type Role struct {
	Id        int         `json:"id,omitempty"`
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	Resources []*Resource `json:"resources"` //nil = all
}

type Resource struct {
	Id       int         `json:"id,omitempty"`
	Alias    string      `json:"alias"`
	Name     string      `json:"name"`
	Type     int         `json:"type"`
	Children []*Resource `json:"children"`
	ParentId int         `json:"parentId,omitempty"`
}
type RoleListFilter struct {
	ListRequest
	Name string `query:"name"`
	Sort string `query:"sort"` //desc,asc
}
