package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	Identifier
	FullName      string          `json:"fullName" gorm:"not null;size:4;size:32"`
	City          string          `json:"city" gorm:"not null;size:4;size:32"`
	StreetAddress string          `json:"streetAddress" gorm:"not null;size:4;size:64"`
	State         *string         `json:"state" gorm:"size:4;size:32"`
	ZipCode       *string         `json:"zipCode" gorm:"size:3;size:12"`
	Country       string          `json:"country" gorm:"not null;size:4;size:32"`
	UserID        uint            `json:"userId" gorm:"not null"`
	User          *User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time       `json:"createdAt"`
	DeletedAt     *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

func (a *Address) GetUserId() uint {
	return a.UserID
}
