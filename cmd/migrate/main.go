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
	defer func() { _ = con.Close() }()

	con.MustExec(`CREATE TABLE IF NOT EXISTS users (
  	id TEXT PRIMARY KEY,
  	username TEXT,
  	password TEXT,
  	role TEXT
	);`)

	con.MustExec(`CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		account_type TEXT,
		parent_id TEXT,
		title TEXT
	);`)

	con.MustExec(`CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		transaction_id TEXT,
		account_id TEXT,
		date TEXT,
		description TEXT,
		debit REAL,
		credit REAL
	);`)

	if len(os.Args) == 1 || os.Args[1] != "test" {
		createAdmin(con)
		createAccounts(con)
	}

	fmt.Printf("migration success\n")
}

func createAdmin(con *sqlx.DB) {
	rows, err := con.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

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
	_, _ = fmt.Scanf("%s", &username)
	fmt.Printf("admin password: ")
	_, _ = fmt.Scanf("%s", &password)

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

func createAccounts(con *sqlx.DB) {
	rows, err := con.Query("SELECT COUNT(*) FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = rows.Close() }()

	var count int

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Fatal(err)
		}
	}

	if count > 0 {
		return
	}

	accounts := []db.Account{
		// Equity
		{
			Id:          "e4c0ea91-3a1f-4b63-9f64-bed977ccede3",
			AccountType: db.AccountTypeEquity,
			ParentId:    "",
			Title:       "Opening Balance",
			Childrens:   []db.Account{},
		},
		// ASSET
		{
			Id:          "82bed124-443d-42ba-9b81-8b306868edeb",
			AccountType: db.AccountTypeAsset,
			ParentId:    "",
			Title:       "Money",
			Childrens:   []db.Account{},
		},
		{
			Id:          "b9a47cc8-6bcb-46c5-a04c-65aed68b44b4",
			AccountType: db.AccountTypeAsset,
			ParentId:    "82bed124-443d-42ba-9b81-8b306868edeb",
			Title:       "Bank Account",
			Childrens:   []db.Account{},
		},
		{
			Id:          "17b3aa12-611c-4ffe-8a81-31978115e7d4",
			AccountType: db.AccountTypeAsset,
			ParentId:    "82bed124-443d-42ba-9b81-8b306868edeb",
			Title:       "Wallet",
			Childrens:   []db.Account{},
		},
		{
			Id:          "e7580092-8484-4d22-9b0c-c5832c64f8ee",
			AccountType: db.AccountTypeAsset,
			ParentId:    "",
			Title:       "Loan",
			Childrens:   []db.Account{},
		},
		{
			Id:          "ecc7fc24-9f3b-42bf-88c0-f40c4495fd48",
			AccountType: db.AccountTypeAsset,
			ParentId:    "e7580092-8484-4d22-9b0c-c5832c64f8ee",
			Title:       "Home",
			Childrens:   []db.Account{},
		},
		// LIABILITY
		{
			Id:          "3c7b20c2-4edd-4a27-b326-a78cf6cc0218",
			AccountType: db.AccountTypeLiability,
			ParentId:    "",
			Title:       "Loan",
			Childrens:   []db.Account{},
		},
		{
			Id:          "234e13fa-8dfe-4483-9c16-8093b1500d41",
			AccountType: db.AccountTypeLiability,
			ParentId:    "3c7b20c2-4edd-4a27-b326-a78cf6cc0218",
			Title:       "Home",
			Childrens:   []db.Account{},
		},
		// INCOME
		{
			Id:          "68661064-e7ae-4829-bce7-f43e4746eaf8",
			AccountType: db.AccountTypeIncome,
			ParentId:    "",
			Title:       "Salary",
			Childrens:   []db.Account{},
		},
		{
			Id:          "47b101ea-2bfc-4e03-963b-0ab0db0a7310",
			AccountType: db.AccountTypeIncome,
			ParentId:    "",
			Title:       "Freelance",
			Childrens:   []db.Account{},
		},
		{
			Id:          "4acc2868-5308-4860-a36c-70e55889d22a",
			AccountType: db.AccountTypeIncome,
			ParentId:    "",
			Title:       "Other",
			Childrens:   []db.Account{},
		},
		// EXPENSES
		{
			Id:          "c2fe4b8c-2efd-481e-9936-a690c8470e71",
			AccountType: db.AccountTypeExpense,
			ParentId:    "",
			Title:       "Work",
			Childrens:   []db.Account{},
		},
		{
			Id:          "6a616027-0385-4e1d-8f84-8827c92dda08",
			AccountType: db.AccountTypeExpense,
			ParentId:    "",
			Title:       "Bank",
			Childrens:   []db.Account{},
		},
		{
			Id:          "4c445eb5-d378-451a-be06-f2a73a9555de",
			AccountType: db.AccountTypeExpense,
			ParentId:    "",
			Title:       "Food",
			Childrens:   []db.Account{},
		},
	}

	store := db.NewAccountStore(con)
	for _, account := range accounts {
		err = store.Insert(&account)
		if err != nil {
			log.Fatal(err)
		}
	}
}
