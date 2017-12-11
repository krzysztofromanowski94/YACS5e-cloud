package main

import (
	"database/sql"
	"fmt"
	"log"

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

	exchangeCharInfo := true
	for exchangeCharInfo {

		// get uuid and timestamp
		streamIn, err := stream.Recv()
		if err != nil {
			return utils.ErrorStatus(err)
		}

		log.Println("Synchronize: streamIn:", streamIn)

		switch ttalk := streamIn.Union.(type) {
		case *pb.TTalk_Character:

			var (
				name      string
				timestamp uint64
				data      []byte
			)

			log.Println("Synchronize: Trying to get data for user", user.Login, ttalk.Character.Uuid)

			err := db.QueryRow(
				"SELECT name, timestamp, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?) AND uuid=? LIMIT 1",
				user.Login,
				ttalk.Character.Uuid,
			).Scan(&name, &timestamp, &data)

			if err == sql.ErrNoRows {
				log.Println("Synchronize: character not found on server, ask for complete data")
				// 3 - not on server
				onCharacterNotFound(stream, *user)
			} else if err != nil {
				return utils.ErrorStatus(err)
			}

			// 0. Characters are synced
			if timestamp == ttalk.Character.GetTimestamp() {
				log.Println("Synchronize: Characters are even")
				err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Blob: data, Timestamp: timestamp}}})
				if err != nil {
					log.Println("Synchronize return updated character error:", err)
				}
			}

			//log.Println("Synchronize: Character match:", timestamp, data)
			//err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Blob: data, Timestamp: timestamp}}})
			//if err != nil {
			//	log.Println("Synchronize return updated character error:", err)
			//}

			// 	clientTimestampList = append(clientTimestampList, ttalk.Character)

		case *pb.TTalk_Good:
			log.Println("Synchronize: no more characters on client")
			exchangeCharInfo = false
			continue

		default:
			return status.Errorf(125, "Expected type TTalk_Character")
		}

	}

	return status.Errorf(0, "")
}

func onCharacterFound(stream pb.YACS5E_SynchronizeServer) {

}

// 3 - not on server - send timestamp == 0, receive complete character
func onCharacterNotFound(stream pb.YACS5E_SynchronizeServer, user pb.TUser) error {
	err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Timestamp: 0}}})
	if err != nil {
		return utils.ErrorStatus(err)
	}

	streamIn, err := stream.Recv()
	if err != nil {
		return utils.ErrorStatus(err)
	}

	switch tCharacter := streamIn.Union.(type) {
	case *pb.TTalk_Character:
		char := tCharacter.Character
		_, err := db.Exec(
			"INSERT INTO characters "+
				"SET uuid=?, users_id=(SELECT id FROM users WHERE login=?), timestamp=?, data=?",
			char.Uuid,
			user,
			char.Timestamp,
			char.Blob,
		)
		if err == sql.ErrNoRows {
			log.Println("Synchronize: internar error:", err)
			return status.Errorf(2, "Insert new character to db internal error")
		} else if err != nil {
			return utils.ErrorStatus(err)
		}
	}

	return nil
}
