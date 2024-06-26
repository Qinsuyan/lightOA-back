package entity

import "time"

// 用户表
type UserRaw struct {
	Id         int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Username   string    `xorm:"varchar(20) notnull index 'username'" json:"username"`
	Password   string    `xorm:"varchar(80) notnull 'password'" json:"password"`
	Department int       `xorm:"int notnull 'department'" json:"department"`
	Phone      string    `xorm:"varchar(20) notnull unique index 'phone'" json:"phone"`
	Role       int       `xorm:"int notnull 'role'" json:"role"`
	CreatedAt  time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt  time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
}

//部门表
type Department struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Name      string    `xorm:"varchar(20) notnull index 'name'" json:"name"`
	Desc      string    `xorm:"longtext 'description'" json:"description"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'"`
}

// 登录状态表
type Online struct {
	Phone  string    `xorm:"varchar(20) notnull unique index 'phone'" json:"phone"`
	Token  string    `xorm:"varchar(64) notnull 'token'" json:"token"`
	Expire time.Time `xorm:"datetime notnull 'expire'" json:"expire"`
}

//资源表
type ResourceRaw struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Name      string    `xorm:"varchar(20) notnull index 'name'" json:"name"`    //资源的名称 e.g.添加用户
	Alias     string    `xorm:"varchar(50) notnull unique 'alias'" json:"alias"` //资源的别名 e.g.user:add
	Type      int       `xorm:"int notnull 'type'" json:"type"`                  //资源的类型,目前没有使用，默认为1
	ParentId  int       `xorm:"int 'parentId'" json:"parentId"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
}

//系统整型变量
type SystemVariableInts struct {
	Name  string `xorm:"varchar(20) pk 'name'" json:"name"`
	Value string `xorm:"int notnull 'value'" json:"value"`
}

//系统字符串变量
type SystemVariableTexts struct {
	Name  string `xorm:"varchar(20) pk 'name'" json:"name"`
	Value string `xorm:"longtext notnull 'value'" json:"value"`
}

//角色表
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

//用户日志表
type UserLog struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Content   string    `xorm:"longtext 'description'" json:"description"`
	Type      string    `xorm:"varchar(20) notnull index 'type'" json:"type"`
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
}

//文件信息表（除了本来就是用来展示的文件，还包括报销的发票、缺陷管理的附件和图片）
type File struct {
	Id        int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	UUID      string    `xorm:"varchar(36) unique index 'uuid'" json:"uuid"`   //磁盘上存储的文件名是UUID
	Name      string    `xorm:"varchar(50) notnull unique 'name'" json:"name"` //真实的文件名
	Owner     int       `xorm:"bigint(11) index 'owner'" json:"owner"`
	Type      int       `xorm:"int notnull 'type'" json:"type"` //文件类型，1:展示文件，2:报销发票，3:缺陷管理
	CreatedAt time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
	UpdatedAt time.Time `xorm:"datetime notnull updated 'updatedAt'" json:"updatedAt"`
}

//报销表
type Reimburse struct {
	Id         int       `xorm:"bigint(11) pk autoincr 'id'" json:"id"`
	Applicant  int       `xorm:"bigint(11) notnull 'applicant'" json:"applicant"`
	Amount     int       `xorm:"int notnull 'amount'" json:"amount"`
	Status     int       `xorm:"int notnull 'status'" json:"status"` //报销状态，1:待审核，3:审核通过，2:审核不通过
	Auditor    int       `xorm:"bigint(11) notnull 'auditor'" json:"auditor"`
	Title      string    `xorm:"varchar(50) notnull unique 'alias'" json:"alias"`
	Desc       string    `xorm:"longtext 'description'" json:"description"`
	CreatedBy  int       `xorm:"bigint(11) notnull 'createdBy'" json:"createdBy"`
	HappenedAt time.Time `xorm:"datetime notnull 'happenedAt'" json:"happenedAt"`
	CreatedAt  time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	DeletedAt  time.Time `xorm:"datetime notnull deleted 'deletedAt'" json:"deletedAt,omitempty"`
	UpdatedAt  time.Time `xorm:"datetime notnull updated 'updatedAt'" json:"updatedAt"`
}

type ReimbuiseFiles struct {
	ReimburseId int `xorm:"bigint(11) notnull index 'reimburseId'" json:"reimburseId"`
	FileId      int `xorm:"bigint(11) notnull index 'fileId'" json:"fileId"`
}

type ReimburseComments struct {
	Id          int       `xorm:"bigint(11) pk 'Id'" json:"Id"`
	ReimburseId int       `xorm:"bigint(11) notnull index 'reimburseId'" json:"reimburseId"`
	Comment     string    `xorm:"longtext 'comment'" json:"comment"`
	Creator     int       `xorm:"bigint(11) notnull 'creator'" json:"creator"`
	CreatedAt   time.Time `xorm:"datetime notnull created 'createdAt'" json:"createdAt"`
	UpdatedAt   time.Time `xorm:"datetime notnull updated 'updatedAt'" json:"updatedAt"`
}
