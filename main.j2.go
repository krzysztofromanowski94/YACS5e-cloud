package main

import (
	"crypto/tls"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

var (
	useTls, parseErr = strconv.ParseBool("{{ use_tls }}")
	tlsCertFile      = "{{ tls_cert_file }}"
	tlsKeyFile       = "{{ tls_key_file }}"
	serverAddress    = ":{{ server_port }}"
	serverOpts       []grpc.ServerOption

	exitProgramChannel chan bool = make(chan bool, 1)
)

func main() {
	log.Println("Server starting...")

	if parseErr != nil {
		log.Fatal("Error parsing use_tls ansible variable. ERROR:", parseErr)
	}
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal("Failed to create listen service. ERROR: ", err)
	}

	if useTls {
		//cred, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		//if err != nil {
		//	log.Fatal("Failed to generate credentials. ERROR: ", err)
		//}

		BackendCert, err := ioutil.ReadFile(tlsCertFile)
		if err != nil {
			log.Fatal("Can't get BackendCert ERROR: ", err)
		}
		BackendKey, err := ioutil.ReadFile(tlsKeyFile)
		if err != nil {
			log.Fatal("Can't get BackendKey ERROR: ", err)
		}

		cert, err := tls.X509KeyPair(BackendCert, BackendKey)
		if err != nil {
			log.Fatal("Failed to generate credentials. ERROR: ", err)
		}

		cred := credentials.NewServerTLSFromCert(&cert)

		serverOpts = append(serverOpts, grpc.Creds(cred))
	}

	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterYACS5EServer(grpcServer, newServer())

	go func() {
		log.Println("Server Started")
		grpcServer.Serve(lis)
	}()

	<-exitProgramChannel
	return
}
