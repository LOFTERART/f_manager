package model

type Admin struct {
	BaseInfo
	Name     string `json:"name" gorm:"size:20" json:"name,omitempty"`
	Password string `json:"password"  gorm:"size:20" json:"password,omitempty"`
	Tokens   string `json:"tokens"  gorm:"size:20" json:"tokens,omitempty"`
	Roles    string `json:"roles,omitempty"`
	Avatar   string `json:"avatar"`
}
