package entity

type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type UserPayload struct {
	Id              int    `json:"id,omitempty"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Role            int    `json:"role"`
}

type UserToken struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}
