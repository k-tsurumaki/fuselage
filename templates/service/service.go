package service

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id int) error
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(user *User) error {
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}
