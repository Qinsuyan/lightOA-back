package entity

type Role struct {
	Id        int        `json:"id,omitempty"`
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"` //nil = all
}

type Resource struct {
	Id       int        `json:"id,omitempty"`
	Alias    string     `json:"alias"`
	Name     string     `json:"name"`
	Type     int        `json:"type"` //1 menu 2 button 3 root
	Children []Resource `json:"children"`
	ParentId int        `json:"parentId"`
}
