package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/jangidRkt08/go-Ecom/cmd/api"
	"github.com/jangidRkt08/go-Ecom/config"
	"github.com/jangidRkt08/go-Ecom/db"
)

func main() {
	db, err := db.NewMySqlStorage(mysql.Config{
		User: config.Envs.DBUser,
		Passwd : config.Envs.DBPassword,
		Addr : config.Envs.DBAddress,
		DBName: config.Envs.DBName,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	})

	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	server := api.NewAPIserver(":8080", db)
	log.Println("Starting server on port 8080")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}


func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB:Successfully Connected!")
}

