package models

import (
	"time"
)

type Organization struct {
  ID        	 uint           `gorm:"primaryKey"`
  Name      	 string         `gorm:"uniqueIndex;not null"`
  CreatedAt 	 time.Time
  UpdatedAt 	 time.Time
  Teams     	 []Team
  Administrators []OrgAdmin 	`gorm:"constraint:OnDelete:CASCADE"`
}

type Team struct {
  ID             uint           `gorm:"primaryKey"`
  Name           string         `gorm:"not null"`
  OrganizationID uint           `gorm:"index;not null"`
  Organization   Organization   `gorm:"constraint:OnDelete:CASCADE"`
  CreatedAt      time.Time
  UpdatedAt      time.Time
  Members        []Membership
}

type User struct {
  ID        uint           `gorm:"primaryKey"`
  Email     string         `gorm:"uniqueIndex;not null"`
  Password  string         `gorm:"not null"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt time.Time
  // optionally: Memberships []Membership
}

type OrgAdmin struct {
	ID             uint           `gorm:"primaryKey"`
	UserID         uint           `gorm:"index;not null"`
	OrganizationID uint           `gorm:"index;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}	

type Membership struct {
  ID        uint           `gorm:"primaryKey"`
  UserID    uint           `gorm:"index;not null"`
  TeamID    uint           `gorm:"index;not null"`
  Role      string         `gorm:"not null"` // e.g. "admin", "member"
  CreatedAt time.Time
  UpdatedAt time.Time
}

