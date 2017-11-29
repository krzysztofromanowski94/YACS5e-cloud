package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"github.com/krzysztofromanowski94/YACS5e-cloud/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

var (
	sqlDBName   = "{{ sql_dbname }}"
	sqlHostname = "{{ sql_hostname }}"
	sqlPassword = "{{ sql_password }}"
	sqlPort     = "{{ sql_port }}"
	sqlUser     = "{{ sql_user }}"
	db          *sql.DB
)

func (server *YACS5eServer) Registration(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	log.Println("Trying to register user ", user.Password)

	_, err := db.Exec("INSERT INTO users VALUES (NULL, ?, ?, ?)", user.Login, user.Password, user.VisibleName)
	if err != nil {
		switch err {

		case sql.ErrNoRows:
			returnStr := fmt.Sprint("User ", user.Login, " exists.")
			log.Println(returnStr)
			return &pb.Empty{}, status.Errorf(103, returnStr)

		default:
			utils.LogUnknownError(err)
			return &pb.Empty{}, status.Errorf(100, "Unknown error: ", err)
		}
	}

	returnStr := fmt.Sprint("Registered user: ", user.Login)
	log.Println(returnStr)
	return &pb.Empty{}, status.Errorf(0, returnStr)
}

func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	log.Println("Trying to login: ", user.Login, " ", user.Password)

	row, err := db.Query("SELECT login, visible_name FROM users WHERE login=? AND password=? LIMIT 1", user.Login, user.Password)

	if err != nil {
		utils.LogUnknownError(err)
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
			utils.LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return &pb.Empty{}, status.Errorf(110, returnStr)
		}

		log.Println("User logged in")
		return &pb.Empty{}, status.Errorf(0, "User found")
	}

	returnStr := fmt.Sprint("User ", user.Login, " not found")
	log.Println(returnStr)
	return &pb.Empty{}, status.Errorf(112, returnStr)
}

func init() {
	sqlMaxOpenConnections, err := strconv.ParseInt("{{ sql_max_open_connections }}", 10, 64)
	if err != nil {
		log.Fatalln("error parsing sql_max_open_connections")
	}

	dataSourceName := sqlUser + ":" + sqlPassword + "@tcp(" + sqlHostname + ":" + sqlPort + ")/" + sqlDBName

	log.Println("Connecting to: ", dataSourceName)

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

func partialLogin(tTalk *pb.TTalk) (user *pb.TUser, err error) {
	switch tCharacterUnion := tTalk.Union.(type) {

	case *pb.TTalk_User:
		err := db.QueryRow(
			"SELECT id FROM users WHERE login=? AND password=? LIMIT 1",
			tCharacterUnion.User.Login,
			tCharacterUnion.User.Password,
		).Scan(&tCharacterUnion.User.Id)

		if err != nil {
			return nil, utils.ErrorStatus(err)
		}

		return tCharacterUnion.User, nil

	default:
		return nil, status.Errorf(53, "EXPECTED TYPE IS TTalk_User")
	}

	return nil, status.Errorf(2, "UNEXPECTED RETURN AT PARTIAL LOGIN")
}

type YACS5eServer struct {
}

//type CharacterList struct {
//	Characters []CharacterInfo `json:"characters"`
//}

//type CharacterInfo struct {
//	ID          int    `json:"id"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}
