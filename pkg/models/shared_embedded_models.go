package models

import (
	"time"

	"gorm.io/gorm"
)

type Identifier struct {
	ID uint `json:"id" gorm:"primarykey;autoIncrement"`
}

type TimeSpans struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ModelBasics struct {
	ID        uint      `json:"id" gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ModelBasicsTrackedDel struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type DeletedAtCol struct {
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
