package models

type User struct {
	ID       uint    `gorm:"primaryKey"`
	Email    string  `gorm:"uniqueIndex;not null"`
	Password string  `gorm:"not null"`
	Events   []Event `gorm:"foreignKey:CreatorID"`
}

type Event struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Org       string
	Team      string
	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID"`
}
