package test

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type TestInfo struct {
	passed  int
	failed  int
	baseUrl string
	token   string
}

func (t *TestInfo) setup() error {
	err := t.createUser()
	if err != nil {
		return err
	}

	return nil
}

func (t *TestInfo) cleanup() {
}

func RunTests() {
	t := TestInfo{
		passed:  0,
		failed:  0,
		baseUrl: fmt.Sprintf("http://127.0.0.1:%s", os.Getenv("PORT")),
		token:   "",
	}

	log.Info("setup tests")
	err := t.setup()
	if err != nil {
		log.Fatal("error in tests setup", "error", err)
	}

	log.Info("running tests")
	t.testLogin()

	log.Info("celanup tests")
	t.cleanup()

	if t.failed > 0 {
		log.Error("all tests finished", "failed", t.passed, "failed", t.failed)
	} else {
		log.Info("all tests finished", "passed", t.passed, "failed", t.failed)
	}
}
