package models

type SeatType string

const (
	SeatTypeDesk        SeatType = "desk"
	SeatTypeMeetingRoom SeatType = "meeting_room"
	SeatTypeOffice      SeatType = "office"
)

type Seat struct {
	Base

	// Название места / переговорки
	Name string `gorm:"size:255;not null" json:"name"`

	// Тип места: рабочее место, переговорная, "кабинет" отражается через office/is_private
	Type SeatType `gorm:"type:varchar(50);not null;default:'desk'" json:"type"`

	// Вместимость — полезно для переговорок и офисов
	Capacity int `gorm:"default:1" json:"capacity"`

	// Признак приватного кабинета (решено не выносить в отдельную сущность)
	IsPrivate bool `gorm:"default:false" json:"is_private"`

	// Удобства — хранится как запятую-разделённый список (можно заменить на JSON при необходимости)
	Amenities string `gorm:"type:text" json:"amenities"`

	// Цена в копейках/центах за час (или другая единица в проекте)
	PriceCents int64 `gorm:"default:0" json:"price_cents"`

	// Можно ли бронировать (например, некоторые места могут быть только по записи/админом)
	IsReservable bool `gorm:"default:true" json:"is_reservable"`

	// Связи — бронь и отзывы
	Bookings []Booking `gorm:"foreignKey:SeatID" json:"-"`
	Reviews  []Review  `gorm:"foreignKey:SeatID" json:"-"`
}
