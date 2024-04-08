package user

import "time"

type RoleName string

const RoleUser RoleName = "ROLE_USER"
const RoleAdmin RoleName = "ROLE_ADMIN"

type Role struct {
	id        uint 
	createdAt time.Time
	updatedAt time.Time
	name      RoleName 
}

func NewRole(name RoleName) *Role {
	return &Role{
		name:      name,
	}
}

func RoleFromDB(id uint, name RoleName) *Role {
	return &Role{
		id:        id,
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (r *Role) GetId() uint {
	return r.id
}

func (r *Role) GetCreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) GetUpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Role) GetName() RoleName {
	return r.name
}

