package main

import (
	pb "github.com/krzysztofromanowski94/YACS5e-cloud/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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
// ERROR CODES:
// 100: UNKNOWN ERROR
// 101: INVALID LOGIN
// 102: INVALID PASSWORD
// 103: USER EXISTS
func (server *YACS5eServer) Registration(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Println("Registration Context: ", ctx)
	return &pb.Empty{}, status.Errorf(0, "Got user: ", user.Login)
}

// ERROR CODES:
// 110: UNKNOWN ERROR
// 111: INVALID LOGIN
// 112: INVALID PASSWORD
// 113: USER EXISTS
func (server *YACS5eServer) Login(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	log.Println("Login Context: ", ctx)
	return &pb.Empty{}, status.Errorf(0, "User ", user.Login, " may exists. I don't know yet ;x")
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
