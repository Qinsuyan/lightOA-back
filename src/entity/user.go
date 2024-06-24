package entity

type UserInfo struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Phone      string `json:"phone,omitempty"`
	Department int    `json:"department"`
	Role       Role   `json:"role,omitempty"`
}

type UserPayload struct {
	Id              int    `json:"id,omitempty"`
	Username        string `json:"username"`
	Department      int    `json:"department"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
	Role            int    `json:"role"`
}

type UserToken struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}
