package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/zivlakmilos/perfin/db"
)

func (t *TestInfo) createUser() error {
	con := db.GetInstance()
	user := db.User{
		Id:       "",
		Username: "admin",
		Password: "password",
		Role:     "admin",
	}

	store := db.NewUserStore(con)
	err := store.Insert(&user)
	if err != nil {
		log.Error("error while creating user", "err", err)
		return err
	}

	log.Info("created user", "username", user.Username, "password", user.Password)
	return nil
}

func (t *TestInfo) testLogin() error {
	log.Info("testing user login")
	req := map[string]any{
		"username": "admin",
		"password": "password",
	}
	reqData, err := json.Marshal(req)
	if err != nil {
		log.Error("error while creating login request", "err", err)
		t.failed++
		return nil
	}
	bodyReader := bytes.NewReader(reqData)
	resp, err := http.Post(fmt.Sprintf("%s/auth/login", t.baseUrl), "application/json", bodyReader)
	if err != nil {
		log.Error("error while makeing login request", "err", err)
		t.failed++
		return nil
	}

	var res struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error("error while parsing login response", "err", err)
		t.failed++
		return nil
	}

	if res.Token == "" {
		log.Error("token missing")
		t.failed++
		return nil
	}

	t.token = res.Token
	t.passed++
	log.Info("test user login pass", "token", t.token)

	return nil
}
