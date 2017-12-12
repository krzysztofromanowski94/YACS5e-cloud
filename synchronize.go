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

	// 2a. Create slice of uuids'. If after app-sync there will be any left, app does not have them.
	uuidQuery, err := db.Query("SELECT uuid FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?)", user.Login)
	if err != nil {
		utils.ErrorStatus(err)
	}

	uuidSlice := make([]string, 0)
	for uuidQuery.Next() {
		var uuid string

		err := uuidQuery.Scan(&uuid)
		if err != nil {
			utils.ErrorStatus(err)
		}

		uuidSlice = append(uuidSlice, uuid)
	}

	uuidSlice = append(uuidSlice, "asd")

	log.Println("This user has theese characters on server dv:")
	log.Println(uuidSlice)

	uuidSlice = utils.RemoveFromSlice(uuidSlice, "asd")

	log.Println("After remove:")
	log.Println(uuidSlice)

	// 2b. Perform char sync one-by-one

	log.Println("Synchronize: Perform char sync one-by-one")

	exchangeCharInfo := true
	for exchangeCharInfo {

		// get login, uuid
		streamIn, err := stream.Recv()
		if err != nil {
			return utils.ErrorStatus(err)
		}

		switch ttalk := streamIn.Union.(type) {
		case *pb.TTalk_Character:

			var (
				uuid     string
				lastSync uint64
				lastMod  uint64
				data     []byte
			)

			log.Println("Synchronize: Trying to get data for user", user.Login, ttalk.Character.Uuid)

			err := db.QueryRow(
				"SELECT uuid, last_sync, last_mod, data FROM characters WHERE users_id=(SELECT id FROM users WHERE login=?) AND uuid=? LIMIT 1",
				user.Login,
				ttalk.Character.Uuid,
			).Scan(&uuid, &lastSync, &lastMod, &data)

			if err == sql.ErrNoRows {
				log.Println("Synchronize: character not found on server, ask for complete data")
				// 4 - not on server - receive empty uuid, send complete character
				onCharacterNotFound(stream, *user)
				break
			} else if err != nil {
				return utils.ErrorStatus(err)
			}

			err = stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Uuid: uuid, LastSync: lastSync, LastMod: lastMod}}})
			if err != nil {
				return utils.ErrorStatus(err)
			}

			// Character is even
			if lastSync == ttalk.Character.GetLastSync() && lastMod == ttalk.Character.GetLastMod() {
				log.Println("Synchronize: Character is even (0)", uuid, lastSync, data)
				break
			}

			// if not even - app wants to send data
			streamIn, err = stream.Recv()
			if err != nil {
				return utils.ErrorStatus(err)
			}
			switch ttalk := streamIn.Union.(type) {
			case *pb.TTalk_Character:
				tChar := ttalk.Character
				if tChar.GetLastSync() != 0 && tChar.GetLastMod() != 0 && tChar.Uuid != "" && len(tChar.Blob) > 0 {

					log.Println("Synchronize: app wants to insert / update character uuid: " + tChar.Uuid)

					_, err := db.Exec("INSERT INTO characters "+
						"SET uuid=?, users_id=(SELECT id FROM users WHERE login=?), last_sync=?, last_mod=?, data=? "+
						"ON DUPLICATE KEY UPDATE last_sync=?, last_mod=?, data=?",
						tChar.Uuid, user.Login, tChar.LastSync, tChar.LastMod, tChar.Blob, tChar.LastSync, tChar.LastMod, tChar.Blob)

					if err != nil {
						return utils.ErrorStatus(err)
					}
					continue

				} else if tChar.LastMod == 0 && tChar.LastSync == 0 {
					log.Println("Synchronize: app asks for data")
					err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{
						Uuid:     uuid,
						LastSync: lastSync,
						LastMod:  lastMod,
						Blob:     data,
					}}})
					if err != nil {
						return utils.ErrorStatus(err)
					}
					continue
				}
			}

			log.Println("Synchronize: Unimplemented...")
			log.Println(streamIn)
			log.Println(lastSync, lastMod)

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

// 4 - not on server - receive empty uuid, send complete character
func onCharacterNotFound(stream pb.YACS5E_SynchronizeServer, user pb.TUser) error {
	err := stream.Send(&pb.TTalk{Union: &pb.TTalk_Character{Character: &pb.TCharacter{Uuid: ""}}})
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
				"SET uuid=?, users_id=(SELECT id FROM users WHERE login=?), last_sync=?, last_mod=?, data=?",
			char.Uuid,
			user.Login,
			char.LastSync,
			char.LastMod,
			char.Blob,
		)
		if err == sql.ErrNoRows {
			log.Println("Synchronize: internar error:", err)
			return status.Errorf(2, "Insert new character to db internal error")
		} else if err != nil {
			return utils.ErrorStatus(err)
		}
		log.Println("Synchronize: New character uuid:", char.Uuid)
	}

	return nil
}
