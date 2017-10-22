package main

import (
	"fmt"
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	tls, parseErr = strconv.ParseBool("{{ use_tls }}")
	tlsCertFile   = "{{ tls_cert_file }}"
	tlsKeyFile    = "{{ tls_key_file }}"
	serverAddress = ":{{ server_port }}"
	serverOpts    []grpc.ServerOption
)

type YACS5eServer struct {
}

// rpc Registration (stream Register) returns (stream Register)
func (server *YACS5eServer) Registration(stream pb.YACS5E_RegistrationServer) error {
	for {
		streamIn, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Println("Error in Registration Server stream in. ERROR: ", err)
			return err
		}
		fmt.Println(streamIn)
		err = stream.Send(&pb.Register{&pb.Register_Valid{true}})
		if err != nil {
			log.Println("Error in Registration Server stream out. ERROR: ", err)
			return err
		}
	}
}

func newServer() *YACS5eServer {
	return new(YACS5eServer)

}

func main() {
	if parseErr != nil {
		log.Fatal("Error parsing use_tls ansible variable. ERROR:", parseErr)
	}
	lis, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal("Failed to create listen service. ERROR: ", err)
	}

	if tls {
		cred, err := credentials.NewServerTLSFromFile(tlsCertFile, tlsKeyFile)
		if err != nil {
			log.Fatal("Failed to generate credentials. ERROR: ", err)
		}
		serverOpts = append(serverOpts, grpc.Creds(cred))
	}

	grpcServer := grpc.NewServer(serverOpts...)
	pb.RegisterYACS5EServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
