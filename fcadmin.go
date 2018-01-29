package main

import (
	"log"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func hello() {

	opt := option.WithCredentialsFile("{{ firebase_account_key }}")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Printf("ERROR: can't initialize fcadmin: %v\n", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		LogUnknownError(err)
	}

	client.CustomToken("asd")
}