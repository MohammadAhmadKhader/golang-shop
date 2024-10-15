package models

type Role struct {
	Identifier
	Name      string      `json:"role" gorm:"size:32;not null"`
	Users     []User `json:"users,omitempty" gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:UserID;constraint:OnDelete:CASCADE;"`
	TimeSpans
}
