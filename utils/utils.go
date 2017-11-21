package utils

import (
	"log"
	"runtime"
)

func LogUnknownError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("ERROR in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}
