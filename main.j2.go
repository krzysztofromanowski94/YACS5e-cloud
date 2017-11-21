package main

import (
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"strconv"
)

var (
	useTls        bool
	tlsCertFile   = "{{ tls_cert_file }}"
	tlsKeyFile    = "{{ tls_key_file }}"
	serverAddress = ":{{ server_port }}"
	serverOpts    []grpc.ServerOption

	exitProgramChannel = make(chan bool, 1)
)

func main() {
	log.Println("Server starting...")

	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal("Failed to create listen service. ERROR: ", err)
	}

	if useTls {
		cred, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		if err != nil {
			log.Fatal("Failed to generate credentials. ERROR: ", err)
		}
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

func init() {
	var (
		err error
	)
	useTls, err = strconv.ParseBool("{{ use_tls }}")
	if err != nil {
		log.Fatal("Error parsing use_tls ansible variable. ERROR:", err)
	}

}
