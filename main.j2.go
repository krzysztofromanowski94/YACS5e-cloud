package main

import (
	"fmt"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

var (
	tlsCertFile = "{{ tls_cert_file }}"
	tlsKeyFile  = "{{ tls_key_file }}"
	port        = "{{ port }}"
)

func main() {
	cred, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
	if err != nil {
		log.Fatal("Failed to generate credentials %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to create listen service %v", err)
	}

	server := grpc.NewServer(grpc.Creds(cred))

	server.Serve(lis)

	dupa := pb.Register{&pb.Register_Valid{true}}
	fmt.Println(dupa)
}
