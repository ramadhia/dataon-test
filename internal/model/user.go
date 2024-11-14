package model

import "time"

type Claim struct {
	ID       string `json:"user_id" binding:"required"`
	Role     string `json:"role"`
	IsClient bool   `json:"is_client"`
	Token    string `json:"token"`
}

type User struct {
	ID          *string    `json:"id"`
	FirstName   *string    `json:"first_name"`
	LastName    *string    `json:"last_name"`
	PhoneNumber *string    `json:"phone_number"`
	Address     *string    `json:"address"`
	Pin         *string    `json:"pin"`
	CreatedDate *time.Time `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
}

type AccessToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:varchar(50);not null" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null" json:"token"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Relationship with User table
}

func (user Claim) IsClientToken() bool {
	return user.IsClient
}
