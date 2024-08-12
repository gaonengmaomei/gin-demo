package models

type User struct {
	ID   int
	Name string
}

type Role struct {
	ID     int
	Name   string
	Active bool
}

func (User) TableName() string {
	return "gin_user"
}
func (Role) TableName() string {
	return "gin_role"
}
