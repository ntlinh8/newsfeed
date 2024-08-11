package model

// gorm model
type DbUser struct {
	Id           int    `gorm:"column:id"`
	Username     string `gorm:"column:username"`
	HashPassword string `gorm:"column:hash_password"`
	FirstName    string `gorm:"column:first_name"`
	LastName     string `gorm:"column:last_name"`
	DOB          int    `gorm:"column:dob"`
	Email        string `gorm:"column:email"`
}
