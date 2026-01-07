package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/zivlakmilos/perfin/db"
	"github.com/zivlakmilos/perfin/test"
)

func main() {
	log.SetReportCaller(true)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := os.Getenv("DB_URL")
	log.Infof("remove old db: %s", dbPath)
	err = os.Remove(dbPath)
	if err != nil && os.IsNotExist(err) {
		log.Fatal(err)
		log.Fatal("error while removing db", "err", err)
		return
	}

	executable, _ := os.Executable()

	migratePath := path.Join(filepath.Dir(executable), "migrate")

	log.Infof("running migrations: %s test", migratePath)
	cmd := exec.Command(migratePath, "test")
	err = cmd.Run()
	if err != nil {
		log.Fatal("error while running migrations", "err", err)
		return
	}

	apiPath := path.Join(filepath.Dir(executable), "api")
	log.Infof("starting api: %s", apiPath)

	cmd = exec.Command(apiPath)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Fatal("starting api failed", "err", err)
		}
	}()

	time.Sleep(1000 * time.Millisecond)

	err = db.CreateConnection(dbPath)
	if err != nil {
		log.Fatal(err)
		log.Fatal("error while connecting to db", "err", err)
	}
	log.Info("new database created")

	test.RunTests()

	_ = cmd.Process.Kill()
}
