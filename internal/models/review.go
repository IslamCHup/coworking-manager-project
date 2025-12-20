package models

type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

type Review struct {
	Base

	UserID uint   `gorm:"not null" json:"user_id"`
	SeatID uint   `gorm:"not null" json:"seat_id"`
	Rating uint8  `gorm:"type:tinyint;default:5" json:"rating"` // 1-5
	Text   string `gorm:"type:text" json:"text"`

	// Модерация
	Status         ReviewStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ModeratedByID  *uint        `json:"moderated_by_id,omitempty"`
	ModerationNote string       `gorm:"type:text" json:"moderation_note,omitempty"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"-"`
	Seat Seat `gorm:"foreignKey:SeatID" json:"-"`
}
