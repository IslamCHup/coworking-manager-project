package models

type Admin struct {
	Base

	Provider       string `json:"provider" gorm:"not null"`
	ProviderUserID string `json:"provider_user_id" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`

	Reviews []Review `json:"-"`
}
