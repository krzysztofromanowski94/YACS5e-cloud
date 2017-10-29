package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"strings"
)

var (
	sqlDBName                = "{{ sql_dbname }}"
	sqlHostname              = "{{ sql_hostname }}"
	sqlPassword              = "{{ sql_password }}"
	sqlPort                  = "{{ sql_port }}"
	sqlUser                  = "{{ sql_user }}"
	sqlMaxOpenConnections, _ = strconv.ParseInt("{{ sql_max_open_connections }}", 10, 64)
	db                       *sql.DB
	dataSourceName           = sqlUser + ":" + sqlPassword +
		"@tcp(" + sqlHostname + ":" + sqlPort + ")/" + sqlDBName
)

type YACS5eServer struct {
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}

func init() {
	log.Println("Connecting to: ", dataSourceName)

	var err error

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

// rpc Registration (User) returns (Empty)
// ERROR CODES:
// 100: UNKNOWN ERROR
// 101: INVALID LOGIN
// 102: INVALID PASSWORD
// 103: USER EXISTS
func (server *YACS5eServer) Registration(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	_, err := db.Exec("INSERT INTO users VALUES (NULL, ?, ?, ?)", user.Login, user.Password, user.VisibleName)
	if err != nil {
		switch strErr := err.Error(); {

		case strings.Contains(strErr, "Error 1062"):
			returnStr := fmt.Sprint("User ", user.Login, " exists.")
			return &pb.Empty{}, status.Errorf(103, returnStr)

		default:
			log.Fatal("Registering user caused unknown ERROR: ", err)
			return &pb.Empty{}, status.Errorf(100, "Unknown error: ", err)
		}
	}

	returnStr := fmt.Sprint("Registered user: ", user.Login)
	return &pb.Empty{}, status.Errorf(0, returnStr)
}

// rpc Login (User) returns (Empty)
// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID CREDENTIALS
// 112: USER NOT FOUND
func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	row, err := db.Query("SELECT login, visible_name FROM users WHERE login=? AND password=? LIMIT 1", user.Login, user.Password)

	if err != nil {
		returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
		return &pb.Empty{}, status.Errorf(110, returnStr)
	}

	for row.Next() {
		var (
			login       string
			visibleName string
		)

		err := row.Scan(&login, &visibleName)
		if err != nil {
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return &pb.Empty{}, status.Errorf(110, returnStr)
		}

		return &pb.Empty{}, status.Errorf(0, "User found")
	}

	returnStr := fmt.Sprint("User ", user.Login, " not found")
	return &pb.Empty{}, status.Errorf(112, returnStr)
}
