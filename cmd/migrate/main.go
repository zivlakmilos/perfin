package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/zivlakmilos/perfin/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error while loading .env file: %s\n", err)
	}

	err = db.CreateConnection(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("error while connecting to db: %s\n", err)
	}

	con := db.GetInstance()
	defer con.Close()

	con.MustExec(`CREATE TABLE IF NOT EXISTS User (
  	id TEXT PRIMARY KEY,
  	username TEXT,
  	password TEXT,
  	role TEXT
	);`)

	createAdmin(con)

	fmt.Printf("migration success\n")
}

func createAdmin(con *sqlx.DB) {
	rows, err := con.Query("SELECT COUNT(*) FROM User")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	if count > 0 {
		return
	}

	var username string
	var password string

	fmt.Printf("admin username: ")
	fmt.Scanf("%s", &username)
	fmt.Printf("admin password: ")
	fmt.Scanf("%s", &password)

	user := db.User{
		Id:       "",
		Username: username,
		Password: password,
		Role:     "admin",
	}

	store := db.NewUserStore(con)
	err = store.Insert(&user)
	if err != nil {
		log.Fatal(err)
	}
}
