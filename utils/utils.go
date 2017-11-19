package utils

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
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
