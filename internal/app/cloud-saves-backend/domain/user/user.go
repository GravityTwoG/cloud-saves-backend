package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/domain_errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id uint

	username  string
	email     string
	password  string
	isBlocked bool
	role      Role

	createdAt time.Time
	updatedAt time.Time
}

func NewUser(
	username string,
	email string,
	password string,
	role *Role,
) (*User, error) {
	user := User{
		role:      *role,
		isBlocked: false,
	}

	if len(email) < 3 || len(email) > 256 {
		return nil, domain_errors.NewErrInvalidInput(
			"email length must be between 3 and 256 characters",
		)
	}
	user.email = email

	if len(username) < 3 || len(username) > 32 {
		return nil, domain_errors.NewErrInvalidInput(
			"username length must be between 3 and 32 characters",
		)
	}
	user.username = username

	err := user.ChangePassword(password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserFromDB(
	id uint,
	username,
	email,
	password string,
	role *Role,
	isBlocked bool,
) *User {
	return &User{
		id:        id,
		username:  username,
		email:     email,
		password:  password,
		role:      *role,
		isBlocked: isBlocked,
	}
}

func (u *User) GetId() uint {
	return u.id
}

func (u *User) SetId(id uint) {
	u.id = id
}

func (u *User) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *User) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetRole() *Role {
	return &u.role
}

func (u *User) IsBlocked() bool {
	return u.isBlocked
}

func (u *User) Block() {
	u.isBlocked = true
}

func (u *User) Unblock() {
	u.isBlocked = false
}

func (u *User) ChangePassword(rawPassword string) error {
	if len(rawPassword) < 8 || len(rawPassword) > 64 {
		return domain_errors.NewErrInvalidInput(
			"password length must be between 8 and 64 characters",
		)
	}

	hashedPassword, err := hashPassword(rawPassword)
	if err != nil {
		return err
	}
	u.password = hashedPassword
	return nil
}

func (u *User) ComparePassword(rawPassword string) bool {
	return comparePasswords(u.password, rawPassword)
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
