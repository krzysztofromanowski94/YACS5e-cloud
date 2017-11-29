package utils

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"runtime"

	"google.golang.org/grpc/status"
)

func LogUnknownError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("ERROR in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

func ErrorStatus(err error) error {
	switch err {

	case nil:
		return nil

	case sql.ErrNoRows:
		return status.Errorf(52, "No result")

	case io.EOF:
		log.Println("Synchronize: EOF")
		return status.Errorf(2, "EOF too soon")

	default:
		LogUnknownError(err)
		returnStr := fmt.Sprint("UNKNOWN ERROR: ", err)
		return status.Errorf(51, returnStr)
	}
}
