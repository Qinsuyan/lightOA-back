package entity

import "time"

// 用户表
type UserRaw struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Username  string    `xorm:"varchar(20) notnull unique index 'username'" json:"username"`
	Password  string    `xorm:"varchar(50) notnull 'password'" json:"password"`
	Role      int       `xorm:"int notnull 'role'" json:"role"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
}

// 登录状态表
type Online struct {
	Username string    `xorm:"varchar(20) notnull unique index 'username'" json:"username"`
	Token    string    `xorm:"varchar(64) notnull 'token'" json:"token"`
	Expire   time.Time `xorm:"datetime notnull 'expire'" json:"expire"`
}

type ResourceRaw struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Name      string    `xorm:"varchar(20) notnull index 'name'" json:"name"`
	Alias     string    `xorm:"varchar(20) notnull unique 'alias'" json:"alias"`
	Type      int       `xorm:"int notnull 'type'" json:"type"`
	ParentId  int       `xorm:"int 'parentId'" json:"parentId"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
}

type RoleRaw struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Name      string    `xorm:"varchar(20) notnull index 'name'" json:"name"`
	Desc      string    `xorm:"longtext 'description'" json:"description"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
}

// RoleResource 结构体定义了角色与资源之间的关系
// 该结构体包含两个字段，RoleId 和 ResourceId，分别表示角色的ID和资源的ID。
// 这两个字段都是主键，用于唯一标识一个角色资源关系。
type RoleResource struct {
	RoleId     int `xorm:"bigint(11) pk 'roleId'" json:"roleId"`         // 角色ID，bigint类型，作为主键
	ResourceId int `xorm:"bigint(11) pk 'resourceId'" json:"resourceId"` // 资源ID，bigint类型，作为主键
}

type UserLog struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Content   string    `xorm:"longtext 'description'" json:"description"`
	Type      string    `xorm:"varchar(20) notnull index 'type'" json:"type"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
}
