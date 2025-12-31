package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       string `json:"id" xml:"id" form:"id" query:"id"`
	Username string `json:"username" xml:"username" form:"username" query:"username"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
	Role     string `json:"role" xml:"role" form:"role" query:"role"`
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

	err := s.con.Get(user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("wrong username or password")
	}

	return user, nil
}

func (s *UserStore) ChangePassword(id, newPassword string) error {
	_, err := s.con.Exec("UPDATE users SET password=? WHERE id=?", newPassword, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Insert(m *User) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	_, err := s.con.NamedExec(`INSERT INTO users (
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
