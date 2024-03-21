package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsBlocked bool   `gorm:"not null"`
	RoleID    uint   `gorm:"not null"`
	Role      Role   `gorm:"onDelete:RESTRICT"`
}

func NewUser(
	username string,
	email string,
	password string,
	role *Role,
) (*User, error) {
	user := User{
		Username:  username,
		Email:     email,
		Role:      *role,
		IsBlocked: false,
	}

	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) SetPassword(rawPassword string) error {
	hashedPassword, err := hashPassword(rawPassword)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) ComparePassword(rawPassword string) bool {
	return comparePasswords(u.Password, rawPassword)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
