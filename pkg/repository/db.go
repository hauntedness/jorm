package repository

import (
	"database/sql"
	"log"
)

func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "root:295449448@tcp(localhost:3306)/helloworld?parseTime=true")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("err ping")
		panic(err)
	}
	log.Println("ping success")
	return db
}

var db = GetDB()
