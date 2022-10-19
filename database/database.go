package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "****"
	password = "****"
	dbname   = "shop"
)

var (
	Db            *sql.DB
	connectionERR error
)

func ConnectDb() {
	psqlCon := fmt.Sprintf("host= %s port= %s user= %s password= %s dbname= %s sslmode=disable", host, port, user, password, dbname)

	Db, connectionERR = sql.Open("postgres", psqlCon)
	if connectionERR != nil {
		log.Fatal("psql connection err :", connectionERR)
	}

}
