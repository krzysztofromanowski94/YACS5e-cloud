package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"log"
	"runtime"
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
			LogUnknownError(err)
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
		LogUnknownError(err)
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
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return &pb.Empty{}, status.Errorf(110, returnStr)
		}

		return &pb.Empty{}, status.Errorf(0, "User found")
	}

	returnStr := fmt.Sprint("User ", user.Login, " not found")
	return &pb.Empty{}, status.Errorf(112, returnStr)
}

// ERROR CODES:
// 120: UNKNOWN ERROR
// 121: INVALID CREDENTIALS
// 122: CHARACTER DOES NOT EXISTS
// 123: ERROR GETTING DATA FROM STREAM
// 124: ERROR SENDING DATA TO STREAM
// 125: INCORRECT FLOW
// 126: USER DON'T HAVE THIS CHARACTER
func (server *YACS5eServer) GetCharacter(stream pb.YACS5E_GetCharacterServer) error {
	var (
		user          *pb.TUser
		characterList CharacterList
	)

	streamIn, err := stream.Recv()
	if err != nil {
		returnErr := fmt.Sprint("ERROR GETTING DATA FROM INPUT STREAM: ", err)
		log.Println(returnErr)
		return status.Errorf(123, returnErr)
	}

	// 1. Check credentials

	switch tCharacterUnion := streamIn.Union.(type) {

	case *pb.TCharacter_User:
		var (
			userID int
		)
		err := db.QueryRow(
			"SELECT id FROM users WHERE login=? AND password=? LIMIT 1",
			tCharacterUnion.User.Login,
			tCharacterUnion.User.Password,
		).Scan(&userID)

		switch err {

		case sql.ErrNoRows:
			// User does not exists
			return status.Errorf(121, "INVALID CREDENTIALS")

		case nil:
			// User exists
			user = tCharacterUnion.User
			break

		default:
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return status.Errorf(120, returnStr)
		}

	default:
		return status.Errorf(125, "EXPECTED TYPE IS TCharacter_User")
	}

	// 2. Return character list

	rows, err := db.Query("SELECT id, name, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?)", user.Login)
	if err != nil {
		LogUnknownError(err)
		returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
		return status.Errorf(120, returnStr)
	}

	for rows.Next() {
		var (
			id   int
			name string
			data []byte
		)

		err := rows.Scan(&id, &name, &data)
		if err != nil {
			LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
			return status.Errorf(120, returnStr)
		}

		// unmarshall character blob data to map[string]interface{}
		var characterDataIF interface{}
		err = json.Unmarshal(data, &characterDataIF)
		if err != nil {
			LogUnknownError(err)
		}
		characterData := characterDataIF.(map[string]interface{})

		// add character to list
		switch desc := characterData["description"].(type) {
		case string:
			characterList.Characters = append(characterList.Characters, CharacterInfo{
				id,
				name,
				desc,
			})
		}
	}

	// Select the value of key `characters` to be send. It will result as an anonymous list.
	blob, err := json.Marshal(characterList.Characters)
	if err != nil {
		LogUnknownError(err)
		returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
		return status.Errorf(120, returnStr)
	}

	// Send list
	err = stream.Send(&pb.TCharacter{&pb.TCharacter_Blob{blob}})
	if err != nil {
		LogUnknownError(err)
		returnStr := fmt.Sprint("ERROR SENDING DATA TO STREAM: ", err)
		return status.Errorf(124, returnStr)
	}

	// Get character id
	var characterId int

	streamIn, err = stream.Recv()
	if err != nil {
		returnErr := fmt.Sprint("ERROR GETTING DATA FROM INPUT STREAM: ", err)
		log.Println(returnErr)
		return status.Errorf(123, returnErr)
	}

	switch tCharacterUnion := streamIn.Union.(type) {
	case *pb.TCharacter_Id:
		characterId = int(tCharacterUnion.Id)

	default:
		return status.Errorf(125, "EXPECTED TYPE IS TCharacter_Id")
	}

	// Get selected character from database
	var characterData []byte

	err = db.QueryRow(
		"SELECT data FROM characters WHERE id=? AND users_id=(SELECT id FROM users WHERE login=?)",
		characterId,
		user.Login,
	).Scan(&characterData)

	if err == sql.ErrNoRows {
		return status.Errorf(126, "USER DON'T HAVE THIS CHARACTER")
	}

	// Send character to client
	err = stream.Send(&pb.TCharacter{&pb.TCharacter_Blob{characterData}})
	if err != nil {
		LogUnknownError(err)
		returnStr := fmt.Sprint("ERROR SENDING DATA TO STREAM", err)
		return status.Errorf(124, returnStr)
	}

	return status.Errorf(0, "Sent the character")
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

	//testServer := newServer()
	//
	//log.Println(testServer.GetCharacter(nil))

}

func LogUnknownError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("ERROR in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

type YACS5eServer struct {
}

type CharacterList struct {
	Characters []CharacterInfo `json:"characters"`
}

type CharacterInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)
}
