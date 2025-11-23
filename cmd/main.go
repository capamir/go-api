package main

import (
	"log"
	"database/sql"

	"github.com/capamir/go-api/cmd/api"
	"github.com/capamir/go-api/configs"
	"github.com/capamir/go-api/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr: 				  configs.Envs.DBAddress,
		DBName: 			  configs.Envs.DBName,
		Net: 				  "tcp",
		AllowNativePasswords: true,
		ParseTime: 			  true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}