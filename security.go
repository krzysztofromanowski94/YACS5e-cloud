package main

import (
	"github.com/krzysztofromanowski94/YACS5e-cloud/utils"
	"google.golang.org/grpc/status"

	"context"
	"database/sql"
	"fmt"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/crypto/blake2b"
	"log"
)

func partialLogin(tTalk *pb.TTalk) (user *pb.TUser, err error) {

	// Here should be checking if recaptcha is right

	switch tUnion := tTalk.Union.(type) {

	case *pb.TTalk_User:

		if tUnion.User.Login == "" || tUnion.User.Password == ""{
			log.Println("partialLogin: user login or passwd is null")
		}

		hashedPasswd := blake2b.Sum512([]byte(tUnion.User.Password))

		err := db.QueryRow(
			"SELECT id FROM users WHERE login=? AND password=? LIMIT 1",
			tUnion.User.Login,
			append(make([]byte, 0), hashedPasswd[:]...),
		).Scan(&tUnion.User.Id)

		if err != nil {
			return nil, utils.ErrorStatus(err)
		}

		return tUnion.User, nil

	default:
		return nil, status.Errorf(53, "Unexpected type")
	}

	return nil, status.Errorf(2, "Unexpected return at partial login")
}

func (server *YACS5eServer) Registration(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	log.Println("Registration: Trying to register user ", user.Login)

	hashedPasswd := blake2b.Sum512([]byte(user.Password))

	_, err := db.Exec("INSERT INTO users VALUES (NULL, ?, ?, ?)",
		user.Login,
		append(make([]byte, 0), hashedPasswd[:]...),
		user.VisibleName,
		)
	if err != nil {
		switch err {

		case sql.ErrNoRows:
			returnStr := fmt.Sprint("User ", user.Login, " exists.")
			log.Println(returnStr)
			return &pb.Empty{}, status.Errorf(103, returnStr)

		default:
			utils.LogUnknownError(err)
			return &pb.Empty{}, status.Errorf(100, "Unknown error")
		}
	}

	returnStr := fmt.Sprint("Registered user: ", user.Login)
	log.Println(returnStr)
	return &pb.Empty{}, status.Errorf(0, returnStr)
}

func (server *YACS5eServer) Login(ctx context.Context, user *pb.TUser) (*pb.Empty, error) {

	// Here should be checking if recaptcha is right

	log.Println("Trying to login: ", user.Login)

	hashedPasswd := blake2b.Sum512([]byte(user.Password))

	row, err := db.Query("SELECT login, visible_name FROM users WHERE login=? AND password=? LIMIT 1",
		user.Login,
		append(make([]byte, 0), hashedPasswd[:]...),
		)

	if err != nil {
		utils.LogUnknownError(err)
		return &pb.Empty{}, status.Errorf(110, "Unknown error")
	}

	for row.Next() {
		var (
			login       string
			visibleName string
		)

		err := row.Scan(&login, &visibleName)
		if err != nil {
			utils.LogUnknownError(err)
			return &pb.Empty{}, status.Errorf(110, "Unknown error")
		}

		log.Println("User logged in")
		return &pb.Empty{}, status.Errorf(0, "User found")
	}

	returnStr := fmt.Sprint("User ", user.Login, " not found")
	log.Println(returnStr)
	return &pb.Empty{}, status.Errorf(112, returnStr)
}
