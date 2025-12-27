package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       string
	Username string
	Password string
	Role     string
}

type UserStore struct {
	con *sqlx.DB
}

func NewUser() *User {
	return &User{}
}

func NewUserStore(con *sqlx.DB) *UserStore {
	return &UserStore{
		con: con,
	}
}

func (s *UserStore) Login(username, password string) (*User, error) {
	user := NewUser()

	err := s.con.Get(user, "SELECT * FROM User WHERE username=?", username)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("wrong username or password")
	}

	return user, nil
}

func (s *UserStore) ChangePassword(id, newPassword string) error {
	_, err := s.con.Exec("UPDATE User SET password=? WHERE id=?", newPassword, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Insert(m *User) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	_, err := s.con.NamedExec(`INSERT INTO User (
		id,
		username,
		password,
		role
	) VALUES (
		:id,
		:username,
		:password,
		:role
	)`, m)
	if err != nil {
		return err
	}

	return nil
}
