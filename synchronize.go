package main

import (
	"fmt"
	"log"
	"strings"

	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"github.com/krzysztofromanowski94/YACS5e-cloud/utils"
	"google.golang.org/grpc/status"
)

// Synchronize...
func (server *YACS5eServer) Synchronize(stream pb.YACS5E_SynchronizeServer) error {
	var (
		user *pb.TUser
	)

	log.Println("Start synchronize task")

	// Check recaptcha

	streamIn, err := stream.Recv()
	if err != nil {
		utils.LogUnknownError(err)
		returnStr := fmt.Sprint("Synchronize: ERROR GETTING DATA FROM INPUT STREAM: ", err)
		return status.Errorf(54, returnStr)
	}

	// 1. Check credentials

	user, err = partialLogin(streamIn)
	if err != nil {
		log.Println("Synchronize: Error logging user:", err)
		return err
	}

	log.Println("Synchronize: Logged user", user.Login)

	err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Good{Good: true}})
	if err != nil {
		utils.LogUnknownError(err)
		returnStr := fmt.Sprint("Synchronize: ERROR SENDING DATA FROM INPUT STREAM: ", err)
		return status.Errorf(55, returnStr)
	}

	// 2. Perform char sync one-by-one

	log.Println("Synchronize: Perform char sync one-by-one")

	// var (
	// 	clientTimestampList = make([]*pb.TCharacter, 0)
	// 	serverCharList = make()
	// )

	// example list of characters
	// TTalk_Character[] characterList =

	// change character for character
	exchangeCharInfo := true
	for exchangeCharInfo {

		// get uuid and timestamp
		streamIn, err := stream.Recv()
		if err != nil {
			return utils.ErrorStatus(err)
		}

		log.Println("streamIn:", streamIn)

		switch ttalk := streamIn.Union.(type) {

		case *pb.TTalk_Character:

			var (
				timestamp uint64
				data      []byte
			)

			log.Println("Synchronize: Trying to get data for user", user.Login, ttalk.Character.Uuid)

			err := db.QueryRow(
				"SELECT timestamp, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?) AND uuid=? LIMIT 1",
				user.Login,
				ttalk.Character.Uuid,
			).Scan(&timestamp, &data)

			if err != nil {
				return utils.ErrorStatus(err)
			}

			log.Println("Character match:", timestamp, data)
			err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Blob: data, Timestamp: timestamp}}})
			if err != nil {
				log.Println("Synchronize return updated character error:", err)
			}

			// if ttalk.Character.Timestamp != 0 {
			// 	clientTimestampList = append(clientTimestampList, ttalk.Character)
			// 	err = stream.Send(streamIn)
			// 	if err != nil {
			// 		log.Println("Synchronize return updated character error:", err)
			// 	}
			// } else {
			// 	exchangeCharInfo = false
			// 	break
			// }

		case *pb.TTalk_Good:
			log.Println("Synchronize: no more characters on client")
			exchangeCharInfo = false
			continue

		default:
			return status.Errorf(125, "Expected type TTalk_Character")
		}

	}

	// 3a. Get timestamps from database

	var (
		timestampsDB []*pb.TCharacter
	)

	log.Println(user)
	timestampsQuery, err := db.Query("SELECT timestamp, uuid, data FROM characters WHERE users_id=?", user.Id)
	if err != nil {
		switch strErr := err.Error(); {
		case strings.Contains(strErr, "Error 1062"):
			// User don't have any characters in database
			break
		default:
			utils.LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR:", err)
			return status.Errorf(120, returnStr)
		}
	}

	for timestampsQuery.Next() {
		var (
			timestamp uint64
			uuid      string
			blob      []byte
		)
		err = timestampsQuery.Scan(
			&timestamp,
			&uuid,
			&blob,
		)
		if err != nil {
			utils.LogUnknownError(err)
			returnStr := fmt.Sprint("UNKNOWN ERROR:", err)
			return status.Errorf(120, returnStr)
		}
		timestampsDB = append(timestampsDB, &pb.TCharacter{Uuid: uuid, Timestamp: timestamp, Blob: blob})
	}

	// 3b. Decide what characters need to be updated on server

	log.Println(timestampsDB)

	return status.Errorf(0, "")
}
