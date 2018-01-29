package main

import (
	"github.com/appleboy/go-fcm"
	"fmt"
)


func pushLoop() {
	deviceToken := "dupadupa"

	msg := &fcm.Message{
		To: deviceToken,
		Data: map[string]interface{}{
			"myfield1": "myVal1",
		},
	}

	fmt.Println(msg)
}

func serve(){

}

