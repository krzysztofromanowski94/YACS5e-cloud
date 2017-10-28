package main

import (
	"database/sql"
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

	user := &pb.TUser{"testlogin", "12345", "token", "Tester"}
	user2 := &pb.TUser{"testlogin2", "12345", "token2", "Tester2"}

	log.Println("Testing the registration functionality with user:\n", user)

	testserver := YACS5eServer{}
	testserver.Registration(nil, user)
	testserver.Login(nil, user)
	testserver.Login(nil, user2)

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
			log.Println("User ", user.Login, " exists")
			return &pb.Empty{}, status.Errorf(103, "User", user.Login, "exists.")

		default:
			log.Fatal("Registering user caused unknown ERROR: ", err)
			return &pb.Empty{}, status.Errorf(1, "Unknown error: ", err)
		}
	}

	return &pb.Empty{}, status.Errorf(0, "Registered user: ", user.Login)
}

// rpc Login (User) returns (Empty)
// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID CREDENTIALS
func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	row := db.QueryRow("SELECT login, visible_name FROM users WHERE login=? AND password=?", user.Login, user.Password)

	var (
		login       string
		visibleName string
	)

	err := row.Scan(&login, &visibleName)

	if err != nil {
		switch strErr := err.Error(); {

		case strings.Contains(strErr, "sql: no rows in result set"):
			log.Println("Wrong login and/or password for user: ", user.Login)
			return &pb.Empty{}, status.Errorf(111, "INVALID CREDENTIALS")

		default:
			log.Fatal("Logging user caused unknown ERROR: ", err)
			return &pb.Empty{}, status.Errorf(110, "UNKNOWN ERROR: ", err)
		}
	}

	log.Println("After query: ", login, visibleName)

	return &pb.Empty{}, status.Errorf(0, "User ", user.Login, " may exists. I don't know yet ;x")
}
