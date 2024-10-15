package models

import "time"

type UserRoles struct {
	UserID     uint      `json:"userId" gorm:"primaryKey;not null"`
	User       *User     `json:"user,omitempty"`
	RoleID     uint      `json:"roleId" gorm:"primaryKey;not null"`
	Role       *Role     `json:"role,omitempty"`
	AssignedAt time.Time `json:"assignedAt,omitempty" gorm:"autoCreateTime"`
}
