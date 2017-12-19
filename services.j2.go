package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

var (
	sqlDBName   = "{{ sql_dbname }}"
	sqlHostname = "{{ sql_hostname }}"
	sqlPassword = "{{ sql_password }}"
	sqlPort     = "{{ sql_port }}"
	sqlUser     = "{{ sql_user }}"
	db          *sql.DB
)

func init() {
	sqlMaxOpenConnections, err := strconv.ParseInt("{{ sql_max_open_connections }}", 10, 64)
	if err != nil {
		log.Fatalln("error parsing sql_max_open_connections")
	}

	dataSourceName := sqlUser + ":" + sqlPassword + "@tcp(" + sqlHostname + ":" + sqlPort + ")/" + sqlDBName

	log.Println("Connecting to: ", sqlHostname)

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Preparing the database connection caused ERROR: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Connecting to database caused ERROR: ", err)
	}

	db.SetMaxOpenConns(int(sqlMaxOpenConnections))

	log.Println("Connection estabilished")
}

type YACS5eServer struct {
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}
