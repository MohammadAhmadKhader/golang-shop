package models

type Message struct {
	Identifier
	From     uint   `json:"from" gorm:"not null;index"`
	FromUser *User  `json:"sender,omitempty" gorm:"foreignKey:From;OnDelete:CASCADE"`
	To       uint   `json:"to" gorm:"not null;index"`
	Content  string `json:"content" gorm:"not null;size:256"`
	ToUser   *User  `json:"recipient,omitempty" gorm:"foreignKey:To;OnDelete:CASCADE"`
	Status   string `json:"status"  gorm:"not null"`
	TimeSpans
}

func (m *Message) GetUserId() uint {
	return m.From
}