package utils

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"runtime"
)

func GetTLS(host, dirCache string) *tls.Config {
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(dirCache),
		HostPolicy: autocert.HostWhitelist(host),
		Email:      "theamazingptp@gmail.com",
	}
	return &tls.Config{GetCertificate: manager.GetCertificate}
}

func LogUnknownError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("ERROR in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}
