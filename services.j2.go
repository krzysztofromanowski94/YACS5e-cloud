package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

var (
	sqlDBName                = "{{ sql_dbname }}"
	sqlHostname              = "{{ sql_hostname }}"
	sqlPassword              = "{{ sql_password }}"
	sqlPort                  = "{{ sql_port }}"
	sqlUser                  = "{{ sql_user }}"
	sqlMaxOpenConnections, _ = strconv.ParseInt("{{ sql_max_open_connections }}", 10, 64)

	db             *sql.DB
	dataSourceName = sqlUser + ":" + sqlPassword +
		"@tcp(" + sqlHostname + ":" + sqlPort + ")/" + sqlDBName
)

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

type YACS5eServer struct {
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}

// rpc Registration (User) returns (Empty)
// ERROR CODES:
// 100: UNKNOWN ERROR
// 101: INVALID LOGIN
// 102: INVALID PASSWORD
// 103: USER EXISTS
func (server *YACS5eServer) Registration(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {
	log.Println("Registration Context: ", ctx)

	return &pb.Empty{}, status.Errorf(0, "Registered user: ", user.Login)
}

// rpc Login (User) returns (Empty)
// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID LOGIN
// 112: INVALID PASSWORD
// 113: USER EXISTS
func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {
	log.Println("Login Context: ", ctx)
	return &pb.Empty{}, status.Errorf(0, "User ", user.Login, " may exists. I don't know yet ;x")
}
