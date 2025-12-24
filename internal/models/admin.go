package models

type Admin struct {
	Base

	Login        string `json:"login" gorm:"not null;uniqueIndex"`
	PasswordHash string `json:"-" gorm:"not null"`

	Reviews []Review `json:"-"`
}
