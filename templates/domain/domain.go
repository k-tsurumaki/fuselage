package domain

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}